// apiVersion: deployto.dev/v1beta1
// kind: Envirement
// meta:
//   name: dev
// spec:
//   targets:
//   - dev-target
//   workflow:
//     type: target-by-target
//   script:
//     repository: file://./envirements-helm/dev

package types

type Environment struct {
	Base `json:",inline" yaml:",inline"`
	Spec *EnvironmentSpec `json:"spec,omitempty" yaml:"spec,omitempty"`
}

type EnvironmentSpec struct {
	Targets []string `json:"targets,omitempty" yaml:"targets,omitempty"`
	Script  Values   `json:"values,omitempty" yaml:"values,omitempty"`
}
