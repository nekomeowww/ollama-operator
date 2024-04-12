package model

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
)

type WrappedRecorder[T runtime.Object] struct {
	recorder record.EventRecorder
	t        T
}

func NewWrappedRecorder[T runtime.Object](recorder record.EventRecorder, object T) *WrappedRecorder[T] {
	return &WrappedRecorder[T]{
		recorder: recorder,
		t:        object,
	}
}

func (r *WrappedRecorder[T]) Event(eventType, reason, message string) {
	r.recorder.Event(r.t, eventType, reason, message)
}

// Eventf is just like Event, but with Sprintf for the message field.
func (r *WrappedRecorder[T]) Eventf(eventType, reason, messageFmt string, args ...any) {
	r.recorder.Eventf(r.t, eventType, reason, messageFmt, args...)
}

// AnnotatedEventf is just like eventf, but with annotations attached
func (r *WrappedRecorder[T]) AnnotatedEventf(annotations map[string]string, eventType, reason, messageFmt string, args ...any) {
	r.recorder.AnnotatedEventf(r.t, annotations, eventType, reason, messageFmt, args...)
}
