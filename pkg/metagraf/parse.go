/*
Copyright 2019 The metaGraf Authors

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

package metagraf

import (
	"encoding/json"
	log "k8s.io/klog"
	"io/ioutil"
	"os"
)

func Parse(filepath string) MetaGraf {
	b, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	err = nil
	var mg MetaGraf

	err = json.Unmarshal(b, &mg)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	return mg
}
