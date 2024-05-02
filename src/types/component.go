package types

type Component struct {
	Base `json:",inline" yaml:",inline"`
	Spec Values `json:"spec,omitempty" yaml:"spec,omitempty"`
}

// type ComponentSpec struct {
// 	BuildJob     string        `json:"buildjob,omitempty" yaml:"buildjob,omitempty"`
// 	Dependencies []*Dependency `json:"dependencies,omitempty" yaml:"dependencies,omitempty"`
// 	Values       Values        `json:"values,omitempty" yaml:"values,omitempty"`
// }
