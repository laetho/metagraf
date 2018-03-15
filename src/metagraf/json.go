package metagraf

// JSON structure for a MetaGraf entity
type MetaGraf struct {
	Kind     string
	Version  string
	Metadata struct {
		Name              string
		Version           string
		ResourceVersion   string
		Namespace         string
		CreationTimestamp string
		Description       string
		Labels            map[string]string
		Annotations       map[string]string
	}
	Spec struct {
		Resources   []Resource
		Environment struct {
			Local    []EnvironmentVar
			External []EnvironmentVar
		}
	}
	Config []ConfigParam
}

type Resource struct {
	Name     string
	Type     string
	Version  string
	Match    string
	Required bool
}
type ConfigParam struct {
	Name        string
	Required    bool
	Description string
	Type        string
	Default     string
}
type EnvironmentVar struct {
	Name        string
	Required    bool
	Type        string
	Description string
}
