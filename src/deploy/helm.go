package deploy

import (
	"deployto/src/types"

	"github.com/rs/zerolog/log"
)

func init() {
	RunScriptFuncImplementations["helm"] = Helm
}

func Helm(kubeconfig string, workdir string, aliases []string, rootValues, input types.Values) (output types.Values, err error) {
	log.Error().Msg("Заглушка, для деплоя Helm")
	return nil, nil
}
