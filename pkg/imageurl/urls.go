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

package imageurl

import (
	"net/url"
	"fmt"
	"strings"
)

type ImageURL struct {
	URL *url.URL
	Namespace string
	Image string
	Tag string
}

func (u *ImageURL) Parse( url string) error {

	// Assume https when schema is unspecified
	if !strings.HasPrefix(url, "https") {
		url = "https://" + url
	}

	tmpurl, err := u.URL.Parse(url)
	if err != nil {
		fmt.Println("Error", err)
		return err
	}
	u.URL = tmpurl

	err = u.ParsePath()
	if err != nil {
		return err
	}
	return nil
}

func (u *ImageURL) ParsePath() error {
	sa := strings.Split(u.URL.Path, "/")
	if len(sa) < 3 {
		err := fmt.Errorf("Unexpected path format: %v", u.URL.Path )
		return err
	}
	sb := strings.Split(sa[2], ":")
	u.Namespace = sa[1]
	u.Image = sb[0]

	// Assume latest when no tag is provided
	if len(sb) < 2 {
		u.Tag = "latest"
	} else {
		u.Tag = sb[1]
	}
	return nil
}

func (u *ImageURL) IsValid() bool {
	if len(u.Namespace) < 1 {return false}
	if len(u.Image) < 1 {return false}
	if len(u.Tag) < 1 {return false}
	return true
}