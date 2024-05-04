package yaml

import (
	"deployto/src/filesystem"
	"deployto/src/types"
	"strings"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

func Get[T types.Component | types.Environment | types.Target | types.Job](filesystem *filesystem.Filesystem, deploytoDir string) (result []*T) {
	files, err := filesystem.FS.ReadDir("/")
	if err != nil {
		log.Error().Err(err).Msg("ReadDir error")
	}
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(strings.ToLower(file.Name()), ".yaml") {
			result = append(result, GetFromFile[T](filesystem, file.Name())...)
		}
	}
	return
}

func GetFromFile[T types.Component | types.Environment | types.Target | types.Job](filesystem *filesystem.Filesystem, fileName string) (result []*T) {
	file, err := filesystem.FS.Open(fileName)
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
			itemTyped.Base.Status.FileName = fileName
			if itemTyped.Kind == "Component" {
				result = append(result, item.(*T))
			}
		case *types.Environment:
			itemTyped.Base.Status.FileName = fileName
			if itemTyped.Kind == "Environment" {
				result = append(result, item.(*T))
			}
		case *types.Target:
			itemTyped.Base.Status.FileName = fileName
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
	return
}
