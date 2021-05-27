package metagraf

import (
	corev1 "k8s.io/api/core/v1"
)

func (s Secret) VolumeName() string {
	return "vol-"+s.Name
}

// Creates a slice of corev1.Volume{} definitions from buildsecrets defined
// in the metaGraf specification.
func (mg MetaGraf) BuildSecretsToVolumes() ([]corev1.Volume) {
	var vols []corev1.Volume

	for _, s := range mg.Spec.BuildSecret {
		v := corev1.Volume{
			Name:         s.VolumeName(),
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName:  s.Name,
				},
			},
		}
		vols = append(vols,v)
	}
	return vols
}

// Creates a slice of corev1.VolumeMount{} definitions from buildsecrets defined
// in the metaGraf specification.
func (mg MetaGraf) BuildSecretsToVolumeMounts() ([]corev1.VolumeMount) {
	var vols []corev1.VolumeMount

	for _, s := range mg.Spec.BuildSecret {
		v := corev1.VolumeMount{
			Name:             s.VolumeName(),
			ReadOnly:         true,
		}
		if len(s.MountPath) > 0 {
			v.MountPath = s.MountPath
		}
		vols = append(vols,v)
	}
	return vols
}

// Creates a slice of corev1.Volume{} definitions from secrets defined
// in the metaGraf specification.
func (mg MetaGraf) SecretsToVolumes() ([]corev1.Volume) {
	var vols []corev1.Volume

	for _, s := range mg.Spec.Secret {
		v := corev1.Volume{
			Name:         s.VolumeName(),
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName:  s.Name,
				},
			},
		}
		vols = append(vols,v)
	}

	return vols
}
// Creates a slice of corev1.VolumeMount{} definitions from secrets defined
// in the metaGraf specification.
func (mg MetaGraf) SecretsToVolumeMounts() ([]corev1.VolumeMount) {
	var vols []corev1.VolumeMount

	for _, s := range mg.Spec.Secret {
		v := corev1.VolumeMount{
			Name:             s.VolumeName(),
			ReadOnly:         true,
		}
		if len(s.MountPath) > 0 {
			v.MountPath = s.MountPath
		}
		vols = append(vols,v)
	}

	return vols
}
