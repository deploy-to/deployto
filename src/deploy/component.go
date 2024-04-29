package deploy

import (
	"deployto/src/gitclient"
	"deployto/src/types"
	"deployto/src/yaml"
	"path/filepath"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
)

func init() {
	RunScriptFuncImplementations["component"] = Component
}
func Component(kubeconfig string, workdir string, aliases []string, rootValues, input types.Values) (output types.Values, err error) {
	output = make(types.Values)
	repository := types.Get(input, "", "repository")
	path := types.Get(input, "", "path")
	if repository == "" || strings.HasPrefix(repository, "file://") {
		if repository == "" {
			workdir = filepath.Join(workdir, path)
		} else {
			workdir = filepath.Join(workdir, repository[len("file://"):], path)
		}
	} else {
		panic("not implimented")
	}

	// COMPONENTS
	comps, err := yaml.GetComponent(workdir)
	if err != nil {
		log.Error().Err(err).Str("path", workdir).Msg("Components search error")
		return nil, err
	}
	for _, c := range comps {
		workdir = c.GetDir()

		var aliases []string
		alias := types.Get(input, c.Meta.Name, "alias")
		if types.Get(input, false, "root") {
			aliases = []string{alias}
		} else {
			aliases = append(aliases, alias)
		}

		l := log.With().Strs("aliases", aliases).Logger()
		dependenciesOutput := make(types.Values)

		// зависимость git  выполняется всегда
		l.Debug().Msg("Get commit hash and tags")
		dependenciesOutput["git"] = gitclient.GetValues(workdir)

		dependencies := types.Get(c.Spec, map[string]any(nil), "dependencies")
		for alias, dependencyAsMap := range dependencies {
			var d types.Dependency
			err := mapstructure.Decode(dependencyAsMap, &d)
			if err != nil {
				l.Error().Str("alias", alias).Err(err).Msg("dependency is not types.Dependency")
				return nil, err
			}
			dependencyAliases := append(aliases, alias)
			if d.Root {
				dependencyAliases = []string{alias}
			}

			if theDependencyWasDeployedEarlier, ok := rootValues[buildAlias(aliases)]; ok {
				l.Info().Strs("alias", dependencyAliases).Msg("Deployed earlier")
				dependenciesOutput[alias] = theDependencyWasDeployedEarlier
				continue
			}

			dependencyOutput, e := RunScript(kubeconfig, workdir,
				dependencyAliases,
				rootValues,
				d.Script, input)
			if e != nil {
				l.Error().Err(e).Msg("RunScript error")
			}
			dependenciesOutput[alias] = dependencyOutput
			if d.Root {
				rootValues[alias] = dependencyOutput
			}
		}

		scriptContext := types.MergeValues(dependenciesOutput, input)
		scriptOutput, e := RunScript(kubeconfig, workdir,
			aliases,
			rootValues,
			types.Get(c.Spec, types.Values(nil), "script"), scriptContext)
		if e != nil {
			l.Error().Err(e).Msg("RunScript error")
		}
		if scriptOutput == nil {
			scriptOutput = make(types.Values)
		}
		scriptOutput["dependencies"] = dependenciesOutput
		output[buildAlias(aliases)] = scriptOutput
	}

	return output, err
}

func buildAlias(names []string) string {
	return strings.Join(names, "-")
}
