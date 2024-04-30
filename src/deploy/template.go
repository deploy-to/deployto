package deploy

import (
	"deployto/src/types"
	"strings"

	"github.com/rs/zerolog/log"
)

func init() {
	RunScriptFuncImplementations["template"] = Template
	RunScriptFuncImplementations[""] = Template //default script type
}
func Template(kubeconfig string, workdir string, aliases []string, rootValues, input types.Values) (output types.Values, err error) {
	output = types.Values{
		"ConnectionString": "http://" + strings.Join(aliases, "."),
	}
	log.Error().Strs("aliases", aliases).Any("input", input).Any("output", output).Msg("Заглушка, для деплоя Template")
	return output, nil
}
