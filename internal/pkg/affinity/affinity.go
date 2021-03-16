package affinity

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"metagraf/internal/pkg/params/params"
)


func AntiAffinityLabelSelector(key string, op metav1.LabelSelectorOperator, value string) *metav1.LabelSelector {

	var me []metav1.LabelSelectorRequirement
	var values []string

	values = append(values, value)

	me = append(me, metav1.LabelSelectorRequirement{
		Key:      key,
		Operator: metav1.LabelSelectorOpIn,
		Values:   values,
	})

	return &metav1.LabelSelector{
		MatchExpressions: me,
	}
}

func SoftPodAntiAffinity(app string, topologyKey string, weight int32) *corev1.Affinity {

	var terms  []corev1.WeightedPodAffinityTerm
	var namespaces []string

	namespaces = append(namespaces, params.NameSpace)

	terms = append(terms, corev1.WeightedPodAffinityTerm{
		Weight:          weight,
		PodAffinityTerm: corev1.PodAffinityTerm{
			LabelSelector: AntiAffinityLabelSelector("app", metav1.LabelSelectorOpIn, app),
			Namespaces:    namespaces,
			TopologyKey:   topologyKey,
		},
	})

	aff := corev1.Affinity{
		PodAntiAffinity: &corev1.PodAntiAffinity{
			PreferredDuringSchedulingIgnoredDuringExecution: terms,
		},
	}
	return &aff
}

func HardPodAntiAffinity(app string, topologyKey string) *corev1.Affinity {
	var terms  []corev1.PodAffinityTerm
	var namespaces []string

	namespaces = append(namespaces, params.NameSpace)

	terms = append(terms, corev1.PodAffinityTerm{
			LabelSelector: AntiAffinityLabelSelector("app", metav1.LabelSelectorOpIn, app),
			Namespaces:    namespaces,
			TopologyKey:   topologyKey,
		},
	)

	aff := corev1.Affinity{
		PodAntiAffinity: &corev1.PodAntiAffinity{
			RequiredDuringSchedulingIgnoredDuringExecution: terms,

		},
	}
	return &aff
}