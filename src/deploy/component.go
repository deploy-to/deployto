package deploy

import "deployto/src/types"

func init() {
	RunScripts["component"] = ComponentRunScript
}

func ComponentRunScript(kubeconfig string, names []string, kind string, script *types.Script, target *types.Target, input map[string]any) (output map[string]any, err error) {
	return nil, nil
}
