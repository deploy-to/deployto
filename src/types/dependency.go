package types

type Dependency struct {
	Alias  string `json:"alias,omitempty" yaml:"alias,omitempty"`
	Root   bool   `json:"root,omitempty" yaml:"root,omitempty"`
	Script Values `json:"script,omitempty" yaml:"script,omitempty"`
}
