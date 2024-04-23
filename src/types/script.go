package types

type Script struct {
	Type       string            `json:"type,omitempty" yaml:"type,omitempty"`
	Version    string            `json:"version,omitempty" yaml:"version,omitempty"`
	Repository string            `json:"repository,omitempty" yaml:"repository,omitempty"`
	Values     map[string]string `json:"values,omitempty" yaml:"values,omitempty"`
}
