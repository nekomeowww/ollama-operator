package model

import (
	"context"
	"fmt"
	"path/filepath"

	namepkg "github.com/google/go-containerregistry/pkg/name"
	"github.com/samber/lo"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	ollamav1 "github.com/nekomeowww/ollama-operator/api/ollama/v1"
	"github.com/nekomeowww/xo"
)

func getServiceByLabels(ctx context.Context, c client.Client, namespace string, l labels.Set) (*corev1.Service, error) {
	var service corev1.ServiceList

	err := c.List(ctx, &service, &client.ListOptions{
		Namespace:     namespace,
		LabelSelector: labels.SelectorFromValidatedSet(l),
	})
	if err != nil {
		return nil, err
	}
	if len(service.Items) == 0 {
		return nil, nil
	}

	return &service.Items[0], nil
}

func ModelServiceName(name string) string {
	return fmt.Sprintf("ollama-srv-%s", xo.RandomHashString(6))
}

func ModelAppName(name string) string {
	return fmt.Sprintf("ollama-model-%s", name)
}

func ModelLabels(name string) map[string]string {
	return map[string]string{
		"app":                        ModelAppName(name),
		"ollama.ayaka.io/type":       "model",
		"model.ollama.ayaka.io":      name,
		"model.ollama.ayaka.io/name": name,
	}
}

func OllamaModelNameFromNameReference(ref namepkg.Reference) string {
	parsedModelName := filepath.Base(ref.Context().RepositoryStr())
	if ref.Identifier() != "latest" {
		parsedModelName += ":" + ref.Identifier()
	}

	return parsedModelName
}

func ImageStoreLabels() map[string]string {
	return map[string]string{
		"app":                  "ollama-image-store",
		"ollama.ayaka.io/type": "image-store",
	}
}

func ModelAnnotations(name string, imageStore bool) map[string]string {
	return map[string]string{}
}

func getDeployment(ctx context.Context, c client.Client, namespace string, name string) (*appsv1.Deployment, error) {
	var deployment appsv1.Deployment

	err := c.Get(ctx, client.ObjectKey{Namespace: namespace, Name: ModelAppName(name)}, &deployment)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, nil
		}

		return nil, err
	}

	return &deployment, nil
}

func EnsureDeploymentCreated(
	ctx context.Context,
	namespace string,
	name string,
	image string,
	replicas *int32,
	model *ollamav1.Model,
) (*appsv1.Deployment, error) {
	c := ClientFromContext(ctx)
	modelRecorder := WrappedRecorderFromContext[*ollamav1.Model](ctx)

	deployment, err := getDeployment(ctx, c, namespace, name)
	if err != nil {
		return nil, err
	}
	if deployment != nil {
		return deployment, nil
	}

	ref, err := namepkg.ParseReference(
		image,
		namepkg.Insecure,
		namepkg.WithDefaultRegistry("https://registry.ollama.ai"),
		namepkg.WithDefaultTag("latest"),
	)
	if err != nil {
		return nil, err
	}

	deployment = &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:        ModelAppName(name),
			Namespace:   namespace,
			Labels:      ModelLabels(name),
			Annotations: ModelAnnotations(ModelAppName(name), false),
			OwnerReferences: []metav1.OwnerReference{{
				APIVersion:         model.APIVersion,
				Kind:               model.Kind,
				Name:               model.Name,
				UID:                model.UID,
				BlockOwnerDeletion: lo.ToPtr(true),
			}},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: lo.Ternary(replicas == nil, lo.ToPtr(int32(1)), replicas),
			Selector: &metav1.LabelSelector{
				MatchLabels: ModelLabels(name),
			},
			Template: MergePodTemplate(ctx, namespace, name, image, OllamaModelNameFromNameReference(ref), replicas, model),
		},
	}

	err = c.Create(ctx, deployment)
	if err != nil {
		return nil, err
	}

	modelRecorder.Eventf("Normal", "DeploymentCreated", "Deployment %s created", deployment.Name)

	return deployment, nil
}

