package types

import (
	"github.com/go-viper/mapstructure/v2"
	"github.com/rs/zerolog/log"
)

type Dependency = Script
type Script = struct {
	Order         int
	Type          string
	Repository    string
	Path          string
	Shared        bool
	Alias         string // служит для формерования контекста / хоста
	Name          string // указывает на имя ресурса или job, или на имя сhart для helm
	OutputMapping Values
	Values        Values `mapstructure:",remain"` //Values хранятся на том же уровне, что и Order, Type...
}

func DecodeScript(defaultAlias string, values Values) (script *Script) {
	if values == nil {
		log.Debug().Str("alias", defaultAlias).Msg("DecodeScript - input values is nil")
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
	if _, nameExists := values["name"]; !nameExists {
		script.Name = defaultAlias
	}
	if script.Values == nil {
		script.Values = make(Values)
	}
	return script
}
