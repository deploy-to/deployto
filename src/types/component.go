package types

type Component struct {
	Base `json:",inline" yaml:",inline"`
	Spec *ComponentSpec `json:"spec,omitempty" yaml:"spec,omitempty"`
}

type ComponentSpec struct {
	BuildJob     string        `json:"buildjob,omitempty" yaml:"buildjob,omitempty"`
	Dependencies []*Dependency `json:"dependencies,omitempty" yaml:"dependencies,omitempty"`
	Script       *Script       `json:"script,omitempty" yaml:"script,omitempty"`
}
