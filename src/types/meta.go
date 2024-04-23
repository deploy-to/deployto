package types

import (
	"github.com/rs/zerolog/log"
)

const DeploytoAPIVersion = "deployto.dev/v1beta1"

type Meta struct {
	Kind       string    `json:"kind,omitempty" yaml:"kind,omitempty"`
	APIVersion string    `json:"apiVersion,omitempty" yaml:"apiVersion,omitempty"`
	Meta       *MetaData `json:"meta,omitempty" yaml:"meta,omitempty"`
}

type MetaData struct {
	Name        string            `json:"name,omitempty" yaml:"name,omitempty"`
	Labels      map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty" yaml:"annotations,omitempty"`
}

func (m *Meta) Is(kind string) bool {
	if m == nil {
		log.Error().Msg("Meta is nil")
		return false
	}
	if m.APIVersion != DeploytoAPIVersion {
		log.Debug().Str("apiVersion", m.APIVersion).Str("want", DeploytoAPIVersion).Msg("The apiVersion does not match")
		return false
	}
	if m.Kind != kind {
		log.Debug().Str("kind", m.Kind).Str("want", kind).Msg("The kind does not match")
		return false
	}
	return true
}
