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

package metagraf

import "github.com/pkg/errors"

// Returns a map of all parameterized fields in a metaGraf
// specification.
// @todo need to look for parameterized fields in more places?
func (mg *MetaGraf) GetVars() MGVars {
	vars := MGVars{}

	// Environment Section
	for _,env := range mg.Spec.Environment.Local {
		vars[env.Name] = ""
	}
	for _,env := range mg.Spec.Environment.External.Introduces {
		vars[env.Name] = ""
	}
	for _,env := range mg.Spec.Environment.External.Consumes {
		vars[env.Name] = ""
	}

	// Config section, find parameters from
	for _,conf := range mg.Spec.Config {
		if len(conf.Options) == 0 {
			continue
		}

		switch conf.Type {
		case "parameters":
			for _, opts := range conf.Options {
				vars[opts.Name] = ""
			}
		case "JVM_SYS_PROP":
			for _, opts := range conf.Options {
				vars[opts.Name] = ""
			}
		}

	}
	return vars
}

func (mg *MetaGraf) GetVarsFromSource(defaults bool) MGVars {
	vars := MGVars{}

	// Environment Section
	for _,env := range mg.Spec.Environment.Local {
		if env.Required == false {continue}
		if defaults {
			vars["local="+env.Name] = env.Default
		} else {
			vars["local="+env.Name] = ""
		}
	}
	for _,env := range mg.Spec.Environment.External.Introduces {
		if env.Required == false {continue}
		if defaults {
			vars["external="+env.Name] = env.Default
		} else {
			vars["external="+env.Name] = ""
		}
	}
	for _,env := range mg.Spec.Environment.External.Consumes {
		if env.Required == false {continue}
		if defaults {
			vars["external="+env.Name] = env.Default
		} else {
			vars["external="+env.Name] = ""
		}
	}

	// Config section, find parameters from
	for _,conf := range mg.Spec.Config {
		if len(conf.Options) == 0 {
			continue
		}

		switch conf.Type {
			case "parameters":
				for _,opts := range conf.Options {
					if opts.Required == false {continue}
					vars[conf.Name+"="+opts.Name] = opts.Default
				}
			case "JVM_SYS_PROP":
				for _,opts := range conf.Options {
					if opts.Required == false {continue}
					vars[conf.Type+"="+opts.Name] = opts.Default
				}
		}
	}

	return vars
}

func (mg *MetaGraf) GetResourceByName(name string) (Resource, error) {
	for _,r := range mg.Spec.Resources{
		if r.Name == name {
			return r, nil
		}
	}
	return Resource{}, errors.New("Resource{} not found, name: "+name)
}

//
func (mg *MetaGraf) GetSecretByName(name string) (Secret, error) {
	for _,s := range mg.Spec.Secret{
		if s.Name == name {
			return s, nil
		}
	}
	return Secret{}, errors.New("Secret{} not found, name: "+name)
}

//
func (mg *MetaGraf) GetConfigByName(name string) (Config, error) {
	for _,c := range mg.Spec.Config{
		if c.Name == name {
			return c, nil
		}
	}
	return Config{}, errors.New("Config{} not found, name: "+name)
}
