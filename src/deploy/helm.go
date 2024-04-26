package deploy

import "deployto/src/types"

func init() {
	RunScripts["helm"] = HelmRunScript
}

func HelmRunScript(kubeconfig string, names []string, kind string, script *types.Script, input map[string]any) (output map[string]any, err error) {
	// эта функци будет вызыватсья только для script.type = helm
	// для script.type == helm, атрибут kind можно игнорировать
	return nil, nil
}
