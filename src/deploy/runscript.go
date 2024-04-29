package deploy

import (
	"deployto/src/types"
	"errors"
	"strings"

	"github.com/rs/zerolog/log"
)

type RunScriptFuncImplementationType = func(kubeconfig string, workdir string, aliases []string, rootValues, values types.Values) (output map[string]any, err error)

var RunScriptFuncImplementations = map[string]RunScriptFuncImplementationType{}

func RunScript(kubeconfig string, workdir string, aliases []string, rootValues, script, scriptContext types.Values) (output types.Values, err error) {
	l := log.With().Strs("aliases", aliases).Logger()
	if script == nil {
		l.Debug().Msg("Script not defined")
		return nil, nil
	}
	scriptType := types.Get(script, "template", "type")
	l.Debug().Str("scriptType", scriptType).Any("input", script).Msg("RunScript")

	input := lookupValues(types.Get(script, types.Values(nil), "values"), scriptContext)

	if runScript, ok := RunScriptFuncImplementations[scriptType]; ok {
		return runScript(kubeconfig, workdir,
			aliases,
			rootValues, input)
	}
	l.Error().Str("scriptType", scriptType).Msg("RunScript function not found")
	return nil, errors.New("RUNSCRIPT FUNCTION NOT FOUND")
}

func lookupValues(scripValues, scriptContext types.Values) types.Values {
	if scripValues == nil {
		return nil
	}
	result := make(types.Values, len(scripValues))

	for k, v := range scripValues {
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
