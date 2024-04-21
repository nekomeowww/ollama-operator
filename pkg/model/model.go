package model

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/samber/lo"

	ollamav1 "github.com/nekomeowww/ollama-operator/api/ollama/v1"
)

func ModelAppName(name string) string {
	return fmt.Sprintf("ollama-model-%s", name)
}

func ModelLabels(appName, name string, imageStore bool) map[string]string {
	return map[string]string{
		"app":                        name,
		"model.ollama.ayaka.io":      name,
		"model.ollama.ayaka.io/type": lo.Ternary(imageStore, "image-store", "model"),
	}
}

func ModelAnnotations(name string, imageStore bool) map[string]string {
	return map[string]string{
		"model.ollama.ayaka.io/name": name,
		"model.ollama.ayaka.io/type": lo.Ternary(imageStore, "image-store", "model"),
	}
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
	c client.Client,
	namespace string,
	name string,
	image string,
	replicas *int32,
	model *ollamav1.Model,
	modelRecorder *WrappedRecorder[*ollamav1.Model],
) (*appsv1.Deployment, error) {
	deployment, err := getDeployment(ctx, c, namespace, name)
	if err != nil {
		return nil, err
	}
	if deployment != nil {
		return deployment, nil
	}

	deployment = &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Labels:      ModelLabels(ModelAppName(name), name, false),
			Annotations: ModelAnnotations(ModelAppName(name), false),
			Name:        ModelAppName(name),
			Namespace:   namespace,
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
				MatchLabels: ModelLabels(ModelAppName(name), name, false),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      ModelLabels(ModelAppName(name), name, false),
					Annotations: ModelAnnotations(ModelAppName(name), false),
				},
				Spec: corev1.PodSpec{
					InitContainers: []corev1.Container{
						NewOllamaPullerContainer(image, namespace),
					},
					Containers: []corev1.Container{
						NewOllamaServerContainer(true),
					},
					Volumes: []corev1.Volume{
						{
							Name: "image-storage",
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: imageStorePVCName,
									ReadOnly:  true,
								},
							},
						},
					},
				},
			},
		},
	}

	err = c.Create(ctx, deployment)
	if err != nil {
		return nil, err
	}

	modelRecorder.Eventf("Normal", "DeploymentCreated", "Deployment %s created", deployment.Name)

	return deployment, nil
}

func IsDeploymentReady(
	ctx context.Context,
	c client.Client,
	namespace string,
	name string,
	modelRecorder *WrappedRecorder[*ollamav1.Model],
) (bool, error) {
	log := log.FromContext(ctx)

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
	c client.Client,
	model *ollamav1.Model,
	modelRecorder *WrappedRecorder[*ollamav1.Model],
) (bool, error) {
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

func getService(ctx context.Context, c client.Client, namespace string, name string) (*corev1.Service, error) {
	var service corev1.Service

	err := c.Get(ctx, client.ObjectKey{Namespace: namespace, Name: ModelAppName(name)}, &service)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, nil
		}

		return nil, err
	}

	return &service, nil
}

func NewServiceForModel(namespace, name string, deployment *appsv1.Deployment, serviceType corev1.ServiceType) *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Labels:      ModelLabels(ModelAppName(name), name, false),
			Annotations: ModelAnnotations(ModelAppName(name), false),
			Name:        ModelAppName(name),
			Namespace:   namespace,
			OwnerReferences: []metav1.OwnerReference{{
				APIVersion:         "apps/v1",
				Kind:               "Deployment",
				Name:               deployment.Name,
				UID:                deployment.UID,
				BlockOwnerDeletion: lo.ToPtr(true),
			}},
		},
		Spec: corev1.ServiceSpec{
			Type: serviceType,
			Selector: map[string]string{
				"app": ModelAppName(name),
			},
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
	c client.Client,
	namespace string,
	name string,
	deployment *appsv1.Deployment,
	modelRecorder *WrappedRecorder[*ollamav1.Model],
) (*corev1.Service, error) {
	service, err := getService(ctx, c, namespace, name)
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
	c client.Client,
	namespace string,
	name string,
	modelRecorder *WrappedRecorder[*ollamav1.Model],
) (bool, error) {
	log := log.FromContext(ctx)

	service, err := getService(ctx, c, namespace, name)
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
	c client.Client,
	ollamaModelResource ollamav1.Model,
) (bool, error) {
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

	err := c.Status().Update(ctx, &ollamaModelResource)
	if err != nil {
		return false, err
	}

	return true, nil
}

func ShouldSetReplicas(
	ctx context.Context,
	ollamaModelResource ollamav1.Model,
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
	c client.Client,
	ollamaModelResource ollamav1.Model,
	replicas int32,
	readyReplicas int32,
	availableReplicas int32,
	unavailableReplicas int32,
) (bool, error) {
	ollamaModelResource.Status.Replicas = replicas
	ollamaModelResource.Status.ReadyReplicas = readyReplicas
	ollamaModelResource.Status.AvailableReplicas = availableReplicas
	ollamaModelResource.Status.UnavailableReplicas = unavailableReplicas

	err := c.Status().Update(ctx, &ollamaModelResource)
	if err != nil {
		return false, err
	}

	return true, nil
}
