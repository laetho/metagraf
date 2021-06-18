package metagraf

import v1 "k8s.io/api/core/v1"

func (mg MetaGraf) Volumes() []v1.Volume {
	var vols []v1.Volume

	for _, v := range mg.Spec.Volume {

		vol := v1.Volume{
			Name:         v.Name,
			VolumeSource: v1.VolumeSource{
				HostPath: hostPathVolumeSource(v),
			},
		}
		vols = append(vols, vol)
	}

	return vols
}

// Returns a HostPathVolumeSoruce given a Volume struct.
func hostPathVolumeSource(volume Volume) *v1.HostPathVolumeSource {
	return &v1.HostPathVolumeSource{
		Path: volume.HostPath.Path,
		Type: volume.HostPath.Type,
	}
}

// Retruns a slice of v1.VolumeMount{} for Volumes defined in metaGraf specification.
func(mg MetaGraf) VolumesToVolumeMounts() []v1.VolumeMount {
	var mounts []v1.VolumeMount

	for _, v := range mg.Spec.Volume {
		mnt := v1.VolumeMount{
			Name:             v.Name,
			ReadOnly:         false,
			MountPath:        v.MountPath,
			MountPropagation: nil,
		}
		mounts = append(mounts, mnt)
	}

	return mounts
}