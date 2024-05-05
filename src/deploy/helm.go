package deploy

import (
	"deployto/src/filesystem"
	"deployto/src/types"
	"strings"

	"github.com/rs/zerolog/log"
)

func init() {
	RunScriptFuncImplementations["helm"] = Helm
}

func Helm(target *types.Target, fs *filesystem.Filesystem, workDir string, aliases []string, rootValues, input types.Values) (output types.Values, err error) {
	output = types.Values{
		"ConnectionString": "http://" + strings.Join(aliases, "."),
	}
	log.Error().Strs("aliases", aliases).Any("input", input).Any("output", output).Msg("Заглушка, для деплоя Helm")
	return output, nil
}
