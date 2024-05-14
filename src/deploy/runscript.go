package deploy

import (
	"bytes"
	"deployto/src"
	"deployto/src/filesystem"
	"deployto/src/types"
	"errors"
	"html/template"
	"net/url"

	"github.com/rs/zerolog/log"
)

type RunScriptFuncImplementationType = func(target *types.Target, fs *filesystem.Filesystem, workDir string, aliases []string, rootValues, values types.Values, ContextDump *src.ContextDump) (output map[string]any, err error)

var RunScriptFuncImplementations = map[string]RunScriptFuncImplementationType{}

func RunScript(target *types.Target, repositoryFS *filesystem.Filesystem, workdir string, aliases []string, script *types.Script, rootContext, parentContext types.Values, ContextDump *src.ContextDump) (output types.Values, err error) {
	l := log.With().Strs("aliases", aliases).Logger()

	if theDependencyWasDeployedEarlier, ok := rootContext[buildAlias(aliases)]; ok {
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
		script.Values["resource"] = aliases[len(aliases)-1]
	}
	// script.Path prepare for go-git function
	// remove Schema from path
	u, err := url.Parse(script.Path)
	if err != nil {
		log.Error().Err(err).Msg("Url parsing  error")
		return nil, err
	}
	script.Path = u.Host + u.Path
	if script.Repository != "" {
		if filesystem.Supported(script.Repository) {
			repositoryFS = filesystem.Get(script.Repository)
			workdir = script.Path
		} else {
			script.Values["repository"] = script.Repository
			script.Values["path"] = script.Path
		}
	} else {
		workdir = repositoryFS.FS.Join(workdir, script.Path)
	}

	l.Info().Str("type", script.Type).Str("repository", script.Repository).Str("path", script.Path).Bool("root", script.Root).Msg("RunScript")

	context, err := prepareInput(script.Values, rootContext, parentContext, aliases)
	if err != nil {
		l.Error().Err(err).Msg("templating error")
		return nil, err
	}
	l.Info().Any("values", script.Values).Any("rootOutput", rootContext).Any("scriptContext", parentContext).Any("input", context).Msg("RunScript - values")

	ContextDump.Push("context", context)
	ContextDump.Push("script", script)

	if RunScriptFuncImplementation, ok := RunScriptFuncImplementations[script.Type]; ok {
		output, err = RunScriptFuncImplementation(target, repositoryFS, workdir,
			aliases,
			rootContext, context, ContextDump)
		if err != nil {
			l.Error().Err(err).Msg("RunScriptFuncImplementation error")
			return nil, err
		}

		output, err = prepareOutput(script.OutputMapping, output, context, aliases)
		if err != nil {
			l.Error().Err(err).Msg("prepareOutput error")
			return nil, err
		}
		ContextDump.Push("output", output)

		//TODO подумать, возможноли и нужно ли избегать безконечного цикла, когда в компоненте вызывается зависимость на саму себя (возможно неявно через цепочку)
		//например, добавить в начало Component(...), счётчик вызовов определённого пути, и не допускать вызова более 10 раз
		if script.Root {
			rootContext[buildAlias(aliases)] = output
		}

		l.Info().Any("output", output).Msg("RunScript - result")
		return output, err
	}

	l.Error().Str("scriptType", script.Type).Msg("RunScript function not found")
	return nil, errors.New("RUNSCRIPT FUNCTION NOT FOUND")
}

func mergeContext(rootContext, parentContext types.Values, aliases []string) map[string]any {
	return types.MergeValues(
		rootContext,
		parentContext,
		types.Values{
			"aliases": aliases,
			"alias":   buildAlias(aliases),
		},
	)
}

func prepareInput(values, rootContext, parentContext types.Values, aliases []string) (types.Values, error) {
	fullContext := mergeContext(rootContext, parentContext, aliases)
	templated, err := templating(values, fullContext)
	if err != nil {
		log.Error().Err(err).Strs("aliases", aliases).Msg("prepareValues error")
		return nil, err
	}
	result := types.MergeValues(fullContext, templated)
	return result, nil
}

func prepareOutput(outputMapping, output, context types.Values, aliases []string) (types.Values, error) {
	templated, err := templating(outputMapping, types.Values{"output": output, "context": context})
	if err != nil {
		log.Error().Err(err).Strs("aliases", aliases).Msg("prepareValues error")
		return nil, err
	}
	result := types.MergeValues(output, templated)
	return result, nil
}

func templating(values, context types.Values) (types.Values, error) {
	result := make(types.Values)
	for k, v := range values {
		switch vTyped := v.(type) {
		case types.Values:
			subResult, err := templating(vTyped, context)
			if err != nil {
				log.Error().Err(err).Str("key", k).Msg("Template subValues execute with scriptContext error")
				return nil, err
			}
			result[k] = subResult
		case string:
			res, err := templatingString(vTyped, context)
			if err != nil {
				log.Error().Err(err).Str("key", k).Str("template", vTyped).Msg("Templating error")
				return nil, err
			}
			result[k] = res
		default:
			result[k] = v
		}
	}
	return result, nil
}

var funcMap = template.FuncMap{
	"inc": func(i int) int { return i + 1 },
	"dec": func(i int) int { return i + 1 },
	"add": func(i1 int, i2 int) int { return i1 + i2 },
	"sub": func(i1 int, i2 int) int { return i1 - i2 },
}

func templatingString(templ string, context types.Values) (string, error) {
	t, err := template.New("letter").Funcs(funcMap).Parse(templ)
	if err != nil {
		log.Error().Err(err).Str("template", templ).Msg("Template parse error")
		return "nil", err
	}
	buf := new(bytes.Buffer)
	err = t.Execute(buf, context)
	if err != nil {
		log.Error().Err(err).Str("template", templ).Msg("Template execute with scriptContext error")
		return "", err
	}
	return buf.String(), nil
}
