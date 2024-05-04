package deploy

import (
	"deployto/src/filesystem"
	"deployto/src/gitclient"
	"deployto/src/types"
	"deployto/src/yaml"
	"strings"

	"github.com/rs/zerolog/log"
)

func init() {
	RunScriptFuncImplementations["component"] = Component
}
func Component(target *types.Target, filesystem *filesystem.Filesystem, aliases []string, rootValues, input types.Values) (output types.Values, err error) {
	output = make(types.Values)

	// workdir = filepath.Join(workdir, repository[len("file://"):], path)

	// COMPONENTS
	comps := yaml.Get[types.Component](filesystem, "/")
	for _, c := range comps {
		cOutput, err := RunSingleComponent(target, filesystem, aliases, rootValues, input, c)
		if err != nil {
			return nil, err
		}
		output[buildAlias(aliases)] = cOutput
	}

	return output, err
}

func buildAlias(names []string) string {
	return strings.Join(names, "-")
}

func RunSingleComponent(target *types.Target, fs *filesystem.Filesystem, aliases []string, rootValues, input types.Values, c *types.Component) (types.Values, error) {
	compScript := types.DecodeScript(types.Get(c.Spec, types.Values(nil), "script"))
	if compScript.Root {
		log.Debug().Strs("aliases", aliases).Str("component", c.Meta.Name).Msg("Component script is root")
		aliases = []string{c.Meta.Name}
	} else {
		aliases = append(aliases, c.Meta.Name)
	}

	l := log.With().Strs("aliases", aliases).Logger()
	dependenciesOutput := make(types.Values)

	// зависимость git  выполняется всегда
	l.Debug().Msg("Get commit hash and tags")
	dependenciesOutput["git"] = gitclient.GetValues(fs, "/")

	dependencies := types.Get(c.Spec, map[string]any(nil), "dependencies")
	for alias, dependencyAsMap := range dependencies {
		d := types.DecodeScript(dependencyAsMap)

		var dependencyAliases []string
		if compScript.Root {
			log.Debug().Strs("aliases", aliases).Str("dependency", alias).Msg("Dependency is root")
			dependencyAliases = []string{alias}
		} else {
			dependencyAliases = append(aliases, alias)
		}

		dependencyOutput, e := RunScript(target, fs,
			dependencyAliases,
			rootValues,
			d, input)
		if e != nil {
			l.Error().Err(e).Msg("RunScript error")
		}
		dependenciesOutput[alias] = dependencyOutput
		if d != nil && d.Root {
			rootValues[alias] = dependencyOutput
		}
	}

	scriptContext := types.MergeValues(dependenciesOutput, input)
	scriptOutput, e := RunScript(target, fs,
		aliases,
		rootValues,
		compScript, scriptContext)
	if e != nil {
		l.Error().Err(e).Msg("RunScript error")
	}
	if scriptOutput == nil {
		scriptOutput = make(types.Values)
	}
	scriptOutput["dependencies"] = dependenciesOutput
	return scriptOutput, nil
}
