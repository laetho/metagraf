package generators

import (
	"metagraf/internal/metagraf"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"fmt"
	"encoding/json"
	"strconv"
	"github.com/blang/semver"
	"strings"
)

type ConfigMap struct {
	metav1.TypeMeta			`json:",inline"`
	metav1.ObjectMeta		`json:"metadata"`
	Data map[string]string	`json:"data,omitempty"`
}

/*
Entry function for creating a slew of configmaps, this will be very
specific to NT internal workings for now.
 */
func GenConfigMaps( mg *metagraf.MetaGraf) {

	/*

	We need to create the following ConfigMaps:

		* INTPL.Config.properties
		* jvm.params
		* server.xml

	*/

	for _, c := range mg.Spec.Config {
		if c.Type != "parameters" {
			continue
		}
		genConfigMapFromConfig(&c, mg)
	}


}

/*
Generates a configmap for jvm.params file for Liberty java apps
 */
func genConfigMapFromConfig( conf *metagraf.Config, mg *metagraf.MetaGraf) {

	// Parse version with semver library
	sv, err := semver.Parse(mg.Spec.Version)
	if err != nil {
		fmt.Println(err)
	}

	cm := ConfigMap{}
	// Need to initialize the map defined in struct
	// Should be done in a factory maybe.
	cm.Data = make(map[string]string)
	cm.TypeMeta.Kind = "ConfigMap"
	cm.TypeMeta.APIVersion = "v1"

	cm.Name = strings.ToLower(mg.Metadata.Name+ "v" + strconv.FormatUint(sv.Major, 10) +"-"+ conf.Name)
	for _, o  := range conf.Options {
		cm.Data[o.Name] = o.Default
	}

	b, err := json.Marshal(cm)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(b))
}