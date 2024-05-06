package types

import (
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

type Dependency = Script
type Script = struct {
	Type          string
	Repository    string
	Path          string
	Root          bool
	OutputMapping Values
	Values        Values `mapstructure:",remain"` //Values хранятся на том же уровне, что и Type, Root
}

func DecodeScript(values any) (script *Script) {
	if values == nil {
		log.Info().Msg("DecodeScript - input values is nil")
		return &Script{}
	}
	script = &Script{}
	err := mapstructure.Decode(values, script)
	if err != nil {
		log.Error().Err(err).Msg("DecodeScript error")
		return nil
	}
	return script
}
