package deploy

import (
	"deployto/src/types"

	"github.com/lithammer/shortuuid/v3"
	"github.com/rs/zerolog/log"
)

func init() {
	RunScriptFuncImplementations["job"] = Job
}

func Job(kubeconfig string, workdir string, aliases []string, rootValues, input types.Values) (output types.Values, err error) {
	log.Error().Msg("Заглушка, для выполнения Job")
	uuid := shortuuid.New()
	return map[string]any{
		"image":      "ima-dummy" + uuid,
		"repository": "rep-dummy" + uuid,
		"tag":        "tag-dummy" + uuid,
	}, nil
}
