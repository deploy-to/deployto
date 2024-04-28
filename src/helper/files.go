package helper

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
)

const DeploytoPath = ".deployto" //TODO bug this const defined more then in one place

func GetDeploytoPath(path string) string {
	return filepath.Join(path, DeploytoPath)
}

func IsSubPathExists(path string, subpath string) bool {
	return IsDirExists(filepath.Join(path, subpath))
}

func IsDirExists(path string) bool {
	fi, err := os.Stat(path)
	if err == nil {
		return fi.IsDir()
	}
	if !os.IsNotExist(err) {
		log.Info().Err(err).Str("path", path).Msg("Check dir exists")
	}
	return false
}

func IsDeploytoPath(path string) bool {
	_, file := filepath.Split(path)
	return file == DeploytoPath || filepath.Base(path) == DeploytoPath
}

func GetProjectRoot(path string, searchDir string) (string, error) {
	currentPath := path
	for {
		if currentPath == "/" || len(currentPath) < 4 /*TODO need test on windows*/ {
			log.Error().Str("startPath", path).Str("currentPath", currentPath).Msg("getDeployToPaths end - too short path")
			return "", errors.New("ROOT FOLDER NOT FOUND")
		}

		log.Debug().Str("path", currentPath).Msg("check dir")

		if IsSubPathExists(currentPath, searchDir) {
			return filepath.Join(currentPath, searchDir), nil
		}
		currentPath = filepath.Dir(currentPath)
	}
}
