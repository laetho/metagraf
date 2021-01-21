package metagraf

import (
	corev1 "k8s.io/api/core/v1"
)

func (e *EnvironmentVar) ToEnvVar() (corev1.EnvVar) {
	ev := corev1.EnvVar{}
	switch e.GetType() {
	case "secretfrom":
		if len(e.Key) == 0 {
			e.Key = e.SecretFrom
		}
		ev = corev1.EnvVar{
			Name: e.Name,
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: e.SecretFrom,
					},
					Key: e.Key,
				},
			},
		}
	case "envfrom":
		ev = corev1.EnvVar{
			Name: e.Name,
			ValueFrom: &corev1.EnvVarSource{
				ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: e.EnvFrom,
					},
					Key: e.Key,
				},
			},
		}
		return ev
	case "default":
		ev = corev1.EnvVar{
			Name:      e.Name,
			Value:     "",
		}
	}
	return ev
}

func (e *EnvironmentVar) GetType() string {
	if len(e.SecretFrom) > 0 {
		return "secretfrom"
	} else if len(e.EnvFrom) > 0 {
		return "envfrom"
	} else if e.Type == "JVM_SYS_PROP" {
		return "JVM_SYS_PROP"
	} else {
		return "default"
	}
}
