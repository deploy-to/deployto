package deploy

import "deployto/src/types"

/*
 =========   workflows   =========

| Deployto Adapters | Terraform | Helm      |
| ----------------- | --------- | --------- |
|                   | init      | install   |
|                   | validate  |           |
|                   | plan      |           |
| apply             | apply     | upgrade   |
|                   |           | rollback  |
| destroy           | destroy   | uninstall |

*/

type Adapter interface {
	Apply(d *DeployState, script *types.Script, scriptContext types.Values) (output map[string]any, err error)
	Destroy(d *DeployState, script *types.Script, scriptContext types.Values) error
}

var DefaultAdapters = map[string]Adapter{}
