package kollama

import (
	"context"
	"fmt"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"

	ollamav1 "github.com/nekomeowww/ollama-operator/api/ollama/v1"
)

var (
	ErrOllamaModelNotSupported = fmt.Errorf("%s is not supported on the cluster, did you install the Ollama Operator?", modelSchemaGroupVersionResource.String())
)

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
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(obj.UnstructuredContent(), typedObj); err != nil {
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
