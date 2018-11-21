/*
Copyright 2018 The MetaGraph Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package modules

import (
	"bytes"
	"github.com/golang/glog"
	"metagraf/mg/ocpclient"
	"metagraf/pkg/metagraf"
	"os"
	"strings"
	"text/template"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)


/*
	This function will inspect the metaGraf specification
	for which configmaps it will need to create and their
	type in a map.
*/
func FindConfigMaps(mg *metagraf.MetaGraf) map[string]string {
	maps := make(map[string]string)

	for _, c := range mg.Spec.Config {
		maps[strings.ToLower(c.Name)] = "config"
	}

	for _, r := range mg.Spec.Resources {
		if r.Type == "http" {
			continue
		}
		maps[strings.ToLower(r.User)] = "resource"
	}


	return maps
}

/*
Entry function for creating a slew of configmaps, this will be somewhat
specific to NT internal workings for now.
*/
func GenConfigMaps(mg *metagraf.MetaGraf) {
	for _, c := range mg.Spec.Config {
		if c.Type != "parameters" {
			continue
		}
		genConfigMapsFromConfig(&c, mg)
	}
	genConfigMapsFromResources(mg)
}

/*
Generates a configmap for jvm.params file for Liberty java apps
*/
func genConfigMapsFromConfig(conf *metagraf.Config, mg *metagraf.MetaGraf) {

	objname := Name(mg)

	l := make(map[string]string)
	l["app"] = objname

	cm := corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   objname,
			Labels: l,
		},
	}

	// Need to initialize the map defined in struct
	// Should be done in a factory maybe.
	cm.Data = make(map[string]string)
	cm.ObjectMeta.Labels = l

	cm.Name = objname + "-" + strings.Replace(strings.ToLower(conf.Name), ".", "-", -1)
	for _, o := range conf.Options {
		if ValueFromEnv(o.Name) {
			cm.Data[o.Name] = Variables[o.Name]
		} else {
			cm.Data[o.Name] = o.Default
		}
	}

	if !Dryrun {
		StoreConfigMap(cm)
	}
	if Output {
		MarshalObject(cm)
	}

}

func genConfigMapsFromResources(mg *metagraf.MetaGraf) {

	objname := Name(mg)

	for _, r := range mg.Spec.Resources {
		if strings.Contains(r.Type, "oracle") {
			cm := genJDBCOracle(objname, &r)
			if !Dryrun {
				StoreConfigMap(cm)
			}
			if Output {
				MarshalObject(cm)
			}
		}
	}

}

func genJDBCOracle(objname string, r *metagraf.Resource) corev1.ConfigMap {

	l := make(map[string]string)
	l["app"] = objname

	// todo: should this be fetched from EnvironmentVar Template, possibly
	dstemplate := `<dataSource id="{{ .User }}" jndiName="jdbc/{{ .User }}">
<jdbcDriver libraryRef="OracleLib" />
<properties.oracle URL="jdbc:oracle:thin:@(DESCRIPTION=(ADDRESS=(PROTOCOL=TCP)(HOST=$OVIS_HOST)(PORT=$OVIS_PORT))(CONNECT_DATA=(SERVER=DEDICATED)(SERVICE_NAME=OVIS_ACTIVE)))" user="dbadmin" password="$PASSWORD"/>
<connectionManager maxPoolSize="10" minPoolSize="2" />
</dataSource>
`
	t, _ := template.New("ds").Parse(dstemplate)
	var o bytes.Buffer
	if err := t.Execute(&o, r); err != nil {
		glog.Errorf("%v", err)
		os.Exit(1)
	}

	cm := corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   objname + "-" + strings.ToLower(r.User),
			Labels: l,
		},
	}

	cm.Data = make(map[string]string)
	cm.ObjectMeta.Labels = l
	cm.Data["DS"] = o.String()

	// @todo this should really come from secretRef/vault
	cm.Data["PASSWORD"] = "$PASSWORD"
	return cm
}

func StoreConfigMap(m corev1.ConfigMap) {

	glog.Infof("ResourceVersion: %v Length: %v", m.ResourceVersion, len(m.ResourceVersion))

	cmclient := ocpclient.GetCoreClient().ConfigMaps(NameSpace)

	if len(m.ResourceVersion) > 0 {
		// update
		result, err := cmclient.Update(&m)
		if err != nil {
			glog.Error(err)
			//os.Exit(1)
		}
		glog.Infof("Updated configmap: %v(%v)", result.Name, m.Name)
	} else {
		result, err := cmclient.Create(&m)
		if err != nil {
			glog.Error(err)
			//os.Exit(1)
		}
		glog.Infof("Created configmap: %v(%v)", result.Name, m.Name)
	}
}
