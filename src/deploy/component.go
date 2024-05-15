package deploy

import (
	"deployto/src"
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
func Component(target *types.Target, repositoryFS *filesystem.Filesystem, workdir string, aliases []string, rootValues, context types.Values, ContextDump *src.ContextDump) (output types.Values, err error) {
	if context == nil {
		log.Error().Msg("want context")
		return nil, errors.New("WANT CONTEXT")
	}
	output = make(types.Values)

	// COMPONENTS
	comps := yaml.Get[types.Component](repositoryFS, repositoryFS.FS.Join(workdir, filesystem.DeploytoDirName))
	for _, comp := range comps {
		if len(aliases) == 0 {
			aliases = []string{comp.Meta.Name}
		}
		compOutput, err := RunSingleComponent(target, aliases, rootValues, context, comp, ContextDump.Next(buildAlias(aliases)))
		if err != nil {
			return nil, err
		}
		output[buildAlias(aliases)] = compOutput
	}

	return output, err
}

func RunSingleComponent(target *types.Target, aliases []string, rootContext, context types.Values, comp *types.Component, ContextDump *src.ContextDump) (output types.Values, err error) {
	l := log.With().Strs("aliases", aliases).Logger()
	output = make(types.Values)

	// зависимость git  выполняется всегда
	l.Debug().Msg("Get commit hash and tags")
	context["git"] = gitclient.GetValues(comp.Status.Filesystem, filepath.Dir(comp.Status.FileName))
	context["component"] = aliases[len(aliases)-1]
	context["alias"] = buildAlias(aliases)

	ordered := types.GetTheOrderOfResource(comp.Spec)
	for _, sameOrder := range ordered {
		context = types.MergeValues(context, output)
		// TODO sameOrder parallel deploy (+config)
		for _, script := range sameOrder {
			resourceAliases := aliases
			if script.Shared {
				log.Debug().Strs("aliases", aliases).Str("script", script.Alias).Msg("Script is shared")
				resourceAliases = []string{script.Alias}
			} else {
				if script.Alias != "" {
					resourceAliases = append(resourceAliases, script.Alias)
				}
			}

			scriptOutput, err := RunScript(target, comp.Status.Filesystem, filepath.Dir(comp.Status.FileName),
				resourceAliases,
				script,
				rootContext, context, ContextDump.Next(buildAlias(resourceAliases)))
			if err != nil {
				l.Error().Err(err).Strs("aliases", resourceAliases).Msg("Dependency error")
				return output, err
			}
			if script.Alias == "" {
				output = types.MergeValues(output, scriptOutput)
			} else {
				output[script.Alias] = scriptOutput
			}
		}
	}
	return output, nil
}

func buildAlias(names []string) string {
	return strings.Join(names, "-")
}
