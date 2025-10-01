package kollama

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/samber/lo"
	"github.com/spf13/cobra"

	ollamav1 "github.com/nekomeowww/ollama-operator/api/ollama/v1"
	"github.com/nekomeowww/ollama-operator/pkg/model"
)

var (
	ErrOllamaModelNotSupported = fmt.Errorf("%s is not supported on the cluster, did you install the Ollama Operator?", modelSchemaGroupVersionResource.String())
)

func isKubectlPlugin() (bool, error) {
	exec, err := os.Executable()
	if err != nil {
		return false, err
	}

	basename := filepath.Base(exec)
	if strings.HasPrefix(basename, "kubectl-") {
		return true, nil
	}

	return false, nil
}

func command() string {
	is, _ := isKubectlPlugin()
	if is {
		return "kubectl kollama"
	}

	return "kollama"
}

func IsOllamaOperatorCRDSupported(discoveryClient discovery.DiscoveryInterface, resourceName string) (bool, error) {
	groupVersion := schemaGroupVersion.String()

	list, err := discoveryClient.ServerResourcesForGroupVersion(groupVersion)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return false, nil
		}

		// don't record, just attempt again next time in case it's a transient error
		return false, err
	}

	for _, resources := range list.APIResources {
		if resources.Name == resourceName {
			return true, nil
		}
	}

	return false, nil
}

func FromUnstructured[T any](obj *unstructured.Unstructured) (*T, error) {
	typedObj := new(T)

	err := runtime.DefaultUnstructuredConverter.FromUnstructured(obj.UnstructuredContent(), typedObj)
	if err != nil {
		return nil, err
	}

	return typedObj, nil
}

func Unstructured[T runtime.Object](obj T) (*unstructured.Unstructured, error) {
	unstructured := &unstructured.Unstructured{}

	converted, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	if err != nil {
		return nil, err
	}

	unstructured.SetUnstructuredContent(converted)
	unstructured.SetAPIVersion(schemaGroupVersion.String())
	unstructured.SetKind("Model")

	return unstructured, nil
}

func getOllama(ctx context.Context, dynamicClient dynamic.Interface, namespace string, name string) (*ollamav1.Model, error) {
	unstructuredObj, err := dynamicClient.
		Resource(modelSchemaGroupVersionResource).
		Namespace(namespace).
		Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			return nil, nil
		}

		return nil, err
	}

	model, err := FromUnstructured[ollamav1.Model](unstructuredObj)
	if err != nil {
		return nil, err
	}

	return model, nil
}

func getNamespace(clientConfig clientcmd.ClientConfig, cmd *cobra.Command) (string, error) {
	namespace, err := cmd.Flags().GetString("namespace")
	if err != nil {
		return "", err
	}

	if namespace == "" {
		var ok bool

		namespace, ok, err = clientConfig.Namespace()
		if err != nil {
			return "", err
		}

		if !ok {
			namespace = "default"
		}
	}

	return namespace, nil
}

func getImage(cmd *cobra.Command, args []string) (string, error) {
	modelName := args[0]

	modelImage, err := cmd.Flags().GetString("image")
	if err != nil {
		return "", err
	}

	if modelImage == "" {
		return modelName, nil
	}

	return modelImage, nil
}

func createOllamaModel(
	ctx context.Context,
	dynamicClient dynamic.Interface,
	namespace string,
	name string,
	image string,
	resources corev1.ResourceRequirements,
	storageClass string,
	pvAccessMode string,
) (*ollamav1.Model, error) {
	model := &ollamav1.Model{
		TypeMeta: metav1.TypeMeta{
			APIVersion: schemaGroupVersion.String(),
			Kind:       "Model",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Labels: map[string]string{
				"ollama.ayaka.io/managed-by": "kollama",
			},
		},
		Spec: ollamav1.ModelSpec{
			Image:     image,
			Resources: resources,
		},
	}
	if storageClass != "" {
		model.Spec.StorageClassName = lo.ToPtr(storageClass)
	}

	if pvAccessMode != "" {
		model.Spec.PersistentVolume = &ollamav1.ModelPersistentVolumeSpec{
			AccessMode: lo.ToPtr(corev1.PersistentVolumeAccessMode(pvAccessMode)),
		}
	}

	unstructuredObj, err := Unstructured(model)
	if err != nil {
		return nil, err
	}

	_, err = dynamicClient.
		Resource(modelSchemaGroupVersionResource).
		Namespace(namespace).
		Create(ctx, unstructuredObj, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}

	return model, nil
}

func exposeOllamaModel(
	ctx context.Context,
	kubeClient client.Client,
	namespace string,
	name string,
	serviceType corev1.ServiceType,
	serviceName string,
	nodePort int32,
) (*corev1.Service, error) {
	var deployment appsv1.Deployment

	err := kubeClient.Get(ctx, types.NamespacedName{Namespace: namespace, Name: model.ModelAppName(name)}, &deployment)
	if err != nil {
		if apierrors.IsNotFound(err) {
			time.Sleep(1 * time.Second)
			return exposeOllamaModel(ctx, kubeClient, namespace, name, serviceType, serviceName, nodePort)
		}

		return nil, err
	}

	svc := model.NewServiceForModel(namespace, name, &deployment, serviceType)
	if serviceName != "" {
		svc.Name = serviceName
	} else {
		switch serviceType {
		case corev1.ServiceTypeNodePort:
			svc.Name = model.ModelAppName(name) + "-nodeport"
		case corev1.ServiceTypeLoadBalancer:
			svc.Name = model.ModelAppName(name) + "-lb"
		case corev1.ServiceTypeExternalName:
			svc.Name = model.ModelAppName(name) + "-external"
		case corev1.ServiceTypeClusterIP:
			svc.Name = model.ModelAppName(name)
		default:
			return nil, fmt.Errorf("unsupported service type: %s", serviceType)
		}
	}

	if serviceType == "NodePort" && nodePort != 0 {
		if nodePort < 0 || nodePort > 65535 {
			return nil, fmt.Errorf("invalid nodePort: %d", nodePort)
		}

		svc.Spec.Ports[0].NodePort = nodePort
	}

	var existingService corev1.Service

	err = kubeClient.Get(ctx, types.NamespacedName{Namespace: namespace, Name: svc.Name}, &existingService)
	if err != nil && !apierrors.IsNotFound(err) {
		return nil, err
	}

	if err == nil {
		return &existingService, nil
	}

	err = kubeClient.Create(ctx, svc)
	if err != nil {
		return nil, err
	}

	return svc, nil
}
