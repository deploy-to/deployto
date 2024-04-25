// apiVersion: deployto.dev/v1beta1
// kind: Application || Component   -  spec, для Application и Component одинаковоя. Просто рядом с Application будут лежать envirement и targetа
// metadata:
//   name: <name>
// spec:
//   components:  # массив компонент, которые нужно задеплоить, при деплои приложения (у каждой из них, могут быть свои зависимости)
//   - name:       имя зависимости - например, postgresql, s3, serviceB
//     alias:      опционально, если хочется поменять имя.
//     type:       опционально, значение по умолчанию DeploytoComponent , можно указать helm, в будущем teraform
//     version:    опционально, важно для helm
//     repository: file://./envirements-helm/dev
//     values:
//       - key: value   - если хочется переопределить значение по умолчанию
//   script:
//     type:      # повторяет описание component'ы, за исключением полей name, alias
//     version:
//     repository:
//     values:
//       - key: value

package types

type Component struct {
	Base `json:",inline" yaml:",inline"`
	Spec *ApplicationSpec `json:"spec,omitempty" yaml:"spec,omitempty"`
}

type Application struct {
	Base `json:",inline" yaml:",inline"`
	Spec *ApplicationSpec `json:"spec,omitempty" yaml:"spec,omitempty"`
}

type ApplicationSpec struct {
	Dependencies []*Dependency `json:"dependencies,omitempty" yaml:"dependencies,omitempty"`
	Script       *Script       `json:"script,omitempty" yaml:"script,omitempty"`
}
