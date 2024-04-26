package cmd

import (
	"deployto/src/deploy"
	"deployto/src/gitclient"
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
		log.Error().Err(err).Str("path", path).Msg("Components search error")
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
		log.Error().Int("len(environments)", len(environments)).Str("path", comps[0].StatusGetPath()).Str("waitEnvironment", environmentArg).Msg("environment not found")
		return errors.New("ENVIRONMENT NOT FOUND")
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
		e := Deploy("TODO, for each target", []string{c.Base.Meta.Name}, c.StatusGetPath(), c.Spec, make(map[string]any))
		if e != nil {
			log.Error().Err(e).Msg("Component deploy error")
			err = errors.Join(err, e)
		}
	}
	return err
}

// TODO заглушка, для вызова job
// output где взять образ
func DoJob(s string) (map[string]any, error) {
	return nil, nil
}

func Deploy(kubeconfig string, aliases []string, workDir string, as *types.ComponentSpec, values map[string]any) (resultError error) {
	l := log.With().Strs("aliases", aliases).Logger()
	if as == nil {
		l.Debug().Msg("ComponentSpec is nil")
		return nil
	}

	l.Debug().Msg("Get commit hash and tags")
	values["git"] = gitclient.GetValues(workDir)
	l.Debug().Msg("TODO  BUILD && Push")
	values["image"], resultError = DoJob(as.BuildJob)
	//TODO положить репозиторий/образ/тег в values

	//TODO положить в values всё из as.Script.Input
	l.Debug().Msg("TODO Preparing the values")

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

		// for dependencies, if the script is not defined, I will try to get the default script by kind
		l.Debug().Strs("alias", dependencyAliases).Msg("Run dependency")
		o, e := RunScript(kubeconfig, dependencyAliases, d.Kind, d.Script, values)
		if e != nil {
			l.Error().Err(e).Msg("RunScript error")
			resultError = errors.Join(resultError, e)
		}
		values[buildAlias(dependencyAliases)] = o
	}
	// run the script on the component only if it is defined
	if as.Script != nil {
		o, e := RunScript(kubeconfig, aliases, "component", as.Script, values)
		if e != nil {
			l.Error().Err(e).Msg("RunScript error")
			resultError = errors.Join(resultError, e)
		}
		values[buildAlias(aliases)] = o
	}
	return resultError
}

func buildAlias(names []string) string {
	return strings.Join(names, "-")
}

func RunScript(kubeconfig string, names []string, kind string, script *types.Script, input map[string]any) (output map[string]any, err error) {
	l := log.With().Strs("names", names).Logger()
	if script != nil {
		if runScript, ok := deploy.RunScripts[script.Type]; ok {
			l.Debug().Str("scriptType", script.Type).Msg("Run script")
			return runScript(kubeconfig, names, kind, script, input)
		}
	}
	if runScript, ok := deploy.RunScripts["component"]; ok {
		l.Debug().Msg("script.Type not defined, run script for component")
		return runScript(kubeconfig, names, kind, script, input)
	}
	l.Error().Msg("RunScripts function not found")
	return nil, errors.New("RUNSCRIPT FUNCTION NOT FOUND")
}
