// apiVersion: deployto.dev/v1beta1
// kind: Envirement
// meta:
//   name: dev
// spec:
//   targets:
//   - dev-target
//   workflow:
//     type: target-by-target
//   helm:
//     repository: file://./envirements-helm/dev

package types

type Environment struct {
	Meta `json:",inline" yaml:",inline"`
	Spec *EnvironmentSpec `json:"spec,omitempty" yaml:"spec,omitempty"`
}

type EnvironmentSpec struct {
	Targets []string `json:"targets,omitempty" yaml:"targets,omitempty"`
	Helm    *Helm    `json:"helm,omitempty" yaml:"helm,omitempty"`
}
