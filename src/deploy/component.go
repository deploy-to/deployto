package deploy

import (
	"deployto/src/types"
)

func init() {
	RunScripts["component"] = ComponentRunScript
}

func ComponentRunScript(kubeconfig string, names []string, kind string, script *types.Script, input map[string]any) (output map[string]any, err error) {
	//script.repository установлен - если нет, то использовать дефолтный
	return nil, nil
}
