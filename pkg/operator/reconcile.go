package operator

import (
	"context"
	"errors"
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/samber/lo"
)

func ResultFromError(err error) (*ctrl.Result, error) {
	if err == nil {
		return &ctrl.Result{}, nil
	}

	var requeueErr *RequeueError
	if errors.As(err, &requeueErr) {
		return lo.ToPtr(requeueErr.Result()), requeueErr.err
	}

	return nil, err
}

func HandleError(ctx context.Context, result *ctrl.Result, err error) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	if result != nil && err != nil {
		log.Error(err, "Requeue", "after", result.RequeueAfter)
		return *result, err
	}
	if result != nil {
		return *result, nil
	}

	return ctrl.Result{}, err
}

type SubReconciler[T runtime.Object] struct {
	version   string
	group     string
	kind      string
	reconcile func(ctx context.Context, ns string, name string, m T) error
}

type SubReconcilers[T runtime.Object] []SubReconciler[T]

func (s SubReconcilers[T]) Reconcile(ctx context.Context, req ctrl.Request, m T) error {
	log := log.FromContext(ctx)

	for _, subReconciler := range s {
		log.Info("Reconciling", "kind", subReconciler.kind)

		err := subReconciler.reconcile(ctx, req.Namespace, req.Name, m)
		if err != nil {
			return err
		}
	}

	return nil
}

func NewSubReconcilers[T runtime.Object](reconcilers ...SubReconciler[T]) SubReconcilers[T] {
	return reconcilers
}

type ReconcileHandler[T runtime.Object] func(ctx context.Context, ns string, name string, m T) error

func NewPVCReconciler[T runtime.Object](fn ReconcileHandler[T]) SubReconciler[T] {
	return SubReconciler[T]{
		version:   "v1",
		group:     "core",
		kind:      "PersistentVolumeClaim",
		reconcile: fn,
	}
}

func NewStatefulSetReconciler[T runtime.Object](fn ReconcileHandler[T]) SubReconciler[T] {
	return SubReconciler[T]{
		version:   "v1",
		group:     "apps",
		kind:      "StatefulSet",
		reconcile: fn,
	}
}

func NewServiceReconciler[T runtime.Object](fn ReconcileHandler[T]) SubReconciler[T] {
	return SubReconciler[T]{
		version:   "v1",
		group:     "core",
		kind:      "Service",
		reconcile: fn,
	}
}

func NewDeploymentReconciler[T runtime.Object](fn ReconcileHandler[T]) SubReconciler[T] {
	return SubReconciler[T]{
		version:   "v1",
		group:     "apps",
		kind:      "Deployment",
		reconcile: fn,
	}
}

func RequeueAfter(d time.Duration) error {
	return &RequeueError{after: d}
}

func RequeueAfterWithError(d time.Duration, err error) error {
	return &RequeueError{after: d, err: err}
}

type RequeueError struct {
	after time.Duration
	err   error
}

func (r *RequeueError) Error() string {
	if r.err != nil {
		return fmt.Sprintf("encountered an error: %v, requeue after %d", r.err, r.after.Milliseconds())
	}

	return fmt.Sprintf("requeue after %d", r.after.Milliseconds())
}

func (r *RequeueError) Result() reconcile.Result {
	return reconcile.Result{Requeue: true, RequeueAfter: r.after}
}
