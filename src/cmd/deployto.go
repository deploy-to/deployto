package cmd

import (
	"deployto/src/deploy"
	"deployto/src/helper"
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

	path, err = helper.GetProjectRoot(path, helper.DeploytoPath)
	if err != nil {
		log.Error().Err(err).Msg("Get DeploytoPath error")
		return err
	}

	// Envirement
	environments := yaml.Get[types.Environment](path)
	var environment *types.Environment
	for _, e := range environments {
		if e.Base.Meta.Name == environmentArg {
			environment = e
		}
	}
	if environment == nil {
		log.Error().Int("len(environments)", len(environments)).Str("path", path).Str("waitEnvironment", environmentArg).Msg("environment not found")
		return errors.New("ENVIRONMENT NOT FOUND")
	}
	log.Debug().Str("name", environment.Base.Meta.Name).Msg("Environment found")
	// Targets
	var targets []*types.Target
	for _, t := range yaml.Get[types.Target](path) {
		if slices.Contains(environment.Spec.Targets, t.Base.Meta.Name) {
			targets = append(targets, t)
		}
	}
	if len(targets) != len(environment.Spec.Targets) {
		log.Error().Int("len(targets)", len(targets)).Int("len(environment.Spec.Targets)", len(environment.Spec.Targets)).Msg("Target not found")
		return errors.New("TARGET NOT FOUND")
	}
	log.Debug().Int("len(targets)", len(targets)).Msg("Targets found")

	log.Info().Str("file", environment.Base.Status.FileName).Str("name", environment.Base.Meta.Name).Msg("Deploy environment")
	for _, t := range targets {
		log.Info().Str("file", t.Base.Status.FileName).Str("name", t.Base.Meta.Name).Msg("Deploy target")

		rootValues := make(types.Values)
		//TODO позволить пользователю передавать в deploy.Component значения values заданные в командной строке / файле и т.п.
		_, e := deploy.Component(t,
			path,
			nil,
			rootValues, types.Values(nil))
		if e != nil {
			log.Error().Err(e).Msg("Component deploy error")
			err = errors.Join(err, e)
		}

		//TODO Run target script (move this logic to src/deploy/)
	}
	//TODO Run Env script (move this logic to src/deploy/?   Or stop move on target?)
	return err
}
