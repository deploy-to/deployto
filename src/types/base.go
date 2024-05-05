package types

import (
	"deployto/src/filesystem"
)

const DeploytoAPIVersion = "deployto.dev/v1beta1"

type Base struct {
	Kind       string    `json:"kind,omitempty" yaml:"kind,omitempty"`
	APIVersion string    `json:"apiVersion,omitempty" yaml:"apiVersion,omitempty"`
	Meta       *MetaData `json:"metadata,omitempty" yaml:"metadata,omitempty"`
	Status     struct {
		Filesystem *filesystem.Filesystem
		FileName   string
	}
}

type MetaData struct {
	Name        string `json:"name,omitempty" yaml:"name,omitempty"`
	Labels      Labels `json:"labels,omitempty" yaml:"labels,omitempty"`
	Annotations Labels `json:"annotations,omitempty" yaml:"annotations,omitempty"`
}
