package types

type Job struct {
	Base `json:",inline" yaml:",inline"`
	Spec *JobSpec `json:"spec,omitempty" yaml:"spec,omitempty"`
}

type JobSpec struct {
	Steps []Step `json:"steps,omitempty" yaml:"steps,omitempty"`
}

type Step struct {
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
	Run  string `json:"run,omitempty" yaml:"run,omitempty"`
}
