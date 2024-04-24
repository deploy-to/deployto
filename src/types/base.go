package types

import (
	"path/filepath"

	"github.com/rs/zerolog/log"
)

const DeploytoAPIVersion = "deployto.dev/v1beta1"

type Base struct {
	Kind       string    `json:"kind,omitempty" yaml:"kind,omitempty"`
	APIVersion string    `json:"apiVersion,omitempty" yaml:"apiVersion,omitempty"`
	Meta       *MetaData `json:"metadata,omitempty" yaml:"metadata,omitempty"`
	Status     struct {
		FileName string
	}
}

const DeploytoPath = ".deployto" //TODO bug this const defined more then in one place

func (b *Base) StatusGetPath() (path string) {
	path = b.Status.FileName
	for {
		if filepath.Base(path) == DeploytoPath {
			return filepath.Dir(path)
		}
		if len(path) < 4 /*TODO test on windows*/ {
			return ""
		}
		path = filepath.Dir(b.Status.FileName)
	}
}

type MetaData struct {
	Name        string            `json:"name,omitempty" yaml:"name,omitempty"`
	Labels      map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty" yaml:"annotations,omitempty"`
}

func (b *Base) Check(kind string) bool {
	if b == nil {
		log.Error().Msg("Base is nil")
		return false
	}
	if b.APIVersion != DeploytoAPIVersion {
		log.Debug().Str("apiVersion", b.APIVersion).Str("want", DeploytoAPIVersion).Msg("The apiVersion does not match")
		return false
	}
	if b.Kind != kind {
		log.Debug().Str("kind", b.Kind).Str("want", kind).Msg("The kind does not match")
		return false
	}
	if b.Meta == nil || b.Meta.Name == "" {
		log.Debug().Msg("metadata/name not set")
		return false
	}
	return true
}
