package types

import (
	"deployto/src/filesystem"

	"github.com/fatih/structs"
)

const DeploytoAPIVersion = "deployto.dev/v1beta1"

type Base struct {
	Kind       string     `json:"kind,omitempty"       yaml:"kind,omitempty"       structs:"kind,omitempty"`
	APIVersion string     `json:"apiVersion,omitempty" yaml:"apiVersion,omitempty" structs:"apiVersion,omitempty" `
	Meta       MetaData   `json:"metadata,omitempty"   yaml:"metadata,omitempty"   structs:"metadata,omitempty"   mapstructure:"metadata"`
	Status     StatusType `json:"-"                    yaml:"-"                    structs:",omitempty"`
}

type MetaData struct {
	Name        string `json:"name,omitempty"        yaml:"name,omitempty"        structs:"name,omitempty"`
	Labels      Labels `json:"labels,omitempty"      yaml:"labels,omitempty"      structs:"labels,omitempty"`
	Annotations Labels `json:"annotations,omitempty" yaml:"annotations,omitempty" structs:"annotations,omitempty"`
}

type StatusType struct {
	Filesystem *filesystem.Filesystem `json:"-" yaml:"-" structs:",omitempty"`
	FileName   string
}

func (b *Base) AsValues() Values {
	return structs.Map(b)
}
