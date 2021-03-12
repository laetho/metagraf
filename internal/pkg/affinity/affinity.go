package affinity

import (
	corev1 "k8s.io/api/core/v1"
	"metagraf/pkg/metagraf"
)

func MGAffinityRules(mg *metagraf.Metagraf, topologyKey string, ) corev1.Affinity {

	aff := corev1.Affinity{
		PodAffinity: &corev1.PodAffinity{
			RequiredDuringSchedulingIgnoredDuringExecution:  nil,
			PreferredDuringSchedulingIgnoredDuringExecution: nil,
		},
		PodAntiAffinity: &corev1.PodAntiAffinity{
			RequiredDuringSchedulingIgnoredDuringExecution:  nil,
			PreferredDuringSchedulingIgnoredDuringExecution: nil,
		},
	}
	return aff
}