/*
Copyright 2018 The metaGraf Authors

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

package helpers

import (
	"strings"
)

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// Sanitize a string for a k8s label field
func LabelString(s string) string {
	if len(s) > 63 {
		s = s[0:63]
	}
	return s
}

// Checks if a strings in a slice is part of a string.
func SliceInString(list []string, str string) bool {
	for _, filter := range list {
		if strings.Contains(str, filter) {
			return true
		}
	}
	return false
}

func PathToIdentifier(path string) string {
	return strings.Replace(path, "/", "-", -1)
}
