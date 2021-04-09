/*
Copyright 2021 The metaGraf Authors

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

package pdb

import (
	"context"
	"github.com/golang/glog"
	k8sclient2 "github.com/laetho/metagraf/internal/pkg/k8sclient"
	params2 "github.com/laetho/metagraf/internal/pkg/params"
	"github.com/laetho/metagraf/pkg/metagraf"
	"github.com/laetho/metagraf/pkg/modules"
	"k8s.io/api/policy/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/apimachinery/pkg/util/intstr"
	log "k8s.io/klog"
	"k8s.io/kubectl/pkg/scheme"
	"math"
	"os"
	"time"
)

func GenPodDisruptionBudget(mg *metagraf.MetaGraf, replicas int32) v1beta1.PodDisruptionBudget {
	name := modules.Name(mg) // @todo refactor how we create a name.

	var maxunavail int32

	switch {
	case replicas < 2:
		maxunavail = 0
	case replicas >= 2:
		maxunavail = int32(math.Floor(float64(replicas / 2)))
	}

	l := make(map[string]string)
	l["app"] = name

	selector := metav1.LabelSelector{
		MatchLabels: l,
	}

	obj := v1beta1.PodDisruptionBudget{
		TypeMeta: metav1.TypeMeta{
			Kind:       "PodDisruptionBudget",
			APIVersion: "v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			CreationTimestamp: metav1.Time{
				Time: time.Now(),
			},
		},
		Spec: v1beta1.PodDisruptionBudgetSpec{
			MaxUnavailable: &intstr.IntOrString{
				Type:   0,
				IntVal: maxunavail,
			},
			Selector: &selector,
		},
	}

	if !params2.Dryrun {
		StorePodDisruptionBudget(obj)
	}
	if params2.Output {
		MarshalObject(obj.DeepCopyObject())
	}
	return obj
}

func StorePodDisruptionBudget(obj v1beta1.PodDisruptionBudget) {

	glog.Infof("ResourceVersion: %v Length: %v", obj.ResourceVersion, len(obj.ResourceVersion))
	glog.Infof("Namespace: %v", params2.NameSpace)

	client := k8sclient2.GetKubernetesClient().PolicyV1beta1().PodDisruptionBudgets(params2.NameSpace)
	if len(obj.ResourceVersion) > 0 {
		// update
		result, err := client.Update(context.TODO(), &obj, metav1.UpdateOptions{})
		if err != nil {
			glog.Info(err)
		}
		glog.Infof("Updated: %v(%v)", result.Name, obj.Name)
	} else {
		result, err := client.Create(context.TODO(), &obj, metav1.CreateOptions{})
		if err != nil {
			glog.Info(err)
		}
		glog.Infof("Created: %v(%v)", result.Name, obj.Name)
	}
}

// todo: need to restructure code, this is a duplication
// Marshal kubernetes resource to json
func MarshalObject(obj runtime.Object) {
	switch params2.Format {
	case "json":
		opt := json.SerializerOptions{
			Yaml:   false,
			Pretty: true,
			Strict: true,
		}
		s := json.NewSerializerWithOptions(json.DefaultMetaFactory, scheme.Scheme, scheme.Scheme, opt)
		err := s.Encode(obj, os.Stdout)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
	case "yaml":
		opt := json.SerializerOptions{
			Yaml:   true,
			Pretty: true,
			Strict: true,
		}
		s := json.NewSerializerWithOptions(json.DefaultMetaFactory, scheme.Scheme, scheme.Scheme, opt)
		err := s.Encode(obj, os.Stdout)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
	}
}
