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

func FindOllamaServerContainer(container corev1.Container) bool {
	return container.Name == "server"
}

func UniqEnvVar(env []corev1.EnvVar) []corev1.EnvVar {
	return lo.UniqBy(env, func(item corev1.EnvVar) string {
		return item.Name
	})
}

func AssignOllamaServerContainer(readOnly bool, resources corev1.ResourceRequirements, extraEnvFrom []corev1.EnvFromSource, extraEnv []corev1.EnvVar) func(container corev1.Container, _ int) corev1.Container {
	return func(container corev1.Container, _ int) corev1.Container {
		container.Image = OllamaBaseImage
		container.Args = []string{
			"serve",
		}

		container.EnvFrom = append(container.EnvFrom, extraEnvFrom...)
		_, configuredOllamaPortAsEnvFrom := lo.Find(container.EnvFrom, func(item corev1.EnvFromSource) bool {
			return item.Prefix == "OLLAMA_PORT"
		})

		container.Env = append(container.Env, extraEnv...)
		container.Env = AppendIfNotFound(container.Env, func(item corev1.EnvVar) bool {
			return item.Name == "OLLAMA_HOST"
		}, func() corev1.EnvVar {
			return corev1.EnvVar{
				Name:  "OLLAMA_HOST",
				Value: "0.0.0.0",
			}
		})

		_, configuredOllamaPort := lo.Find(container.Env, func(item corev1.EnvVar) bool {
			return item.Name == "OLLAMA_PORT"
		})

		if !configuredOllamaPort && !configuredOllamaPortAsEnvFrom {
			container.Ports = AppendIfNotFound(container.Ports, func(item corev1.ContainerPort) bool {
				return item.ContainerPort == 11434
			}, func() corev1.ContainerPort {
				return corev1.ContainerPort{
					Name:          "ollama",
					Protocol:      corev1.ProtocolTCP,
					ContainerPort: 11434,
				}
			})
		}

		container.VolumeMounts = AppendIfNotFound(container.VolumeMounts, func(item corev1.VolumeMount) bool {
			return item.Name == "image-storage"
		}, func() corev1.VolumeMount {
			return corev1.VolumeMount{
				Name:      "image-storage",
				MountPath: "/root/.ollama",
				ReadOnly:  readOnly,
			}
		})

		container.Resources.Limits = lo.Assign(container.Resources.Limits, resources.Limits)
		container.Resources.Requests = lo.Assign(container.Resources.Requests, resources.Requests)

		container.ReadinessProbe = &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				HTTPGet: &corev1.HTTPGetAction{
					Path: "/api/tags",
					Port: intstr.FromString("ollama"),
				},
			},
			InitialDelaySeconds: 5,
			SuccessThreshold:    1,
			FailureThreshold:    2500,
			TimeoutSeconds:      5,
		}

		container.LivenessProbe = &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				HTTPGet: &corev1.HTTPGetAction{
					Path: "/api/tags",
					Port: intstr.FromString("ollama"),
				},
			},
			InitialDelaySeconds: 5,
			SuccessThreshold:    1,
			FailureThreshold:    2500,
			TimeoutSeconds:      1,
		}

		return container
	}
}

func NewOllamaServerContainer(readOnly bool, resources corev1.ResourceRequirements, extraEnvFrom []corev1.EnvFromSource, extraEnv []corev1.EnvVar) corev1.Container {
	return corev1.Container{
		Name:  "server",
		Image: OllamaBaseImage,
		Args: []string{
			"serve",
		},
		Env: UniqEnvVar(
			append(
				append([]corev1.EnvVar{}, corev1.EnvVar{
					Name:  "OLLAMA_HOST",
					Value: "0.0.0.0",
				}),
				extraEnv...,
			),
		),
		EnvFrom: extraEnvFrom,
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
			InitialDelaySeconds: 5,
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
			InitialDelaySeconds: 5,
			SuccessThreshold:    1,
			FailureThreshold:    2500,
			TimeoutSeconds:      1,
		},
	}
}

func FindOllamaPullerContainer(container corev1.Container) bool {
	return container.Name == "ollama-image-pull"
}

func AssignOllamaPullerContainer(name string, image string, serverLocatedNamespace string, resources corev1.ResourceRequirements, extraEnvFrom []corev1.EnvFromSource, extraEnv []corev1.EnvVar) func(container corev1.Container, _ int) corev1.Container {
	return func(container corev1.Container, _ int) corev1.Container {
		container.Command = []string{
			"bash",
		}

		container.Args = []string{
			"-c",
			// TODO: This is a temporary solution, we need to find a better way to preload the models
			fmt.Sprintf("apt update && apt install curl -y && ollama pull %s && curl http://ollama-models-store:11434/api/generate -d '{\"model\": \"%s\"}'", image, name),
		}

		container.Env = AppendIfNotFound(container.Env, func(item corev1.EnvVar) bool {
			return item.Name == "OLLAMA_HOST"
		}, func() corev1.EnvVar {
			return corev1.EnvVar{
				Name:  "OLLAMA_HOST",
				Value: fmt.Sprintf("ollama-models-store.%s", serverLocatedNamespace),
			}
		})

		return container
	}
}

func NewOllamaPullerContainer(name string, image string, serverLocatedNamespace string, resources corev1.ResourceRequirements, extraEnvFrom []corev1.EnvFromSource, extraEnv []corev1.EnvVar) corev1.Container {
	return corev1.Container{
		Name:  "ollama-image-pull",
		Image: OllamaBaseImage,
		Command: []string{
			"bash",
		},
		Args: []string{
			"-c",
			// TODO: This is a temporary solution, we need to find a better way to preload the models
			fmt.Sprintf("apt update && apt install curl -y && ollama pull %s && curl http://ollama-models-store:11434/api/generate -d '{\"model\": \"%s\"}'", image, name),
		},
		Env: []corev1.EnvVar{
			{
				Name:  "OLLAMA_HOST",
				Value: fmt.Sprintf("ollama-models-store.%s", serverLocatedNamespace),
			},
		},
	}
}
