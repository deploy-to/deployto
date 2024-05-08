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
func Component(target *types.Target, repositoryFS *filesystem.Filesystem, workdir string, aliases []string, rootValues, context types.Values) (output types.Values, err error) {
	if context == nil {
		log.Error().Msg("want context")
		return nil, errors.New("WANT CONTEXT")
	}
	output = make(types.Values)

	// COMPONENTS
	comps := yaml.Get[types.Component](repositoryFS, repositoryFS.FS.Join(workdir, filesystem.DeploytoDirName))
	for _, comp := range comps {
		compOutput, err := RunSingleComponent(target, aliases, rootValues, context, comp)
		if err != nil {
			return nil, err
		}
		output[buildAlias(aliases)] = compOutput
	}

	return output, err
}

func RunSingleComponent(target *types.Target, aliases []string, rootContext, context types.Values, comp *types.Component) (output types.Values, err error) {
	if len(aliases) == 0 { //is first component (application)
		aliases = []string{comp.Meta.Name}
	}

	l := log.With().Strs("aliases", aliases).Logger()
	dependenciesOutput := make(types.Values)

	// зависимость git  выполняется всегда
	l.Debug().Msg("Get commit hash and tags")
	context["git"] = gitclient.GetValues(comp.Status.Filesystem, filepath.Dir(comp.Status.FileName))
	context["component"] = aliases[len(aliases)-1]
	context["alias"] = buildAlias(aliases)

	dependencies := types.Get(comp.Spec, types.Values(nil), "dependencies")
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

		dependencyOutput, err := RunScript(target, comp.Status.Filesystem, filepath.Dir(comp.Status.FileName),
			dependencyAliases,
			d,
			rootContext, context)
		if err != nil {
			l.Error().Err(err).Msg("RunScript error")
		}
		dependenciesOutput[alias] = dependencyOutput
	}

	if types.Exists(comp.Spec, "script") {
		compScript := types.DecodeScript(types.Get(comp.Spec, types.Values(nil), "script"))

		scriptContext := types.MergeValues(dependenciesOutput, context)
		output, err = RunScript(target, comp.Status.Filesystem, filepath.Dir(comp.Status.FileName),
			aliases,
			compScript,
			rootContext, scriptContext)
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