func MergePodTemplate(
	ctx context.Context,
	namespace string,
	name string,
	image string,
	parsedModelName string,
	replicas *int32,
	model *ollamav1.Model,
) corev1.PodTemplateSpec {
	var pod corev1.PodTemplateSpec

	if model.Spec.PodTemplate != nil {
		pod = *model.Spec.PodTemplate
	}

	pod.ObjectMeta.Labels = lo.Assign(pod.ObjectMeta.Labels, ModelLabels(name))
	pod.ObjectMeta.Annotations = lo.Assign(pod.ObjectMeta.Annotations, ModelAnnotations(ModelAppName(name), false))

	pod.Spec.InitContainers = AssignOrAppend(
		pod.Spec.InitContainers,
		FindOllamaPullerContainer,
		AssignOllamaPullerContainer(name, image, parsedModelName, namespace, model.Spec.Resources, model.Spec.ExtraEnvFrom, model.Spec.Env),
		func() corev1.Container {
			return NewOllamaPullerContainer(name, image, parsedModelName, namespace, model.Spec.Resources, model.Spec.ExtraEnvFrom, model.Spec.Env)
		},
	)
	pod.Spec.Containers = AssignOrAppend(
		pod.Spec.Containers,
		FindOllamaServerContainer,
		AssignOllamaServerContainer(true, model.Spec.Resources, model.Spec.ExtraEnvFrom, model.Spec.Env),
		func() corev1.Container {
			return NewOllamaServerContainer(true, model.Spec.Resources, model.Spec.ExtraEnvFrom, model.Spec.Env)
		},
	)

	pod.Spec.Volumes = AppendIfNotFound(pod.Spec.Volumes, func(item corev1.Volume) bool {
		return item.Name == "image-storage"
	}, func() corev1.Volume {
		return corev1.Volume{
			Name: "image-storage",
			VolumeSource: corev1.VolumeSource{
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
					ClaimName: ImageStorePVCName,
					ReadOnly:  true,
				},
			},
		}
	})

	if model.Spec.RuntimeClassName != nil {
		pod.Spec.RuntimeClassName = model.Spec.RuntimeClassName
	}

	return pod
}

func AssignOrAppend[T any](source []T, predicate func(item T) bool, modifier func(item T, index int) T, newFn func() T) []T {
	target, index, found := lo.FindIndexOf(source, predicate)
	if found {
		target = modifier(target, index)
		source[index] = target
	} else {
		target = newFn()
		source = append(source, target)
	}

	return source
}

func AppendIfNotFound[T any](source []T, predicate func(item T) bool, newFn func() T) []T {
	_, _, found := lo.FindIndexOf(source, predicate)
	if !found {
		source = append(source, newFn())
	}

	return source
}

func IsDeploymentReady(
	ctx context.Context,
	namespace string,
	name string,
) (bool, error) {
	log := log.FromContext(ctx)
	c := ClientFromContext(ctx)
	modelRecorder := WrappedRecorderFromContext[*ollamav1.Model](ctx)

	deployment, err := getDeployment(ctx, c, namespace, name)
	if err != nil {
		return false, err
	}
	if deployment == nil {
		return false, nil
	}

	replica := 1
	if deployment.Spec.Replicas != nil {
		replica = int(*deployment.Spec.Replicas)
	}
	if deployment.Status.ReadyReplicas == int32(replica) {
		log.Info("deployment is ready", "deployment", deployment)
		return true, nil
	}

	log.Info("waiting for deployment to be ready", "deployment", deployment)
	modelRecorder.Eventf("Normal", "WaitingForDeployment", "Waiting for deployment %s to become ready", deployment.Name)

	return false, nil
}

func UpdateDeployment(
	ctx context.Context,
	model *ollamav1.Model,
) (bool, error) {
	c := ClientFromContext(ctx)
	modelRecorder := WrappedRecorderFromContext[*ollamav1.Model](ctx)

	deployment, err := getDeployment(ctx, c, model.Namespace, model.Name)
	if err != nil {
		return false, err
	}
	if deployment == nil {
		return false, nil
	}

	replicas := 1

	if model.Spec.Replicas != nil {
		replicas = int(*model.Spec.Replicas)
	}
	if deployment.Spec.Replicas != nil {
		if int(*deployment.Spec.Replicas) == replicas {
			return false, nil
		}

		deployment.Spec.Replicas = lo.ToPtr(int32(replicas))
	} else {
		deployment.Spec.Replicas = lo.ToPtr(int32(replicas))
	}

	err = c.Update(ctx, deployment)
	if err != nil {
		return false, err
	}

	modelRecorder.Eventf(corev1.EventTypeNormal, "ModelScaled", "Model scaled from %d to %d", deployment.Status.Replicas, replicas)

	return true, nil
}

func NewServiceForModel(namespace, name string, deployment *appsv1.Deployment, serviceType corev1.ServiceType) *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:        ModelServiceName(name),
			Namespace:   namespace,
			Labels:      ModelLabels(name),
			Annotations: ModelAnnotations(ModelAppName(name), false),
			OwnerReferences: []metav1.OwnerReference{{
				APIVersion:         "apps/v1",
				Kind:               "Deployment",
				Name:               deployment.Name,
				UID:                deployment.UID,
				BlockOwnerDeletion: lo.ToPtr(true),
			}},
		},
		Spec: corev1.ServiceSpec{
			Type:     serviceType,
			Selector: ModelLabels(name),
			Ports: []corev1.ServicePort{
				{
					Name:       "ollama",
					Port:       11434,
					TargetPort: intstr.FromInt(11434),
				},
			},
		},
	}
}

