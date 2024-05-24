package deploy

import "deployto/src/types"

type AdapterImplementation = func(d *Deploy, script *types.Script, scriptContext types.Values) (output map[string]any, err error)

var DefaultAdapters = map[string]AdapterImplementation{}
