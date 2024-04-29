package deploy

import (
	"deployto/src/types"

	"github.com/rs/zerolog/log"
)

func init() {
	RunScriptFuncImplementations["template"] = Template
}
func Template(kubeconfig string, workdir string, aliases []string, rootValues, input types.Values) (output types.Values, err error) {
	log.Error().Msg("Заглушка, для выполнения Template")
	return nil, nil
}
