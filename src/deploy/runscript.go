package deploy

import (
	"deployto/src/types"
	"errors"

	"github.com/rs/zerolog/log"
)

type RunScriptFunc = func(kubeconfig string, workdir string, name string, alias string, aliases []string, rootValues, input types.Values) (output map[string]any, err error)

var RunScripts = map[string]RunScriptFunc{}

func RunScript(kubeconfig string, workdir string, name string, alias string, aliases []string, rootValues, input types.Values) (output types.Values, err error) {
	l := log.With().Strs("names", aliases).Logger()
	if !types.Exists(input, "script") {
		l.Debug().Any("input", input).Msg("Script not defined")
		return nil, nil
	}
	scriptType := types.Get(input, "component", "script.type")
	l.Debug().Str("scriptType", scriptType).Any("input", input).Msg("RunScript")

	//repository!!!!!!!!!

	//alias := types.Get(input, "", "script.alias")

	if runScript, ok := RunScripts[scriptType]; ok {
		return runScript(kubeconfig, workdir,
			types.Get(input, "", "script.name"), alias, append(aliases, alias),
			rootValues, input)
	}
	l.Error().Str("scriptType", scriptType).Msg("RunScript function not found")
	return nil, errors.New("RUNSCRIPT FUNCTION NOT FOUND")
}
