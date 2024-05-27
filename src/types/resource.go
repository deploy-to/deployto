package types

import (
	"github.com/go-viper/mapstructure/v2"
	"github.com/rs/zerolog/log"
)

type ResourceArg = struct {
	Version string
	Name    string
	Values  Values `mapstructure:",remain"`
}

func DecodeResourceArg(values any) (resourceArg *ResourceArg) {
	resourceArg = &ResourceArg{}
	if values == nil {
		log.Info().Msg("DecodeTemplate - input values is nil")
		return resourceArg
	}
	err := mapstructure.Decode(values, resourceArg)
	if err != nil {
		log.Error().Err(err).Msg("DecodeScript error")
		return nil
	}
	return resourceArg
}
