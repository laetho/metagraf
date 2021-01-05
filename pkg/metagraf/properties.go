package metagraf

import "github.com/pkg/errors"

// Returns a metagraf adressable key for a property.
func (mgp *MGProperty) MGKey() string {
	return mgp.Source+"|"+mgp.Key
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

// Return a slice of property keys. If required == true only return required keys.
func (mgp MGProperties) SourceKeys(required bool) []string {
	var keys []string
	for _, prop := range mgp {
		if required {
			if prop.Required == required {
				keys = append(keys, prop.MGKey())
			}
		} else {
			keys = append(keys, prop.MGKey())
		}
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
		if len(env.SecretFrom) > 0 {continue}
		if len(env.EnvFrom) > 0 {continue}
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

func (mgp MGProperties) GetByKey(key string) (MGProperty, error){
	for _, p := range mgp {
		if p.Key == key {
			return p, nil
		}
	}
	return MGProperty{}, errors.Errorf("Key not found!")
}