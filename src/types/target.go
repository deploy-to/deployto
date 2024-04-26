package types

type Target struct {
	Base      `json:",inline" yaml:",inline"`
	Namespace string  `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	Script    *Script `json:"script,omitempty" yaml:"script,omitempty"`
}
