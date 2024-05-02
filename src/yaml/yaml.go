package yaml

import (
	"deployto/src/helper"
	"deployto/src/types"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/rs/zerolog/log"
)

// Ищу компоненты.
// Начиная с указанной, проверяю все родительские папки на наличие в ней папки .deployto, когда найдена, то это и есть папка компоненты
// Возвращаю все kind: Component из папки компоненты
func GetComponent(path string) (comps []*types.Component, err error) {
	rootPath, err := helper.GetProjectRoot(path, helper.DeploytoPath)
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
	if !helper.IsDeploytoPath(deploytoPath) {
		deploytoPath = helper.GetDeploytoPath(deploytoPath)
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
					case *types.Service:
						if itemTyped.Kind == "Service" {
							result = append(result, item.(*T))
						}
					case *types.Ingress:
						if itemTyped.Kind == "Ingress" {
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

func GetBytes[T types.Service | types.Ingress](yamlb []byte) (result []*T) {
	var yl map[string]interface{}
	err := yaml.Unmarshal(yamlb, &yl)
	if err != nil {
		log.Error().Err(err).Str("path", "yaml").Msg("Parse yaml error")
	}

	var item any = new(T)
	switch itemTyped := item.(type) {
	case *types.Ingress:
		if itemTyped.Kind == "Ingress" {
			result = append(result, item.(*T))
		}
	case *types.Service:
		if itemTyped.Kind == "Service" {
			result = append(result, item.(*T))
		}
	default:
		log.Error().Type("type", item).Msg("yaml crd type not supported")
	}
	return
}

func GetBytes2(yamlb []byte) (services []types.Service, ingresses []types.Ingress) {

	// Convert Kubernetes objects to YAML
	err := yaml.Unmarshal(yamlb, &services)
	if err != nil {
		fmt.Println("Error converting services to YAML:", err)
		return
	}

	err = yaml.Unmarshal(yamlb, &ingresses)
	if err != nil {
		fmt.Println("Error converting ingresses to YAML:", err)
		return
	}
	return
}
