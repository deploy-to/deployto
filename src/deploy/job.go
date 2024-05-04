package deploy

import (
	"deployto/src/filesystem"
	"deployto/src/types"

	"github.com/lithammer/shortuuid/v3"
	"github.com/rs/zerolog/log"
)

func init() {
	RunScriptFuncImplementations["job"] = Job
}

func Job(target *types.Target, fs *filesystem.Filesystem, aliases []string, rootValues, input types.Values) (output types.Values, err error) {
	uuid := shortuuid.New()
	output = types.Values{
		"image":      aliases[len(aliases)-1],
		"repository": "rep-dummy" + uuid,
		"tag":        "tag-dummy" + uuid,
	}
	log.Error().Strs("aliases", aliases).Any("input", input).Any("output", output).Msg("Заглушка, для деплоя Job")
	return output, nil
}
