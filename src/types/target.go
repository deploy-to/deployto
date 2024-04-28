package types

type Target struct {
	Base   `json:",inline" yaml:",inline"`
	Script Values `json:"values,omitempty" yaml:"values,omitempty"`
}
