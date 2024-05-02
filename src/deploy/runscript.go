package deploy

import (
	"deployto/src/types"
	"errors"
	"strings"

	"github.com/rs/zerolog/log"
)

type RunScriptFuncImplementationType = func(target *types.Target, workdir string, aliases []string, rootValues, values types.Values) (output map[string]any, err error)

var RunScriptFuncImplementations = map[string]RunScriptFuncImplementationType{}

func RunScript(target *types.Target, workdir string, aliases []string, rootOutput types.Values, script *types.Script, scriptContext types.Values) (output types.Values, err error) {
	l := log.With().Strs("aliases", aliases).Logger()

	if theDependencyWasDeployedEarlier, ok := rootOutput[buildAlias(aliases)]; ok {
		l.Info().Strs("alias", aliases).Msg("Deployed earlier")
		return theDependencyWasDeployedEarlier.(types.Values), nil
	}

	l.Debug().Any("script", script).Msg("RunScript")
	//TODO если script.Type не указан, но указаны repository и/или path, то определять что находится по этом пути, и задавать script.Type автоматом

	input := lookupValues(script.Values, scriptContext)

	if RunScriptFuncImplementation, ok := RunScriptFuncImplementations[script.Type]; ok {
		output, err = RunScriptFuncImplementation(target, workdir,
			aliases,
			rootOutput, input)

		//TODO подумать, возможноли и нужно ли избегать безконечного цикла, когда в компоненте вызывается зависимость на саму себя
		//например, добавить в начало Component(...), счётчик вызовов определённого пути, и не допускать вызова более 10 раз
		if script.Root {
			rootOutput[buildAlias(aliases)] = output
		}

		return output, err
	}
	l.Error().Str("scriptType", script.Type).Msg("RunScript function not found")
	return nil, errors.New("RUNSCRIPT FUNCTION NOT FOUND")
}

func lookupValues(values types.Values, scriptContext types.Values) types.Values {
	if values == nil {
		return nil
	}
	result := make(types.Values, len(values))

	for k, v := range values {
		if v, ok := v.(types.Values); ok {
			result[k] = lookupValues(v, scriptContext)
			continue
		}
		if deploytoStr, ok := v.(string); ok {
			if deploytoStr, ok = strings.CutPrefix(deploytoStr, "__deployto-lookup:"); ok {
				deploytoStr = strings.Trim(deploytoStr, " ")
				result[k] = types.Get(scriptContext, "", deploytoStr)
				continue
			}
		}
		result[k] = v
	}
	return result
}
