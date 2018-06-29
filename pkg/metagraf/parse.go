/*
Copyright 2018 The MetaGraf Authors

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
	"io/ioutil"
)

func Parse(filepath string) MetaGraf {
	b, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic(err)
	}

	var mg MetaGraf

	json.Unmarshal(b, &mg)
	if err != nil {
		panic(err)
	}

	//fmt.Printf("%T", mg)
	//fmt.Println(mg)
	return mg
}
