package types

type Component struct {
	Base `json:",inline" yaml:",inline"`
	Spec Values `json:"spec,omitempty" yaml:"spec,omitempty"`
}
