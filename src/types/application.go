// apiVersion: deployto.dev/v1beta1
// kind: Application
// metadata:
//   name: <name>
// spec:
//   components:
//   - name: <component name>
//     repository: file://./envirements-helm/dev
//   envirements:
//     - <envirement name>
//   helm:
//     repository: file://./envirements-helm/dev

package types

type Application struct {
	Meta `json:",inline" yaml:",inline"`
	Spec *ApplicationSpec `json:"spec,omitempty" yaml:"spec,omitempty"`
}

type ApplicationSpec struct {
	Components []*Components `json:"components,omitempty" yaml:"components,omitempty"`
	Helm       *Helm         `json:"helm,omitempty" yaml:"helm,omitempty"`
}

type Components struct {
	Name       string `json:"name,omitempty" yaml:"name,omitempty"`
	Repository string `json:"repository,omitempty" yaml:"repository,omitempty"`
}
