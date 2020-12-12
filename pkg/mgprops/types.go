package mgprops

// Map to hold all variables from a specification
type MGVars			map[string]string

// Structure to hold specification section sourced parameters. Should
// solve key collisions and generally be a more workable solution.
type MGProperty struct {
	Source		string	`json:"source"`
	Key			string	`json:"key"`
	Value		string	`json:"value,omitempty"`
	Required	bool	`json:"required,omitempty"`
	Default		string	`json:"default,omitempty"`
}

// map for holding MGProperty structs, should be keyed by
// MGProperty.Source + ":" + MGProperty.Key
type MGProperties map[string]MGProperty
