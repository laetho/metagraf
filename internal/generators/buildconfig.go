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

package generators

import (
	"fmt"
	"strconv"
	"strings"
	"encoding/json"

	"metagraf/internal/metagraf"
	"github.com/blang/semver"

	buildv1 "github.com/openshift/api/build/v1"
)

func GenBuildConfig( mg *metagraf.MetaGraf) {
	sv, err := semver.Parse(mg.Spec.Version)
	if err != nil {
		fmt.Println(err)
	}

	bc := buildv1.BuildConfig{}
	bc.SetName( strings.ToLower(mg.Metadata.Name +"v"+strconv.FormatUint(sv.Major, 10)) )

	btp := buildv1.BuildTriggerPolicy{}
	btp.Type = "GitHub"



	bc.Spec.Triggers = append(bc.Spec.Triggers, btp)


	ba, err := json.Marshal(bc)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(ba))
}
