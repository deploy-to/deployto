package types

type Environment struct {
	Base `json:",inline" yaml:",inline"`
	Spec *EnvironmentSpec `json:"spec,omitempty" yaml:"spec,omitempty"`
}

type EnvironmentSpec struct {
	Targets []string `json:"targets,omitempty" yaml:"targets,omitempty"`
	Script  Values   `json:"values,omitempty" yaml:"values,omitempty"`
}
