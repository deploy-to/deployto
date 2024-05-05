package deploy

import (
	"deployto/src/filesystem"
	"deployto/src/gitclient"
	"deployto/src/types"
	"deployto/src/yaml"
	"errors"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
)

func init() {
	RunScriptFuncImplementations["component"] = Component
}
func Component(target *types.Target, repositoryFS *filesystem.Filesystem, workdir string, aliases []string, rootValues, input types.Values) (output types.Values, err error) {
	output = make(types.Values)

	// COMPONENTS
	comps := yaml.Get[types.Component](repositoryFS, repositoryFS.FS.Join(workdir, filesystem.DeploytoDirName))
	for _, c := range comps {
		cOutput, err := RunSingleComponent(target, aliases, rootValues, input, c)
		if err != nil {
			return nil, err
		}
		output[buildAlias(aliases)] = cOutput
	}

	return output, err
}

func RunSingleComponent(target *types.Target, aliases []string, rootValues, input types.Values, c *types.Component) (output types.Values, err error) {
	if len(aliases) == 0 { //is first component (application)
		aliases = []string{c.Meta.Name}
	}

	l := log.With().Strs("aliases", aliases).Logger()
	dependenciesOutput := make(types.Values)

	// зависимость git  выполняется всегда
	l.Debug().Msg("Get commit hash and tags")
	dependenciesOutput["git"] = gitclient.GetValues(c.Status.Filesystem, c.Status.FileName)

	dependencies := types.Get(c.Spec, types.Values(nil), "dependencies")
	for alias, dependencyAsMap := range dependencies {
		d := types.DecodeScript(dependencyAsMap)
		if d == nil {
			l.Error().Msg("DecodeScript return nil")
			return nil, errors.New("DecodeScript return nil")
		}
		var dependencyAliases []string
		if d.Root {
			log.Debug().Strs("aliases", aliases).Str("dependency", alias).Msg("Dependency is root")
			dependencyAliases = []string{alias}
		} else {
			dependencyAliases = append(aliases, alias)
		}

		dependencyOutput, err := RunScript(target, c.Status.Filesystem, filepath.Dir(c.Status.FileName),
			dependencyAliases,
			rootValues,
			d, input)
		if err != nil {
			l.Error().Err(err).Msg("RunScript error")
		}
		dependenciesOutput[alias] = dependencyOutput
		if d.Root {
			rootValues[alias] = dependencyOutput
		}
	}

	if types.Exists(c.Spec, "script") {
		compScript := types.DecodeScript(types.Get(c.Spec, types.Values(nil), "script"))

		scriptContext := types.MergeValues(dependenciesOutput, input)
		output, err = RunScript(target, c.Status.Filesystem, filepath.Dir(c.Status.FileName),
			aliases,
			rootValues,
			compScript, scriptContext)
		if err != nil {
			l.Error().Err(err).Msg("RunScript error")
		}
	}
	if output == nil {
		output = make(types.Values)
	}
	output["dependencies"] = dependenciesOutput
	return output, nil
}

func buildAlias(names []string) string {
	return strings.Join(names, "-")
}
