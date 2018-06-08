package generators

import (
	"metagraf/internal/metagraf"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"fmt"
	"encoding/json"
	"strconv"
	"github.com/blang/semver"
)

type ConfigMap struct {
	metav1.TypeMeta			`json:",inline"`
	metav1.ObjectMeta		`json:"metadata"`
	Data map[string]string	`json:"data,omitempty"`
}


func GenConfigMaps( mg *metagraf.MetaGraf) {

	/*

	We need to create the following ConfigMaps:

		* INTPL.Config.properties
		* jvm.params
		* server.xml

	*/

	genJvmParams(mg)


}

func genJvmParams( mg *metagraf.MetaGraf ) {

	// Parse version with semver library
	sv, err := semver.Parse(mg.Spec.Version)
	if err != nil {
		fmt.Println(err)
	}

	cm := ConfigMap{}
	// Need to initialize the map defined in struct
	cm.Data = make(map[string]string)

	cm.Name = mg.Metadata.Name + "v" + strconv.FormatUint(sv.Major, 10) + "-jvm.params"
	for _, e  := range mg.Spec.Environment.Local {
		if !e.Required {
			cm.Data["blah"] = "placeholder"
		}
	}


	b, err := json.Marshal(cm)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(b))

}