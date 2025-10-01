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
	"github.com/nekomeowww/ollama-operator/pkg/operator"
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

	ctx = model.WithWrappedRecorder(ctx, model.NewWrappedRecorder(r.Recorder, &m))
	ctx = model.WithClient(ctx, r.Client)

	res, err := operator.ResultFromError(r.reconcile(ctx, req, &m))

	return operator.HandleError(ctx, res, err)
}

// SetupWithManager sets up the controller with the Manager.
func (r *ModelReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&ollamav1.Model{}).
		Complete(r)
}

func (r *ModelReconciler) reconcile(ctx context.Context, req ctrl.Request, m *ollamav1.Model) error {
	client := model.ClientFromContext(ctx)
	recorder := model.WrappedRecorderFromContext[*ollamav1.Model](ctx)

	if !model.IsAvailable(ctx, *m) {
		hasSet, err := model.SetProgressing(ctx, client, *m)
		if err != nil {
			return err
		}

		if hasSet {
			recorder.Eventf("Normal", "ModelProgressing", "Model is progressing")
			return operator.RequeueAfter(time.Second)
		}
	}

	return operator.NewSubReconcilers(
		operator.NewPVCReconciler(r.reconcilePVC),
		operator.NewStatefulSetReconciler(r.reconcileStatefulSet),
		operator.NewServiceReconciler(r.reconcileService),
		operator.NewDeploymentReconciler(r.reconcileDeployment),
		operator.NewServiceReconciler(r.reconcileModelService),
	).Reconcile(ctx, req, m)
}

func (r *ModelReconciler) reconcilePVC(ctx context.Context, ns string, name string, m *ollamav1.Model) error {
	modelStorageClass := m.Spec.StorageClassName
	modelPVC := m.Spec.PersistentVolumeClaim
	modelPV := m.Spec.PersistentVolume

	_, err := model.EnsureImageStorePVCCreated(ctx, ns, modelStorageClass, modelPVC, modelPV)
	if err != nil {
		return err
	}

	return nil
}

func (r *ModelReconciler) reconcileStatefulSet(ctx context.Context, ns string, name string, m *ollamav1.Model) error {
	_, err := model.EnsureImageStoreStatefulSetCreated(ctx, ns, m)
	if err != nil {
		return err
	}

	statefulSetReady, err := model.IsImageStoreStatefulSetReady(ctx, ns)
	if err != nil {
		return err
	}

	if !statefulSetReady {
		return operator.RequeueAfter(time.Second * 5)
	}

	return nil
}

func (r *ModelReconciler) reconcileService(ctx context.Context, ns string, name string, m *ollamav1.Model) error {
	statefulSet, err := model.EnsureImageStoreStatefulSetCreated(ctx, ns, m)
	if err != nil {
		return err
	}

	_, err = model.EnsureImageStoreServiceCreated(ctx, ns, statefulSet)
	if err != nil {
		return err
	}

	serviceReady, err := model.IsImageStoreServiceReady(ctx, ns)
	if err != nil {
		return err
	}

	if !serviceReady {
		return operator.RequeueAfter(time.Second * 5)
	}

	return nil
}

func (r *ModelReconciler) reconcileDeployment(ctx context.Context, ns string, name string, m *ollamav1.Model) error {
	_, err := model.EnsureDeploymentCreated(ctx, ns, name, m.Spec.Image, m.Spec.Replicas, m)
	if err != nil {
		return err
	}

	modelDeploymentUpdated, err := model.UpdateDeployment(ctx, m)
	if err != nil {
		return err
	}

	if modelDeploymentUpdated {
		return operator.RequeueAfter(time.Second * 5)
	}

	modelDeploymentReady, err := model.IsDeploymentReady(ctx, ns, name)
	if err != nil {
		return err
	}

	if !modelDeploymentReady {
		return operator.RequeueAfter(time.Second * 5)
	}

	return nil
}

func (r *ModelReconciler) reconcileModelService(ctx context.Context, ns string, name string, m *ollamav1.Model) error {
	recorder := model.WrappedRecorderFromContext[*ollamav1.Model](ctx)

	deployment, err := model.EnsureDeploymentCreated(ctx, ns, name, m.Spec.Image, m.Spec.Replicas, m)
	if err != nil {
		return err
	}

	_, err = model.EnsureServiceCreated(ctx, ns, name, deployment)
	if err != nil {
		return err
	}

	modelServiceReady, err := model.IsServiceReady(ctx, ns, name)
	if err != nil {
		return err
	}

	if !modelServiceReady {
		return operator.RequeueAfter(time.Second * 5)
	}

	if model.ShouldSetReplicas(ctx, m, deployment.Status.Replicas, deployment.Status.ReadyReplicas, deployment.Status.AvailableReplicas, deployment.Status.UnavailableReplicas) {
		hasSet, err := model.SetReplicas(ctx, m, deployment.Status.Replicas, deployment.Status.ReadyReplicas, deployment.Status.AvailableReplicas, deployment.Status.UnavailableReplicas)
		if err != nil {
			return err
		}

		if hasSet {
			return operator.RequeueAfter(time.Second * 5)
		}
	}

	_, err = model.SetAvailable(ctx, m)
	if err != nil {
		return err
	}

	recorder.Eventf("Normal", "ModelAvailable", "Model is available")

	return nil
}
