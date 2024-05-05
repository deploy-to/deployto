package yaml

import (
	"deployto/src/filesystem"
	"deployto/src/types"
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

func Get[T types.Component | types.Environment | types.Target | types.Job](fs *filesystem.Filesystem, deploytoDir string) (result []*T) {
	files, err := fs.FS.ReadDir(deploytoDir)
	if err != nil {
		log.Error().Err(err).Msg("ReadDir error")
	}
	for _, file := range files {
		// if file.IsDir() && file.Name() == filesystem.DeploytoDirName {
		// 	result = append(result, Get[T](fs, file.Name())...)
		// }
		if !file.IsDir() && strings.HasSuffix(strings.ToLower(file.Name()), ".yaml") {
			result = append(result, GetFromFile[T](fs, fs.FS.Join(deploytoDir, file.Name()))...)
		}
	}
	return
}

func GetFromFile[T types.Component | types.Environment | types.Target | types.Job](fs *filesystem.Filesystem, fileName string) (result []*T) {
	file, err := fs.FS.Open(fileName)
	if err != nil {
		log.Error().Str("file", fileName).Err(err).Msg("Open file error")
		return nil
	}
	dec := yaml.NewDecoder(file)
	for {
		var item any = new(T)
		err := dec.Decode(item)
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			log.Error().Str("file", fileName).Err(err).Msg("yaml decode error")
			if strings.HasPrefix(err.Error(), "yaml: line ") {
				break
			}
			continue
		}

		switch itemTyped := item.(type) {
		case *types.Component:
			itemTyped.Base.Status.Filesystem = fs
			itemTyped.Base.Status.FileName = fileName
			if itemTyped.Kind == "Component" {
				result = append(result, item.(*T))
			}
		case *types.Environment:
			itemTyped.Base.Status.Filesystem = fs
			itemTyped.Base.Status.FileName = fileName
			if itemTyped.Kind == "Environment" {
				result = append(result, item.(*T))
			}
		case *types.Target:
			itemTyped.Base.Status.Filesystem = fs
			itemTyped.Base.Status.FileName = fileName
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
