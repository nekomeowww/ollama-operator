/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	ollamav1 "github.com/nekomeowww/ollama-operator/api/ollama/v1"
	model "github.com/nekomeowww/ollama-operator/pkg/model"
)

// ModelReconciler reconciles a Model object
type ModelReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

//+kubebuilder:rbac:groups=ollama.ayaka.io,resources=models,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=ollama.ayaka.io,resources=models/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=ollama.ayaka.io,resources=models/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps,resources=statefulsets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=persistentvolumeclaims,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=storageclasses,verbs=get;list;watch
//+kubebuilder:rbac:groups=core,resources=persistentvolumes,verbs=get;list;watch
//+kubebuilder:rbac:groups=core,resources=namespaces,verbs=get;list;watch
//+kubebuilder:rbac:groups=batch,resources=jobs,verbs=get;list;watch;create;update;patch;delete;deletecollection
//+kubebuilder:rbac:groups=batch,resources=jobs/status,verbs=get

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.17.2/pkg/reconcile
func (r *ModelReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var m ollamav1.Model

	err := r.Get(ctx, req.NamespacedName, &m)
	if err != nil {
		return reconcile.Result{}, client.IgnoreNotFound(err)
	}

	modelRecorder := model.NewWrappedRecorder(r.Recorder, &m)

	if !model.IsAvailable(ctx, m) {
		hasSet, err := model.SetProgressing(ctx, r.Client, m)
		if err != nil {
			return reconcile.Result{}, err
		}
		if hasSet {
			modelRecorder.Eventf("Normal", "ModelProgressing", "Model is progressing")
			return reconcile.Result{Requeue: true, RequeueAfter: time.Second * 1}, nil
		}
	}

	modelStorageClass := m.Spec.StorageClassName
	modelPVC := m.Spec.PersistentVolumeClaim
	modelPV := m.Spec.PersistentVolume

	_, err = model.EnsureImageStorePVCCreated(ctx, r.Client, req.Namespace, modelStorageClass, modelPVC, modelPV, modelRecorder)
	if err != nil {
		return reconcile.Result{}, err
	}

	statefulSet, err := model.EnsureImageStoreStatefulSetCreated(ctx, r.Client, req.Namespace, modelRecorder)
	if err != nil {
		return reconcile.Result{}, err
	}

	statefulSetReady, err := model.IsImageStoreStatefulSetReady(ctx, r.Client, req.Namespace, modelRecorder)
	if err != nil {
		return reconcile.Result{}, err
	}
	if !statefulSetReady {
		return reconcile.Result{Requeue: true, RequeueAfter: time.Second * 5}, nil
	}

	_, err = model.EnsureImageStoreServiceCreated(ctx, r.Client, req.Namespace, statefulSet, modelRecorder)
	if err != nil {
		return reconcile.Result{}, err
	}

	serviceReady, err := model.IsImageStoreServiceReady(ctx, r.Client, req.Namespace, modelRecorder)
	if err != nil {
		return reconcile.Result{}, err
	}
	if !serviceReady {
		return reconcile.Result{Requeue: true, RequeueAfter: time.Second * 5}, nil
	}

	deployment, err := model.EnsureDeploymentCreated(ctx, r.Client, req.Namespace, req.Name, m.Spec.Image, m.Spec.Replicas, &m, modelRecorder)
	if err != nil {
		return reconcile.Result{}, err
	}

	modelDeploymentUpdated, err := model.UpdateDeployment(ctx, r.Client, &m, modelRecorder)
	if err != nil {
		return reconcile.Result{}, err
	}
	if modelDeploymentUpdated {
		return reconcile.Result{Requeue: true, RequeueAfter: time.Second * 5}, nil
	}

	modelDeploymentReady, err := model.IsDeploymentReady(ctx, r.Client, req.Namespace, req.Name, modelRecorder)
	if err != nil {
		return reconcile.Result{}, err
	}
	if !modelDeploymentReady {
		return reconcile.Result{Requeue: true, RequeueAfter: time.Second * 5}, nil
	}

	_, err = model.EnsureServiceCreated(ctx, r.Client, req.Namespace, req.Name, deployment, modelRecorder)
	if err != nil {
		return reconcile.Result{}, err
	}

	modelServiceReady, err := model.IsServiceReady(ctx, r.Client, req.Namespace, req.Name, modelRecorder)
	if err != nil {
		return reconcile.Result{}, err
	}
	if !modelServiceReady {
		return reconcile.Result{Requeue: true, RequeueAfter: time.Second * 5}, nil
	}

	if model.ShouldSetReplicas(ctx, m, deployment.Status.Replicas, deployment.Status.ReadyReplicas, deployment.Status.AvailableReplicas, deployment.Status.UnavailableReplicas) {
		hasSet, err := model.SetReplicas(ctx, r.Client, m, deployment.Status.Replicas, deployment.Status.ReadyReplicas, deployment.Status.AvailableReplicas, deployment.Status.UnavailableReplicas)
		if err != nil {
			return reconcile.Result{}, err
		}
		if hasSet {
			return reconcile.Result{Requeue: true, RequeueAfter: time.Second * 5}, nil
		}
	}

	_, err = model.SetAvailable(ctx, r.Client, m)
	if err != nil {
		return reconcile.Result{}, err
	}

	modelRecorder.Eventf("Normal", "ModelAvailable", "Model is available")

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ModelReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&ollamav1.Model{}).
		Complete(r)
}
