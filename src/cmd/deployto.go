package cmd

import (
	"deployto/src/types"
	"deployto/src/yaml"
	"errors"
	"os"
	"slices"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

func Deployto(cCtx *cli.Context) error {
	environmentArg := cCtx.Args().First()
	if len(environmentArg) == 0 {
		environmentArg = "local"
	}
	log.Debug().Str("environment", environmentArg).Msg("start deployto")

	path, err := os.Getwd()
	if err != nil {
		log.Error().Err(err).Msg("Get workdir error")
		return err
	}
	// Application
	appPath := yaml.GetAppPath(path)
	apps := yaml.Get[types.Application](appPath)
	if len(apps) != 1 {
		log.Error().Int("len(app)", len(apps)).Str("path", appPath).Msg("wait one app")
		return errors.New("APP NOT FOUND")
	}
	app := apps[0]
	log.Debug().Str("name", app.Meta.Meta.Name).Msg("Application found")
	// Envirement
	environments := yaml.Get[types.Environment](appPath)
	var environment *types.Environment
	for _, e := range environments {
		if e.Meta.Meta.Name == environmentArg {
			environment = e
		}
	}
	if environment == nil {
		log.Error().Int("len(environments)", len(environments)).Str("path", appPath).Str("waitEnvironment", environmentArg).Msg("environment ")
		return errors.New("APP NOT FOUND")
	}
	log.Debug().Str("name", environment.Meta.Meta.Name).Msg("Environment found")
	// Targets
	var targets []*types.Target
	for _, t := range yaml.Get[types.Target](appPath) {
		if slices.Contains(environment.Spec.Targets, t.Meta.Meta.Name) {
			targets = append(targets, t)
		}
	}
	if len(targets) == len(environment.Spec.Targets) {
		log.Error().Int("len(targets)", len(targets)).Int("len(environment.Spec.Targets)", len(environment.Spec.Targets)).Msg("Target not found")
		return errors.New("TARGET NOT FOUND")
	}
	log.Debug().Int("len(targets)", len(targets)).Msg("Targets found")

	return nil
}
