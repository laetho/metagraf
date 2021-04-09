package oam

import (
	"fmt"
	oamv1 "github.com/crossplane/oam-kubernetes-runtime/apis/core/v1alpha2"
	params2 "github.com/laetho/metagraf/internal/pkg/params"
	"github.com/laetho/metagraf/pkg/metagraf"
	"github.com/laetho/metagraf/pkg/modules"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func GetParameters(mg *metagraf.MetaGraf) []oamv1.ComponentParameter {
	props := mg.GetProperties()
	cpars := []oamv1.ComponentParameter{}

	for _, v := range props {
		if v.Source != "local" {
			continue
		}
		param := oamv1.ComponentParameter{
			Name:     v.Key,
			Required: &v.Required,
		}
		cpars = append(cpars, param)
	}
	return cpars
}

func GetEnvs(envs *[]metagraf.EnvironmentVar) {

}

func GenOAMComponent(mg *metagraf.MetaGraf) {
	fmt.Println(mg)

	obj := oamv1.Component{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Component",
			APIVersion: "core.oam.dev/v1alpha2",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:                       mg.Metadata.Name,
			Namespace:                  params2.NameSpace,
			DeletionGracePeriodSeconds: nil,
			Labels:                     nil,
			Annotations:                nil,
			OwnerReferences:            nil,
			Finalizers:                 nil,
		},
		Spec: oamv1.ComponentSpec{
			Workload: runtime.RawExtension{
				Object: &oamv1.ContainerizedWorkload{
					TypeMeta:   metav1.TypeMeta{},
					ObjectMeta: metav1.ObjectMeta{},
					Spec:       oamv1.ContainerizedWorkloadSpec{},
					Status:     oamv1.ContainerizedWorkloadStatus{},
				},
			},
			Parameters: GetParameters(mg),
		},
		Status: oamv1.ComponentStatus{},
	}

	if modules.Output {
		modules.MarshalObject(obj.DeepCopyObject())
	}

}

func GenOAMApplicationConfiguration(mg *metagraf.MetaGraf) {
	fmt.Println(mg)

	obj := oamv1.ApplicationConfiguration{
		TypeMeta:   metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{},
		Spec:       oamv1.ApplicationConfigurationSpec{},
		Status:     oamv1.ApplicationConfigurationStatus{},
	}

	if modules.Output {
		modules.MarshalObject(obj.DeepCopyObject())
	}
}
