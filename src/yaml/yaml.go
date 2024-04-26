package yaml

import (
	"deployto/src/types"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/rs/zerolog/log"
)

const DeploytoPath = ".deployto" //TODO bug this const defined more then in one place

func GetDeploytoPath(path string) string {
	return filepath.Join(path, DeploytoPath)
}

func SubPathExists(path string, subpath string) bool {
	fi, err := os.Stat(filepath.Join(path, subpath))
	if err == nil {
		return fi.IsDir()
	}
	if !os.IsNotExist(err) {
		log.Info().Err(err).Str("path", filepath.Join(path, subpath)).Msg("app path search error")
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

		if SubPathExists(currentPath, searchDir) {
			return filepath.Join(currentPath, searchDir), nil
		}
		currentPath = filepath.Dir(currentPath)
	}
}

// Ищу компоненты.
// Начиная с указанной, проверяю все родительские папки на наличие в ней папки .deployto, когда найдена, то это и есть папка компоненты
// Возвращаю все kind: Component из папки компоненты
func GetComponent(path string) (comps []*types.Component, err error) {
	rootPath, err := GetProjectRoot(path, DeploytoPath)
	if err != nil {
		log.Error().Str("startPath", path).Msg("Project dir not found")
		return nil, err
	}

	comps = Get[types.Component](rootPath)
	if len(comps) == 0 {
		log.Error().Str("startPath", path).Str("currentPath", rootPath).Msg("Component not found")
		return nil, errors.New("COMPONENT NOT FOUND IN COMPONENT FOLDER")
	}
	return comps, nil
}

func Get[T types.Component | types.Environment | types.Target | types.Job](deploytoPath string) (result []*T) {
	if !IsDeploytoPath(deploytoPath) {
		deploytoPath = GetDeploytoPath(deploytoPath)
	}
	err := filepath.Walk(deploytoPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && strings.HasSuffix(strings.ToLower(path), ".yaml") {
				file, err := os.Open(path)
				if err != nil {
					log.Error().Str("file", path).Err(err).Msg("Open file error")
					return err
				}
				dec := yaml.NewDecoder(file)
				for {
					var item any = new(T)
					err := dec.Decode(item)
					if err != nil {
						if err.Error() == "EOF" {
							break
						}
						log.Error().Str("file", path).Err(err).Msg("yaml decode error")
						if strings.HasPrefix(err.Error(), "yaml: line ") {
							break
						}
						continue
					}

					switch itemTyped := item.(type) {
					case *types.Component:
						itemTyped.Base.Status.FileName = path
						if itemTyped.Kind == "Component" {
							result = append(result, item.(*T))
						}
					case *types.Environment:
						itemTyped.Base.Status.FileName = path
						if itemTyped.Kind == "Environment" {
							result = append(result, item.(*T))
						}
					case *types.Target:
						itemTyped.Base.Status.FileName = path
						if itemTyped.Kind == "Target" {
							result = append(result, item.(*T))
						}
					case *types.Job:
						if itemTyped.Kind == "Job" {
							result = append(result, item.(*T))
						}
					default:
						log.Error().Type("type", item).Msg("yaml crd type not supported")
					}
				}
			}
			return nil
		})
	if err != nil {
		log.Error().Err(err).Msg("yaml reading error")
	}
	return
}
