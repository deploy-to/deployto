package types

type Dependency struct {
	Kind   string  `json:"kind,omitempty" yaml:"kind,omitempty"`
	Name   string  `json:"name,omitempty" yaml:"name,omitempty"`
	Root   bool    `json:"root,omitempty" yaml:"root,omitempty"`
	Values *Values `json:"values,omitempty" yaml:"values,omitempty"`
}
