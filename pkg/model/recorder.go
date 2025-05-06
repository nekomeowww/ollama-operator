package model

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
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

type baseWrapperRecorderContextKey string

func NewWrappedRecorderContextKey(key string) baseWrapperRecorderContextKey {
	return baseWrapperRecorderContextKey(key)
}

const (
	defaultBaseWrapperRecorderContextKey baseWrapperRecorderContextKey = "default"
)

func WithWrappedRecorder[T runtime.Object](ctx context.Context, recorder *WrappedRecorder[T], key ...baseWrapperRecorderContextKey) context.Context {
	if len(key) == 0 {
		return context.WithValue(ctx, defaultBaseWrapperRecorderContextKey, recorder)
	}

	return context.WithValue(ctx, key[0], recorder)
}

func WrappedRecorderFromContext[T runtime.Object](ctx context.Context, key ...baseWrapperRecorderContextKey) *WrappedRecorder[T] {
	if len(key) == 0 {
		r, _ := ctx.Value(defaultBaseWrapperRecorderContextKey).(*WrappedRecorder[T])
		return r
	}

	r, _ := ctx.Value(key[0]).(*WrappedRecorder[T])

	return r
}

type baseClientContextKey string

const (
	defaultBaseClientContextKey baseClientContextKey = "default"
)

func NewClientContextKey(key string) baseClientContextKey {
	return baseClientContextKey(key)
}

func WithClient(ctx context.Context, client client.Client, key ...baseClientContextKey) context.Context {
	if len(key) == 0 {
		return context.WithValue(ctx, defaultBaseClientContextKey, client)
	}

	return context.WithValue(ctx, key[0], client)
}

func ClientFromContext(ctx context.Context, key ...baseClientContextKey) client.Client {
	if len(key) == 0 {
		c, _ := ctx.Value(defaultBaseClientContextKey).(client.Client)
		return c
	}

	c, _ := ctx.Value(key[0]).(client.Client)

	return c
}
