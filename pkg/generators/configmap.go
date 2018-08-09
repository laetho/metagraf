package generators

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"github.com/blang/semver"
	"metagraf/pkg/metagraf"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

/*
Entry function for creating a slew of configmaps, this will be somewhat
specific to NT internal workings for now.
*/
func GenConfigMaps(mg *metagraf.MetaGraf) {

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
func genConfigMapFromConfig(conf *metagraf.Config, mg *metagraf.MetaGraf) {

	// Parse version with semver library
	sv, err := semver.Parse(mg.Spec.Version)
	if err != nil {
		fmt.Println(err)
	}
	objname := strings.ToLower(mg.Metadata.Name + "v" + strconv.FormatUint(sv.Major, 10))

	l := make(map[string]string)
	l["app"] = objname

	cm := corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind: "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: objname,
			Labels: l,
		},
	}

	// Need to initialize the map defined in struct
	// Should be done in a factory maybe.
	cm.Data = make(map[string]string)
	//cm.TypeMeta.Kind = "ConfigMap"
	//cm.TypeMeta.APIVersion = "v1"
	cm.ObjectMeta.Labels = l

	cm.Name = strings.ToLower(mg.Metadata.Name + "v" + strconv.FormatUint(sv.Major, 10) + "-" + strings.Replace(conf.Name, ".", "-", -1 ))
	for _, o := range conf.Options {
		cm.Data[o.Name] = o.Default
	}

	b, err := json.Marshal(cm)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(b))
}
