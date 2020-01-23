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

// The same as GetVars() but returns a map of only required
// addressable variables.
func (mg *MetaGraf) GetRequiredProperties() MGProperties {
	vars := MGProperties{}

	// Environment Section
	for _,env := range mg.Spec.Environment.Local {
		if env.Required == false { continue }
		p:= MGProperty{
			Source:   "local",
			Key:      env.Name,
			Value:    "",
			Required: env.Required,
		}
		vars = append(vars,p )
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
				p:= MGProperty{
					Source:   conf.Name,
					Key:      opts.Name,
					Value:    "",
					Required: opts.Required,
				}
				vars = append(vars,p)
			}
		case "JVM_SYS_PROP":
			for _, opts := range conf.Options {
				if opts.Required == false {continue}
				p:= MGProperty{
					Source:   "JVM_SYS_PROP",
					Key:      opts.Name,
					Value:    "",
					Required: opts.Required,
				}
				vars = append(vars,p)
			}
		}
	}
	return vars
}


// Returns a MGVars map of addressable variables found in the specification
// where the variable key gets prepended a string of where it came from in
// the specification.
//
// Optional required variables are returned. (Required: true)
func (mg *MetaGraf) GetSourceKeyedVars(defaults bool) MGProperties {
	vars := MGProperties{}
	// Environment Section
	for _,env := range mg.Spec.Environment.Local {
		if env.Required == false {continue}
		if defaults {
			prop := MGProperty{
				Source: "local",
				Key:    env.Name,
				Value:  env.Default,
				Required: env.Required,
			}
			vars = append(vars, prop)
		} else {
			prop := MGProperty{
				Source: "local",
				Key:    env.Name,
				Value:  "",
				Required: env.Required,
			}
			vars = append(vars, prop)
		}
	}

	for _,env := range mg.Spec.Environment.External.Introduces {
		if env.Required == false {continue}
		if defaults {
			prop := MGProperty{
				Source: "external",
				Key:    env.Name,
				Value:  env.Default,
				Required: env.Required,
			}
			vars = append(vars, prop)
		} else {
			prop := MGProperty{
				Source: "external",
				Key:    env.Name,
				Value:  "",
				Required: env.Required,
			}
			vars = append(vars, prop)
		}
	}

	for _,env := range mg.Spec.Environment.External.Consumes {
		if env.Required == false {continue}
		if defaults {
			prop := MGProperty{
				Source: "external",
				Key:    env.Name,
				Value:  env.Default,
				Required: env.Required,
			}
			vars = append(vars, prop)
		} else {
			prop := MGProperty{
				Source: "external",
				Key:    env.Name,
				Value:  "",
				Required: env.Required,
			}
			vars = append(vars, prop)
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
					if opts.Required == false {continue} // Skip optional
					prop := MGProperty{
						Source: conf.Name,
						Key:    opts.Name,
						Value:  "",
						Required: opts.Required,
					}
					vars = append(vars, prop)
				}
			case "JVM_SYS_PROP":
				for _,opts := range conf.Options {
					if opts.Required == false {continue} // Skip optional
					prop := MGProperty{
						Source: conf.Name,
						Key:    opts.Name,
						Value:  opts.Default,
						Required: opts.Required,
					}
					vars = append(vars, prop)
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
