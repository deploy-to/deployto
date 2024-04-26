package types

type Dependency struct {
	Kind   string  `json:"kind,omitempty" yaml:"kind,omitempty"`
	Name   string  `json:"name,omitempty" yaml:"name,omitempty"`
	Root   bool    `json:"root,omitempty" yaml:"root,omitempty"`
	Script *Script `json:",inline" yaml:",inline"`
}

// TODO подумать, как будет выглядить Script для:
// * Type == helm
// * Type == terraform
// * Type == component
// * Type == do-not-deploy
type Script struct {
	Type       string  `json:"type,omitempty" yaml:"type,omitempty"`
	Version    string  `json:"version,omitempty" yaml:"version,omitempty"`
	Repository string  `json:"repository,omitempty" yaml:"repository,omitempty"`
	Input      *Values `json:"input,omitempty" yaml:"input,omitempty"`
	Output     *Values `json:"output,omitempty" yaml:"output,omitempty"`
}

type Values = map[string]interface{}
