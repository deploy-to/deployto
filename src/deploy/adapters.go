package deploy

import "deployto/src/types"

/*
 =========   workflows   =========

| Deployto | Terraform | Helm      |
| -------- | --------- | --------- |
|          | init      | install   |
|          | validate  |           |
|          | plan      |           |
| apply    | apply     | upgrade   |
|          |           | rollback  |
| destroy  | destroy   | uninstall |

*/

type Adapter interface {
	Apply(d *Deploy, script *types.Script, scriptContext types.Values) (output map[string]any, err error)
	Destroy(d *Deploy, script *types.Script, scriptContext types.Values) error
}

var DefaultAdapters = map[string]Adapter{}
