package src

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigName("config.yaml")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.deployto")
	viper.AddConfigPath("./.deployto")
	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			panic(fmt.Errorf("fatal error config file: %w", err))
		}
	}
}

func GetTemplateRepositories() (result []string) {
	templateRepositories := viper.GetStringSlice("templateRepositories")
	for i := 0; i < len(templateRepositories); i++ {
		r := templateRepositories[i]
		if strings.HasPrefix(r, "file://") || strings.HasPrefix(r, "http://") || strings.HasPrefix(r, "https://") {
			result = append(result, r)
		} else {
			log.Error().Str("repository", r).Msg("Unsupport repository")
		}
	}
	return result
}
