package types

// Примеры описания Dependency
// 1)
// Name: postgresql
// Создастся база postgresql
// Output:
// host, port,

type Dependency struct {
	Name   string  `json:"name,omitempty" yaml:"name,omitempty"`
	Alias  string  `json:"alias,omitempty" yaml:"alias,omitempty"`
	Script *Script `json:",inline" yaml:",inline"`
}

type Script struct {
	Type       string  `json:"type,omitempty" yaml:"type,omitempty"`
	Version    string  `json:"version,omitempty" yaml:"version,omitempty"`
	Repository string  `json:"repository,omitempty" yaml:"repository,omitempty"`
	Input      *Values `json:"input,omitempty" yaml:"input,omitempty"`
	Output     *Values `json:"output,omitempty" yaml:"output,omitempty"`
}

type Values = map[string]interface{}
