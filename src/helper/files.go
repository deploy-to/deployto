package helper

import (
	"deployto/src"
)

const DeploytoPath = ".deployto" //TODO bug this const defined more then in one place

func GetDeploytoPath(fs *src.Filesystem, path string) string {
	return fs.FS.Join(path, DeploytoPath)
}
