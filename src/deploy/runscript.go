package deploy

import "deployto/src/types"

type RunScript = func(kubeconfig string, names []string, kind string, script *types.Script, input map[string]any) (output map[string]any, err error)

var RunScripts = map[string]RunScript{}
