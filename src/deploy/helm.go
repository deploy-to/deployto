package deploy

import (
	"deployto/src/types"
	"strings"

	"github.com/rs/zerolog/log"
)

func init() {
	RunScriptFuncImplementations["helm"] = Helm
}

func Helm(kubeconfig string, workdir string, aliases []string, rootValues, input types.Values) (output types.Values, err error) {
	output = types.Values{
		"ConnectionString": "http://" + strings.Join(aliases, "."),
	}
	log.Error().Strs("aliases", aliases).Any("input", input).Any("output", output).Msg("Заглушка, для деплоя Helm")
	return output, nil
}
