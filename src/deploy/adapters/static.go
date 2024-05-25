package adapters

import (
	"deployto/src/deploy"
	"deployto/src/types"
)

func init() {
	deploy.DefaultAdapters["static"] = Static
}

func Static(d *deploy.Deploy, script *types.Script, scriptContext types.Values) (output types.Values, err error) {
	return nil, nil
}
