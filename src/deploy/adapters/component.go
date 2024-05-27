package adapters

import (
	"deployto/src/deploy"
	"deployto/src/filesystem"
	"deployto/src/gitclient"
	"deployto/src/types"
	"deployto/src/yaml"
	"errors"
	"path/filepath"

	"github.com/rs/zerolog/log"
)

func init() {
	deploy.DefaultAdapters["component"] = (*component)(nil)
}

type component struct{}

func (c *component) Apply(d *deploy.DeployState, script *types.Script, compContext types.Values) (output types.Values, err error) {
	if compContext == nil {
		log.Error().Msg("want context")
		return nil, errors.New("WANT CONTEXT")
	}
	output = make(types.Values)

	// COMPONENTS
	comps := yaml.Get[types.Component](d.FS, d.FS.FS.Join(d.Workdir, filesystem.DeploytoDirName))
	aliases := d.Aliases
	for _, comp := range comps {
		if d.Root == d {
			aliases = []string{comp.Meta.Name}
		}
		compOutput, err := ApplySingleComponent(d.Child(d.FS, d.Workdir, aliases), comp, compContext)
		if err != nil {
			return nil, err
		}
		output[deploy.BuildAlias(aliases)] = compOutput
	}

	return output, err
}

func (c *component) Destroy(d *deploy.DeployState, script *types.Script, compContext types.Values) error {
	panic("NOT IMPLIMENTED")
}

func ApplySingleComponent(d *deploy.DeployState, comp *types.Component, context types.Values) (output types.Values, err error) {
	l := log.With().Strs("aliases", d.Aliases).Logger()
	output = make(types.Values)

	// зависимость git  выполняется всегда
	l.Debug().Msg("Get commit hash and tags")
	context["git"] = gitclient.GetValues(comp.Status.Filesystem, filepath.Dir(comp.Status.FileName))
	context["component"] = d.Aliases[len(d.Aliases)-1]
	context["alias"] = deploy.BuildAlias(d.Aliases)

	ordered := types.GetTheOrderOfResource(comp.Spec)
	for _, sameOrder := range ordered {
		context = types.MergeValues(context, output)
		// TODO sameOrder parallel deploy (+config)
		for _, script := range sameOrder {
			resourceAliases := d.Aliases
			if script.Shared {
				log.Debug().Strs("aliases", d.Aliases).Str("script", script.Alias).Msg("Script is shared")
				resourceAliases = []string{script.Alias}
			} else {
				if script.Alias != "" {
					resourceAliases = append(resourceAliases, script.Alias)
				}
			}

			scriptOutput, err := ApplyScript(d, comp, resourceAliases, script, context)
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
