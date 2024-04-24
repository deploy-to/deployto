package types

type Target struct {
	Base   `json:",inline" yaml:",inline"`
	Script *Script `json:"script,omitempty" yaml:"script,omitempty"`
}
