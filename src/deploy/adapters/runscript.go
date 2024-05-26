package adapters

import (
	"deployto/src/deploy"
	"deployto/src/filesystem"
	"deployto/src/types"
	"errors"
	"net/url"
	"path/filepath"

	"github.com/rs/zerolog/log"
)

func RunScript(d *deploy.Deploy, comp *types.Component, scriptAlias []string, script *types.Script, scriptInput types.Values) (output types.Values, err error) {
	l := log.With().Strs("scriptAlias", scriptAlias).Logger()

	if theDependencyWasDeployedEarlier, ok := d.Root.Values[scriptAlias[len(scriptAlias)-1]]; ok {
		l.Debug().Msg("Deployed earlier")
		return theDependencyWasDeployedEarlier.(types.Values), nil
	}

	if script == nil {
		script = &types.Script{}
	}

	if script.Values == nil {
		script.Values = make(types.Values)
	}
	if _, ok := script.Values["resource"]; !ok {
		script.Values["resource"] = scriptAlias[len(scriptAlias)-1]
	}
	// script.Path prepare for go-git function
	// remove Schema from path
	u, err := url.Parse(script.Path)
	if err != nil {
		log.Error().Err(err).Msg("Url parsing  error")
		return nil, err
	}
	script.Path = u.Host + u.Path

	l.Info().Str("type", script.Type).Str("repository", script.Repository).Str("path", script.Path).Bool("root", script.Shared).Msg("RunScript")

	repositoryFS := comp.Status.Filesystem
	workdir := filepath.Dir(comp.Status.FileName)
	if script.Repository != "" {
		if filesystem.Supported(script.Repository) {
			repositoryFS = filesystem.Get(script.Repository)
			workdir = script.Path
		} else {
			script.Values["repository"] = script.Repository
			script.Values["path"] = script.Path
		}
	} else {
		workdir = d.FS.FS.Join(workdir, script.Path)
	}
	deployChildForScript := d.Child(repositoryFS, workdir, scriptAlias)

	context, err := prepareInput(deployChildForScript, script.Values, d.Root.Values, scriptInput, scriptAlias)
	if err != nil {
		l.Error().Err(err).Msg("templating error")
		return nil, err
	}

	deployChildForScript.Keeper.Push("context", context)
	deployChildForScript.Keeper.Push("script", script)

	if adapter, ok := deploy.DefaultAdapters[script.Type]; ok {
		output, err = adapter(deployChildForScript, script, context)
		if err != nil {
			l.Error().Err(err).Msg("RunScriptFuncImplementation error")
			return nil, err
		}

		deployChildForScript.Keeper.Push("outputBeforeMapping", output)

		output, err = prepareOutput(deployChildForScript, script.OutputMapping, output, context, scriptInput, scriptAlias)
		if err != nil {
			l.Error().Err(err).Msg("prepareOutput error")
			return nil, err
		}
		deployChildForScript.Keeper.Push("output", output)

		//TODO подумать, возможноли и нужно ли избегать безконечного цикла, когда в компоненте вызывается зависимость на саму себя (возможно неявно через цепочку)
		//например, добавить в начало Component(...), счётчик вызовов определённого пути, и не допускать вызова более 10 раз
		if script.Shared {
			deployChildForScript.Root.Values[scriptAlias[len(scriptAlias)-1]] = output
		}

		l.Info().Any("output", output).Msg("RunScript - result")
		return output, err
	}

	l.Error().Str("scriptType", script.Type).Msg("RunScript function not found")
	return nil, errors.New("RUNSCRIPT FUNCTION NOT FOUND")
}

func prepareInput(d *deploy.Deploy, scriptValues, appContext, scriptContext types.Values, aliases []string) (types.Values, error) {
	fullContext := types.MergeValues(
		appContext,
		scriptContext,
		types.Values{
			"aliases": aliases,
			"alias":   deploy.BuildAlias(aliases),
		},
	)

	templated, err := d.Templater.Templating(scriptValues, fullContext)
	if err != nil {
		log.Error().Err(err).Strs("aliases", aliases).Msg("prepareValues error")
		return nil, err
	}
	result := types.MergeValues(fullContext, templated)
	return result, nil
}

func prepareOutput(d *deploy.Deploy, outputMapping, output, context, scriptInput types.Values, aliases []string) (types.Values, error) {
	fullContext := types.MergeValues(
		context,
		output,
		types.Values{
			"aliases": aliases,
			"alias":   deploy.BuildAlias(aliases),
			"input":   scriptInput,
		},
	)

	templated, err := d.Templater.Templating(outputMapping, fullContext)
	if err != nil {
		log.Error().Err(err).Strs("aliases", aliases).Msg("prepareValues error")
		return nil, err
	}
	result := types.MergeValues(output, templated)
	return result, nil
}
