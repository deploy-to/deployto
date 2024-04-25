package cmd

import (
	"deployto/src/types"
	"deployto/src/yaml"
	"errors"
	"fmt"
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
	app, comps, err := yaml.GetAppComps(path)
	if err != nil {
		log.Error().Err(err).Str("path", path).Msg("Application/Components search error")
		return err
	}
	// Envirement
	environments := yaml.Get[types.Environment](app.StatusGetPath())
	var environment *types.Environment
	for _, e := range environments {
		if e.Base.Meta.Name == environmentArg {
			environment = e
		}
	}
	if environment == nil {
		log.Error().Int("len(environments)", len(environments)).Str("path", app.StatusGetPath()).Str("waitEnvironment", environmentArg).Msg("environment ")
		return errors.New("APP NOT FOUND")
	}
	log.Debug().Str("name", environment.Base.Meta.Name).Msg("Environment found")
	// Targets
	var targets []*types.Target
	for _, t := range yaml.Get[types.Target](app.StatusGetPath()) {
		if slices.Contains(environment.Spec.Targets, t.Base.Meta.Name) {
			targets = append(targets, t)
		}
	}
	if len(targets) != len(environment.Spec.Targets) {
		log.Error().Int("len(targets)", len(targets)).Int("len(environment.Spec.Targets)", len(environment.Spec.Targets)).Msg("Target not found")
		return errors.New("TARGET NOT FOUND")
	}
	log.Debug().Int("len(targets)", len(targets)).Msg("Targets found")

	fmt.Printf("-- Application ----------------------------------------------\n")
	fmt.Printf("  Path: %v\n", app.StatusGetPath())
	fmt.Printf("  File: %v\n", app.Base.Status.FileName)
	fmt.Printf("  Name: %v\n", app.Base.Meta.Name)
	fmt.Printf("-- Environment ----------------------------------------------\n")
	fmt.Printf("  File:   %v\n", environment.Base.Status.FileName)
	fmt.Printf("  Name:  %v\n", environment.Base.Meta.Name)
	fmt.Printf("-- Targets --------------------------------------------------\n")
	for _, t := range targets {
		fmt.Printf("  File: %v\n", t.Base.Status.FileName)
		fmt.Printf("  Name: %v\n", t.Base.Meta.Name)
	}
	fmt.Printf("-- Call for components --------------------------------------\n")
	for _, c := range comps {
		fmt.Printf("  File: %v\n", c.Base.Status.FileName)
		fmt.Printf("  Name: %v\n", c.Base.Meta.Name)
	}

	//собираю все зависимости и их пути, начиная с
	if len(comps) > 0 {
		for _, c := range comps {
			Deploy(&c.Base, c.Spec)
		}
	} else {
		Deploy(&app.Base, app.Spec)
	}

	return nil
}

func Deploy(base *types.Base, as *types.ApplicationSpec) {
	l := log.With().Str("name", base.Meta.Name).Logger()

	l.Debug().Msg("Preparing the values")
	//Run dependency
	for _, d := range as.Dependencies {
		l.Debug().Str("DependencyName", d.Name).Str("DependencyAlias", d.Alias).Msg("Run dependency")
	}
	// Run script
	if as.Script == nil {
		l.Debug().Msg("Script is not defined")
	} else {
		l.Debug().Msg("Run script")
	}
}
