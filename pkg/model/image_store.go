package model

import (
	"context"

	"github.com/samber/lo"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	ollamav1 "github.com/nekomeowww/ollama-operator/api/ollama/v1"
)

const (
	ImageStorePVCName         = "ollama-models-store-pvc"
	ImageStoreStatefulSetName = "ollama-models-store"
)

func getImageStorePVC(ctx context.Context, client client.Client, namespace string) (*corev1.PersistentVolumeClaim, error) {
	var pvc corev1.PersistentVolumeClaim

	err := client.Get(ctx, types.NamespacedName{Namespace: namespace, Name: ImageStorePVCName}, &pvc)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, nil
		}

		return nil, err
	}

	return &pvc, nil
}

func EnsureImageStorePVCCreated(
	ctx context.Context,
	client client.Client,
	namespace string,
	storageClassName *string,
	pvcSource *corev1.PersistentVolumeClaimVolumeSource,
	pvSpec *ollamav1.ModelPersistentVolumeSpec,
	modelRecorder *WrappedRecorder[*ollamav1.Model],
) (*corev1.PersistentVolumeClaim, error) {
	log := log.FromContext(ctx)

	pvc, err := getImageStorePVC(ctx, client, namespace)
	if err != nil {
		return nil, err
	}
	if pvc != nil {
		return pvc, nil
	}

	log.Info("no existing image storage PVC found, creating one...")

	accessMode := corev1.ReadWriteOnce
	if pvSpec != nil && pvSpec.AccessMode != nil {
		accessMode = *pvSpec.AccessMode
	}

	pvc = &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:        ImageStorePVCName,
			Namespace:   namespace,
			Labels:      ImageStoreLabels(),
			Annotations: ModelAnnotations(ImageStoreStatefulSetName, true),
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			Resources: corev1.VolumeResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: resource.MustParse("100Gi"),
				},
			},
			StorageClassName: storageClassName,
			AccessModes:      []corev1.PersistentVolumeAccessMode{accessMode},
		},
	}

	err = client.Create(ctx, pvc)
	if err != nil {
		return nil, err
	}

	log.Info("created image storage PVC", "pvc", pvc)
	modelRecorder.Event(corev1.EventTypeNormal, "ProvisionedImageStoragePVC", "Provisioned image storage PVC")

	return pvc, nil
}

func GetImageStorePVByPVC(ctx context.Context, c client.Client, pvc *corev1.PersistentVolumeClaim) (*corev1.PersistentVolume, error) {
	var pv corev1.PersistentVolume

	err := c.Get(ctx, types.NamespacedName{Namespace: pvc.ObjectMeta.Namespace, Name: pvc.Spec.VolumeName}, &pv)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, nil
		}

		return nil, err
	}

	return &pv, nil
}

func getImageStoreStatefulSet(ctx context.Context, client client.Client, namespace string) (*appsv1.StatefulSet, error) {
	var statefulSet appsv1.StatefulSet

	err := client.Get(ctx, types.NamespacedName{Namespace: namespace, Name: ImageStoreStatefulSetName}, &statefulSet)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, nil
		}

		return nil, err
	}

	return &statefulSet, nil
}

func EnsureImageStoreStatefulSetCreated(
	ctx context.Context,
	client client.Client,
	namespace string,
	modelRecorder *WrappedRecorder[*ollamav1.Model],
) (*appsv1.StatefulSet, error) {
	log := log.FromContext(ctx)

	statefulSet, err := getImageStoreStatefulSet(ctx, client, namespace)
	if err != nil {
		return nil, err
	}
	if statefulSet != nil {
		return statefulSet, nil
	}

	log.Info("no existing image store stateful set found, creating one...")

	statefulSet = &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:        ImageStoreStatefulSetName,
			Namespace:   namespace,
			Labels:      ImageStoreLabels(),
			Annotations: ModelAnnotations(ImageStoreStatefulSetName, true),
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas: lo.ToPtr(int32(1)),
			Selector: &metav1.LabelSelector{
				MatchLabels: ImageStoreLabels(),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      ImageStoreLabels(),
					Annotations: ModelAnnotations(ImageStoreStatefulSetName, true),
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						NewOllamaServerContainer(false, corev1.ResourceRequirements{}),
					},
					RestartPolicy: corev1.RestartPolicyAlways,
					Volumes: []corev1.Volume{
						{
							Name: "image-storage",
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: ImageStorePVCName,
									ReadOnly:  false,
								},
							},
						},
					},
				},
			},
		},
	}

	err = client.Create(ctx, statefulSet)
	if err != nil {
		return nil, err
	}

	log.Info("created image store stateful set", "statefulSet", statefulSet)
	modelRecorder.Event(corev1.EventTypeNormal, "ProvisionedImageStoreStatefulSet", "Provisioned image store stateful set")

	return statefulSet, nil
}