func EnsureServiceCreated(
	ctx context.Context,
	namespace string,
	name string,
	deployment *appsv1.Deployment,
) (*corev1.Service, error) {
	c := ClientFromContext(ctx)
	modelRecorder := WrappedRecorderFromContext[*ollamav1.Model](ctx)

	service, err := getServiceByLabels(ctx, c, namespace, ModelLabels(name))
	if err != nil {
		return nil, err
	}
	if service != nil {
		return service, nil
	}

	service = NewServiceForModel(namespace, name, deployment, corev1.ServiceTypeClusterIP)

	err = c.Create(ctx, service)
	if err != nil {
		return nil, err
	}

	modelRecorder.Eventf("Normal", "ServiceCreated", "Service %s created", service.Name)

	return service, nil
}

func IsServiceReady(
	ctx context.Context,
	namespace string,
	name string,
) (bool, error) {
	log := log.FromContext(ctx)
	c := ClientFromContext(ctx)
	modelRecorder := WrappedRecorderFromContext[*ollamav1.Model](ctx)

	service, err := getServiceByLabels(ctx, c, namespace, ModelLabels(name))
	if err != nil {
		return false, err
	}
	if service == nil {
		return false, nil
	}
	if service.Spec.ClusterIP == "" {
		log.Info("waiting for service to have cluster IP", "service", service)
		modelRecorder.Eventf("Normal", "WaitingForService", "Waiting for service %s to have cluster IP", service.Name)

		return false, nil
	}

	log.Info("service is ready", "service", service)

	return true, nil
}

func IsProgressing(ctx context.Context, ollamaModelResource ollamav1.Model) bool {
	return len(lo.Filter(ollamaModelResource.Status.Conditions, func(item ollamav1.ModelStatusCondition, _ int) bool {
		return item.Type == ollamav1.ModelProgressing
	})) > 0
}

func SetProgressing(
	ctx context.Context,
	c client.Client,
	ollamaModelResource ollamav1.Model,
) (bool, error) {
	hasProgressing := len(lo.Filter(ollamaModelResource.Status.Conditions, func(item ollamav1.ModelStatusCondition, _ int) bool {
		return item.Type == ollamav1.ModelProgressing
	})) > 0
	if hasProgressing {
		return false, nil
	}

	ollamaModelResource.Status.Conditions = []ollamav1.ModelStatusCondition{
		{
			Type:               ollamav1.ModelProgressing,
			Status:             corev1.ConditionTrue,
			LastUpdateTime:     metav1.Now(),
			LastTransitionTime: metav1.Now(),
		},
	}

	err := c.Status().Update(ctx, &ollamaModelResource)
	if err != nil {
		return false, err
	}

	return true, nil
}

func IsAvailable(ctx context.Context, ollamaModelResource ollamav1.Model) bool {
	return len(lo.Filter(ollamaModelResource.Status.Conditions, func(item ollamav1.ModelStatusCondition, _ int) bool {
		return item.Type == ollamav1.ModelAvailable
	})) > 0
}

func SetAvailable(
	ctx context.Context,
	ollamaModelResource *ollamav1.Model,
) (bool, error) {
	c := ClientFromContext(ctx)

	hasAvailable := len(lo.Filter(ollamaModelResource.Status.Conditions, func(item ollamav1.ModelStatusCondition, _ int) bool {
		return item.Type == ollamav1.ModelAvailable
	})) > 0
	if hasAvailable {
		return false, nil
	}

	ollamaModelResource.Status.Conditions = []ollamav1.ModelStatusCondition{
		{
			Type:               ollamav1.ModelAvailable,
			Status:             corev1.ConditionTrue,
			LastUpdateTime:     metav1.Now(),
			LastTransitionTime: metav1.Now(),
		},
	}

	err := c.Status().Update(ctx, ollamaModelResource)
	if err != nil {
		return false, err
	}

	return true, nil
}

func ShouldSetReplicas(
	ctx context.Context,
	ollamaModelResource *ollamav1.Model,
	replicas int32,
	readyReplicas int32,
	availableReplicas int32,
	unavailableReplicas int32,
) bool {
	return ollamaModelResource.Status.Replicas != replicas ||
		ollamaModelResource.Status.ReadyReplicas != readyReplicas ||
		ollamaModelResource.Status.AvailableReplicas != availableReplicas ||
		ollamaModelResource.Status.UnavailableReplicas != unavailableReplicas
}

func SetReplicas(
	ctx context.Context,
	ollamaModelResource *ollamav1.Model,
	replicas int32,
	readyReplicas int32,
	availableReplicas int32,
	unavailableReplicas int32,
) (bool, error) {
	c := ClientFromContext(ctx)

	ollamaModelResource.Status.Replicas = replicas
	ollamaModelResource.Status.ReadyReplicas = readyReplicas
	ollamaModelResource.Status.AvailableReplicas = availableReplicas
	ollamaModelResource.Status.UnavailableReplicas = unavailableReplicas

	err := c.Status().Update(ctx, ollamaModelResource)
	if err != nil {
		return false, err
	}

	return true, nil
}
