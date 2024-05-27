package adapters

import (
	"deployto/src/deploy"
	"deployto/src/types"
)

func init() {
	deploy.DefaultAdapters["static"] = (*static)(nil)
}

type static struct{}

func (s *static) Apply(d *deploy.DeployState, script *types.Script, scriptContext types.Values) (output types.Values, err error) {
	return nil, nil
}

func (s *static) Destroy(d *deploy.DeployState, script *types.Script, scriptContext types.Values) error {
	return nil
}
