package cmd

import (
	"deployto/src/deploy"
	"deployto/src/types"
	"deployto/src/yaml"
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"

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
	// COMPONENTS
	comps, err := yaml.GetComponent(path)
	if err != nil {
		log.Error().Err(err).Str("path", path).Msg("Application/Components search error")
		return err
	}
	// Envirement
	environments := yaml.Get[types.Environment](comps[0].StatusGetPath())
	var environment *types.Environment
	for _, e := range environments {
		if e.Base.Meta.Name == environmentArg {
			environment = e
		}
	}
	if environment == nil {
		log.Error().Int("len(environments)", len(environments)).Str("path", comps[0].StatusGetPath()).Str("waitEnvironment", environmentArg).Msg("environment ")
		return errors.New("APP NOT FOUND")
	}
	log.Debug().Str("name", environment.Base.Meta.Name).Msg("Environment found")
	// Targets
	var targets []*types.Target
	for _, t := range yaml.Get[types.Target](comps[0].StatusGetPath()) {
		if slices.Contains(environment.Spec.Targets, t.Base.Meta.Name) {
			targets = append(targets, t)
		}
	}
	if len(targets) != len(environment.Spec.Targets) {
		log.Error().Int("len(targets)", len(targets)).Int("len(environment.Spec.Targets)", len(environment.Spec.Targets)).Msg("Target not found")
		return errors.New("TARGET NOT FOUND")
	}
	log.Debug().Int("len(targets)", len(targets)).Msg("Targets found")

	fmt.Printf("-- Call for components --------------------------------------\n")
	for _, c := range comps {
		fmt.Printf("  Path: %v\n", c.StatusGetPath())
		fmt.Printf("  File: %v\n", c.Base.Status.FileName)
		fmt.Printf("  Name: %v\n", c.Base.Meta.Name)
	}
	fmt.Printf("-- Environment ----------------------------------------------\n")
	fmt.Printf("  File:   %v\n", environment.Base.Status.FileName)
	fmt.Printf("  Name:  %v\n", environment.Base.Meta.Name)
	fmt.Printf("-- Targets --------------------------------------------------\n")
	for _, t := range targets {
		fmt.Printf("  File: %v\n", t.Base.Status.FileName)
		fmt.Printf("  Name: %v\n", t.Base.Meta.Name)
	}

	for _, c := range comps {
		//деплою все найденные компоненты
		values := make(map[string]any)
		values["git"] = GetGitValues(c.StatusGetPath())
		//values["image"] := CI()

		e := Deploy("TODO, for each target", []string{c.Base.Meta.Name}, c.Spec, values)
		if e != nil {
			log.Error().Err(e).Msg("Component deploy error")
			err = errors.Join(err, e)
		}
	}
	return err
}

func GetGitValues(path string) map[string]any {
	//TODO
	return nil
}

func Deploy(kubeconfig string, aliases []string, as *types.ComponentSpec, values map[string]interface{}) (resultError error) {
	l := log.With().Strs("aliases", aliases).Logger()
	//TODO values из
	l.Debug().Msg("TODO  BUILD && Push") // output где взять образ
	//TODO положить репозиторий/образ/тег в values

	//TODO положить в values всё из as.Script.Input
	l.Debug().Msg("TODO Preparing the values")

	//TODO values из

	//Run dependency
	for _, d := range as.Dependencies {
		var dependencyAliases []string
		if !d.Root {
			dependencyAliases = aliases
		}
		if d.Name == "" {
			dependencyAliases = append(dependencyAliases, d.Kind)
		} else {
			dependencyAliases = append(dependencyAliases, d.Name)
		}

		if _, ok := values[buildAlias(aliases)]; ok {
			l.Debug().Strs("alias", dependencyAliases).Msg("Deployed earlier")
			continue
		}

		l.Debug().Strs("alias", dependencyAliases).Msg("TODO Run dependency")
		o, e := RunScript(kubeconfig, dependencyAliases, d.Kind, d.Script, values)
		if e != nil {
			l.Error().Err(e).Msg("RunScript error")
			resultError = errors.Join(resultError, e)
		}
		values[buildAlias(dependencyAliases)] = o
	}
	// Run script
	o, e := RunScript(kubeconfig, aliases, "component", as.Script, values)
	if e != nil {
		l.Error().Err(e).Msg("RunScript error")
		resultError = errors.Join(resultError, e)
	}
	values[buildAlias(aliases)] = o

	return resultError
}

func buildAlias(names []string) string {
	return strings.Join(names, "-")
}

func RunScript(kubeconfig string, names []string, kind string, script *types.Script, input map[string]any) (output map[string]any, err error) {
	l := log.With().Strs("names", names).Logger()
	if script == nil {
		l.Debug().Msg("Script not defined")
		return nil, nil
	}
	l.Debug().Str("scriptType", script.Type).Msg("Run script")
	if runScript, ok := deploy.RunScripts[script.Type]; ok {
		return runScript(kubeconfig, names, kind, script, input)
	}
	if runScript, ok := deploy.RunScripts["component"]; ok {
		return runScript(kubeconfig, names, kind, script, input)
	}

	l.Error().Msg("RunScripts function not found")
	return nil, errors.New("RUNSCRIPT FUNCTION NOT FOUND")
}
