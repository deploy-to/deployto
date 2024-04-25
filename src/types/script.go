package types

type Script struct {
	Type       string  `json:"type,omitempty" yaml:"type,omitempty"`
	Version    string  `json:"version,omitempty" yaml:"version,omitempty"`
	Repository string  `json:"repository,omitempty" yaml:"repository,omitempty"`
	Input      *Values `json:"input,omitempty" yaml:"input,omitempty"`
	Output     *Values `json:"output,omitempty" yaml:"output,omitempty"`
}

type Values = map[string]interface{}