func IsImageStoreStatefulSetReady(
	ctx context.Context,
	client client.Client,
	namespace string,
	modelRecorder *WrappedRecorder[*ollamav1.Model],
) (bool, error) {
	log := log.FromContext(ctx)

	statefulSet, err := getImageStoreStatefulSet(ctx, client, namespace)
	if err != nil {
		return false, err
	}
	if statefulSet == nil {
		return false, nil
	}
	if statefulSet.Status.ReadyReplicas == 1 {
		return true, nil
	}

	log.Info("waiting for image store stateful set to be ready", "statefulSet", statefulSet)
	modelRecorder.Event(corev1.EventTypeNormal, "WaitingForImageStoreStatefulSet", "Waiting for image store stateful set to become ready")

	return false, nil
}

func getImageStoreService(ctx context.Context, client client.Client, namespace string) (*corev1.Service, error) {
	var service corev1.Service

	err := client.Get(ctx, types.NamespacedName{Namespace: namespace, Name: ImageStoreStatefulSetName}, &service)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, nil
		}

		return nil, err
	}

	return &service, nil
}

func EnsureImageStoreServiceCreated(
	ctx context.Context,
	client client.Client,
	namespace string,
	statefulSet *appsv1.StatefulSet,
	modelRecorder *WrappedRecorder[*ollamav1.Model],
) (*corev1.Service, error) {
	log := log.FromContext(ctx)

	service, err := getImageStoreService(ctx, client, namespace)
	if err != nil {
		return nil, err
	}
	if service != nil {
		return service, nil
	}

	log.Info("no existing image store service found, creating one...")

	service = &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:        ImageStoreStatefulSetName,
			Namespace:   namespace,
			Labels:      ImageStoreLabels(),
			Annotations: ModelAnnotations(ImageStoreStatefulSetName, true),
			OwnerReferences: []metav1.OwnerReference{{
				APIVersion:         "apps/v1",
				Kind:               "StatefulSet",
				Name:               statefulSet.Name,
				UID:                statefulSet.UID,
				BlockOwnerDeletion: lo.ToPtr(true),
			}},
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeClusterIP,
			Ports: []corev1.ServicePort{
				{
					Name:       "ollama",
					Protocol:   corev1.ProtocolTCP,
					Port:       11434,
					TargetPort: intstr.FromInt(11434),
				},
			},
			Selector: ImageStoreLabels(),
		},
	}

	err = client.Create(ctx, service)
	if err != nil {
		return nil, err
	}

	log.Info("created image store service", "service", service)
	modelRecorder.Event(corev1.EventTypeNormal, "ProvisionedImageStoreService", "Provisioned image store service")

	return service, nil
}

func IsImageStoreServiceReady(
	ctx context.Context,
	client client.Client,
	namespace string,
	modelRecorder *WrappedRecorder[*ollamav1.Model],
) (bool, error) {
	log := log.FromContext(ctx)

	service, err := getImageStoreService(ctx, client, namespace)
	if err != nil {
		return false, err
	}
	if service == nil {
		return false, nil
	}
	if service.Spec.ClusterIP != "" {
		return true, nil
	}

	log.Info("waiting for image store service to be ready", "service", service, "due to no ClusterIP is set")
	modelRecorder.Event(corev1.EventTypeNormal, "WaitingForImageStoreService", "Waiting for image store service to become ready")

	return false, nil
}
