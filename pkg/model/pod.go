package model

import (
	"fmt"

	"github.com/samber/lo"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const (
	OllamaBaseImage = "ollama/ollama"
)

func NewOllamaServerContainer(readOnly bool, resources corev1.ResourceRequirements) corev1.Container {
	return corev1.Container{
		Name:  "server",
		Image: OllamaBaseImage,
		Args: []string{
			"serve",
		},
		Env: []corev1.EnvVar{
			{
				Name:  "OLLAMA_HOST",
				Value: "0.0.0.0",
			},
		},
		Ports: []corev1.ContainerPort{
			{
				Name:          "ollama",
				Protocol:      corev1.ProtocolTCP,
				ContainerPort: 11434,
			},
		},
		VolumeMounts: []corev1.VolumeMount{
			{
				Name:      "image-storage",
				MountPath: "/root/.ollama",
				ReadOnly:  readOnly,
			},
		},
		Resources: corev1.ResourceRequirements{
			Limits:   lo.Ternary(len(resources.Limits) == 0, nil, resources.Limits),
			Requests: lo.Ternary(len(resources.Requests) == 0, nil, resources.Requests),
		},
		ReadinessProbe: &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				HTTPGet: &corev1.HTTPGetAction{
					Path: "/api/tags",
					Port: intstr.FromString("ollama"),
				},
			},
			InitialDelaySeconds: 15,
			SuccessThreshold:    1,
			FailureThreshold:    2500,
			TimeoutSeconds:      5,
		},
		LivenessProbe: &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				HTTPGet: &corev1.HTTPGetAction{
					Path: "/api/tags",
					Port: intstr.FromString("ollama"),
				},
			},
			InitialDelaySeconds: 15,
			SuccessThreshold:    1,
			FailureThreshold:    2500,
			TimeoutSeconds:      1,
		},
	}
}

func NewOllamaPullerContainer(image string, serverLocatedNamespace string, resources corev1.ResourceRequirements) corev1.Container {
	return corev1.Container{
		Name:  "ollama-image-pull",
		Image: OllamaBaseImage,
		Args: []string{
			"pull",
			image,
		},
		Env: []corev1.EnvVar{
			{
				Name:  "OLLAMA_HOST",
				Value: fmt.Sprintf("ollama-models-store.%s", serverLocatedNamespace),
			},
		},
	}
}
