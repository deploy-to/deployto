package types

import (
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

type TemplateArg = struct {
	Resource string
	Version  string
	Name     string
	Values   Values `mapstructure:",remain"`
}

func DecodeTemplateArg(values any) (templateArg *TemplateArg) {
	templateArg = &TemplateArg{}
	if values == nil {
		log.Info().Msg("DecodeTemplate - input values is nil")
		return templateArg
	}
	err := mapstructure.Decode(values, templateArg)
	if err != nil {
		log.Error().Err(err).Msg("DecodeScript error")
		return nil
	}
	return templateArg
}
