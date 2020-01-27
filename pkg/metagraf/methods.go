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

// Returns a metagraf adressable key for a property.
func (mgp *MGProperty) MGKey() string {
	return mgp.Source+"|"+mgp.Key
}

// Returns a struct (MGProperties) of all MGProperty addressable
// in the metaGraf specification.
func (mg *MetaGraf) GetProperties() MGProperties {
	props := MGProperties{}

	// Config section, find parameters from
	for _,conf := range mg.Spec.Config {
		if len(conf.Options) == 0 {
			continue
		}

		switch conf.Type {
		case "parameters":
			for _, opts := range conf.Options {
				p := MGProperty{
					Source:   conf.Name,
					Key:      opts.Name,
					Value:    "",
					Required: opts.Required,
					Default: opts.Default,
				}
				props[p.MGKey()] = p
			}
		case "JVM_SYS_PROP":
			for _, opts := range conf.Options {
				p := MGProperty{
					Source:   "JVM_SYS_PROP",
					Key:      opts.Name,
					Value:    "",
					Required: opts.Required,
					Default: opts.Default,
				}
				props[p.MGKey()] = p
			}
		}
	}

	// Environment Section
	for _,env := range mg.Spec.Environment.Local {
		p := MGProperty{
			Source:   "local",
			Key:      env.Name,
			Value:    "",
			Required: env.Required,
			Default: env.Default,
		}

		// Environment variables of type JVM_SYS_PROP will
		// be implicitly populated by values from config
		// named JVM_SYS_PROP
		if env.Type == "JVM_SYS_PROP" {
			continue
		}
		props[p.MGKey()] = p
	}
	for _,env := range mg.Spec.Environment.External.Introduces {
		p := MGProperty{
			Source:   "external",
			Key:      env.Name,
			Value:    "",
			Required: env.Required,
			Default: env.Default,
		}
		props[p.MGKey()] = p
	}
	for _,env := range mg.Spec.Environment.External.Consumes {
		p := MGProperty{
			Source:   "external",
			Key:      env.Name,
			Value:    "",
			Required: env.Required,
			Default: env.Default,
		}
		props[p.MGKey()] = p
	}


	return props
}

// Returns a map of all parameterized fields in a metaGraf
// specification.
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

// The same as GetVars() but returns a map of only required
// addressable variables.
func (mg *MetaGraf) GetRequiredVars() MGVars {
	vars := MGVars{}

	// Environment Section
	for _,env := range mg.Spec.Environment.Local {
		if env.Required == false { continue }
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
				if opts.Required == false {continue}
				vars[opts.Name] = ""
			}
		case "JVM_SYS_PROP":
			for _, opts := range conf.Options {
				if opts.Required == false {continue}
				vars[opts.Name] = ""
			}
		}

	}
	return vars
}

// Returns the MGProperty.Required = true
func (mgp MGProperties) GetRequired() MGProperties {
	props := MGProperties{}

	for _, prop := range mgp {
		if prop.Required && prop.Source != "external" {
			props[prop.MGKey()] = prop
		}
	}
	return props
}
// Returns a slice of Keys
func (mgp MGProperties) Keys() []string {
	var keys []string
	for _, prop := range mgp {
		keys = append(keys, prop.Key)
	}
	return keys
}

// Return a slice
func (mgp MGProperties) SourceKeys(required bool) []string {
	var keys []string
	for _, prop := range mgp {
		if prop.Required == required {
			keys = append(keys, prop.MGKey())
			continue
		}
		keys = append(keys, prop.MGKey())
	}
	return keys
}

// Returns a map of key,values
func (mgp MGProperties) KeyMap() map[string]string {
	keys := make(map[string]string)
	for _, prop := range mgp {
		keys[prop.Key] = prop.Value
	}
	return keys
}

// Returns a map of MGProperty.Source+Key, MGProperty.Value
// It takes a boolean as argument. Return only required or all keys?
func (mgp MGProperties) SourceKeyMap(required bool) map[string]string {
	keys := make(map[string]string)
	for _, prop := range mgp {
		if prop.Required == required {
			keys[prop.MGKey()] = prop.Value
			continue
		}
		keys[prop.MGKey()] = prop.Value
	}
	return keys
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
