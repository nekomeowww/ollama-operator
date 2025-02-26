package kollama

import (
	"context"
	"fmt"
	"time"

	"github.com/gookit/color"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/briandowns/spinner"
	"github.com/nekomeowww/ollama-operator/pkg/model"
)

const (
	isReadyString = " is ready"
)

func waitUntilImageStoreReady(ctx context.Context, kubeClient client.Client, namespace string) error {
	var imageStore appsv1.StatefulSet

	err := kubeClient.Get(ctx, types.NamespacedName{Namespace: namespace, Name: model.ImageStoreStatefulSetName}, &imageStore)
	if err != nil {
		if apierrors.IsNotFound(err) {
			time.Sleep(1 * time.Second)
			return waitUntilImageStoreReady(ctx, kubeClient, namespace)
		}

		return err
	}

	return nil
}

func waitUntilImageStoreServiceReady(ctx context.Context, kubeClient client.Client, namespace string) error {
	var imageStoreService corev1.Service

	err := kubeClient.Get(ctx, types.NamespacedName{Namespace: namespace, Name: model.ImageStoreStatefulSetName}, &imageStoreService)
	if err != nil {
		if apierrors.IsNotFound(err) {
			time.Sleep(1 * time.Second)
			return waitUntilImageStoreServiceReady(ctx, kubeClient, namespace)
		}

		return err
	}

	return nil
}

func waitUntilOllamaModelDeploymentPullImageDone(ctx context.Context, kubeClient client.Client, namespace string, name string) error {
	var pods corev1.PodList

	err := kubeClient.List(ctx, &pods, client.InNamespace(namespace), client.MatchingLabels{"app": model.ModelAppName(name)})
	if err != nil {
		return err
	}

	for _, pod := range pods.Items {
		if len(pod.Status.InitContainerStatuses) == 0 || !pod.Status.InitContainerStatuses[0].Ready {
			time.Sleep(1 * time.Second)
			return waitUntilOllamaModelDeploymentPullImageDone(ctx, kubeClient, namespace, name)
		}
		if len(pod.Status.ContainerStatuses) == 0 || !pod.Status.ContainerStatuses[0].Ready {
			time.Sleep(1 * time.Second)
			return waitUntilOllamaModelDeploymentPullImageDone(ctx, kubeClient, namespace, name)
		}
	}

	return nil
}

func waitUntilOllamaModelDeploymentReady(ctx context.Context, kubeClient client.Client, namespace string, name string) error {
	var deployment appsv1.Deployment

	err := kubeClient.Get(ctx, types.NamespacedName{Namespace: namespace, Name: model.ModelAppName(name)}, &deployment)
	if err != nil {
		if apierrors.IsNotFound(err) {
			time.Sleep(1 * time.Second)
			return waitUntilOllamaModelDeploymentReady(ctx, kubeClient, namespace, name)
		}

		return err
	}

	if deployment.Status.ReadyReplicas == 0 {
		time.Sleep(1 * time.Second)
		return waitUntilOllamaModelDeploymentReady(ctx, kubeClient, namespace, name)
	}

	return nil
}

func waitUntilOllamaModelServiceReady(ctx context.Context, kubeClient client.Client, namespace string, name string) error { //nolint:unparam
	var services corev1.ServiceList

	err := kubeClient.List(ctx, &services, client.MatchingLabels(model.ModelLabels(name)))
	if err != nil {
		return err
	}
	if len(services.Items) == 0 {
		time.Sleep(1 * time.Second)
		return waitUntilOllamaModelServiceReady(ctx, kubeClient, namespace, name)
	}

	return nil
}

func waitUntilModelAvailable(kubeClient client.Client, namespace string, modelName string, modelImage string) error {
	s := spinner.New(spinner.CharSets[14], 200*time.Millisecond)
	s.FinalMSG = color.FgGreen.Render("✓") + " image store" + isReadyString
	_ = s.Color("blue")

	s.Start()
	s.Suffix = " preparing image store..."

	waitConditionCtx, waitConditionCancel := context.WithTimeout(context.Background(), 1*time.Hour)
	defer waitConditionCancel()

	err := waitUntilImageStoreReady(waitConditionCtx, kubeClient, namespace)
	if err != nil {
		return err
	}

	s.Stop()
	fmt.Println()

	s = spinner.New(spinner.CharSets[14], 200*time.Millisecond)
	s.FinalMSG = color.FgGreen.Render("✓") + " image store exposed"
	_ = s.Color("blue")

	s.Start()
	s.Suffix = " exposing image store service..."

	err = waitUntilImageStoreServiceReady(waitConditionCtx, kubeClient, namespace)
	if err != nil {
		return err
	}

	s.Stop()
	fmt.Println()

	s = spinner.New(spinner.CharSets[14], 200*time.Millisecond)
	s.FinalMSG = color.FgGreen.Render("✓") + " model pulled and prepared"
	_ = s.Color("blue")

	s.Start()
	s.Suffix = " pulling model image \"" + modelImage + "\"..."

	err = waitUntilOllamaModelDeploymentPullImageDone(waitConditionCtx, kubeClient, namespace, modelName)
	if err != nil {
		return err
	}

	s.Stop()
	fmt.Println()

	s = spinner.New(spinner.CharSets[14], 200*time.Millisecond)
	s.FinalMSG = color.FgGreen.Render("✓") + " model" + isReadyString
	_ = s.Color("blue")

	s.Start()
	s.Suffix = " deploying model..."

	err = waitUntilOllamaModelDeploymentReady(waitConditionCtx, kubeClient, namespace, modelName)
	if err != nil {
		return err
	}

	s.Stop()
	fmt.Println()

	return nil
}
