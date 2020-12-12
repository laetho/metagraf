package mgprops

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
