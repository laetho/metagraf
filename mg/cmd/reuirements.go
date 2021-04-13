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

package cmd

import (
	"github.com/laetho/metagraf/internal/pkg/params"
	"github.com/spf13/viper"
	log "k8s.io/klog"
	"os"
)

// Require atleast one argument that is the metagraf specification.
func requireMetagraf(args []string) {
	if len(args) < 1 {
		log.Error(StrMissingMetaGraf)
		os.Exit(-1)
	}
}

// Require a that a namespace is provided in either config or in
// params.NameSpace variable.
func requireNamespace() {
	if len(params.NameSpace) == 0 {
		params.NameSpace = viper.GetString("namespace")
		if len(params.NameSpace) == 0 {
			log.Error(StrMissingNamespace)
			os.Exit(-1)
		}
	}
}
