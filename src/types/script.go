package types

import (
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

type Dependency = Script
type Script = struct {
	Order         int
	Type          string
	Repository    string
	Path          string
	Shared        bool
	Alias         string
	OutputMapping Values
	Values        Values `mapstructure:",remain"` //Values хранятся на том же уровне, что и Order, Type...
}

func DecodeScript(defaultAlias string, values Values) (script *Script) {
	if values == nil {
		log.Info().Msg("DecodeScript - input values is nil")
		return nil
	}
	script = &Script{}
	err := mapstructure.Decode(values, script)
	if err != nil {
		log.Error().Err(err).Msg("DecodeScript error")
		return nil
	}
	if _, orderExists := values["order"]; !orderExists {
		script.Order = 100
	}
	if _, aliasExists := values["alias"]; !aliasExists {
		script.Alias = defaultAlias
	}

	return script
}
