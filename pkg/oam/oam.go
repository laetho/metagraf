package oam

import (
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"metagraf/pkg/metagraf"
	oamv1 "github.com/oam-dev/oam-go-sdk/apis/core.oam.dev/v1alpha1"
	"metagraf/pkg/modules"
)

func GenOAMComponent(mg *metagraf.MetaGraf) {
	fmt.Println(mg)


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
