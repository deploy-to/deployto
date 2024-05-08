package cmd

import (
	"deployto/src/deploy"
	"deployto/src/types"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
)

var Flags = []cli.Flag{
	altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "log-format",
		Aliases: []string{"lg"},
		Value:   "pretty",
		Usage:   "Log Format: json, pretty",
	}),
	altsrc.NewStringFlag(&cli.StringFlag{
		Name:    "log-level",
		Aliases: []string{"ll"},
		Value:   "info",
		Usage:   "Log level: trace, debug, warn, info, fatal, panic, absent, disable",
	}),
	altsrc.NewStringFlag(&cli.StringFlag{
		Name:        "kubeconfig",
		EnvVars:     []string{"KUBECONFIG"},
		Usage:       "Set to use when target.kubeconfig.usedefault is set.",
		Destination: &types.SystemKubeconfig,
	}),
	altsrc.NewStringSliceFlag(&cli.StringSliceFlag{
		Name:        "template-repositories",
		Aliases:     []string{"tr"},
		Usage:       "URIs to template repository file://XXX or https://XXX.git ",
		Destination: &deploy.TemplateRepositories,
	}),
	&cli.StringFlag{
		Name:    "config",
		Aliases: []string{"c"},
		Usage:   "Set config file uri/path  https://cli.urfave.org/v2/examples/flags/#values-from-alternate-input-sources-yaml-toml-and-others",
	},
}

func LoadYamlConfig(ctx *cli.Context) error {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		log.Error().Err(err).Msg("homeless")
	}

	for _, filePath := range []string{
		ctx.String("config"),
		".deployto/config.yaml",
		filepath.Join(userHomeDir, ".deployto/config.yaml"),
	} {
		if filePath != "" {
			inputSource, err := altsrc.NewYamlSourceFromFile(filePath)
			if err != nil {
				log.Error().Err(err).Str("config", filePath).Msg("Read config error")
				return err
			}
			return altsrc.ApplyInputSourceValues(ctx, inputSource, Flags)
		}
	}
	log.Info().Msg("configuration file not found")
	return nil
}
