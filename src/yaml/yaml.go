package yaml

import (
	"deployto/src/types"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/rs/zerolog/log"
)

const DeploytoPath = ".deployto"

// Ищу корневую папку приложения. Именно в ней, в .deployto хранятся настройки приложения, окружений, таргетов
func GetAppPath(workPath string) string {
	for {
		if workPath == "/" || len(workPath) < 4 /*TODO need test on windows*/ {
			log.Debug().Str("path", workPath).Msg("getDeployToPaths end - too short path")
			return ""
		}

		log.Debug().Str("path", workPath).Msg("check dir")
		tryDeploytoPath := filepath.Join(workPath, DeploytoPath)
		_, err := os.Stat(tryDeploytoPath)
		if err == nil {
			log.Debug().Str("path", tryDeploytoPath).Msg("deployto path found")
			if len(Get[types.Application](tryDeploytoPath)) > 0 {
				return workPath
			}
		} else {
			if !os.IsNotExist(err) {
				return ""
			}
		}
		workPath = filepath.Dir(workPath)
	}
}

func Get[T types.Application | types.Environment | types.Target](appDeploytoPath string) (result []*T) {
	err := filepath.Walk(appDeploytoPath,
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
						if err.Error() == "EOF" || strings.HasPrefix(err.Error(), "yaml: line ") {
							break
						}
						log.Debug().Str("file", path).Err(err).Msg("yaml decode error")
						continue
					}

					switch itemTyped := item.(type) {
					case *types.Application:
						if itemTyped.Kind == "Application" {
							result = append(result, item.(*T))
						}
					case *types.Environment:
						if itemTyped.Kind == "Envirement" {
							result = append(result, item.(*T))
						}
					case *types.Target:
						if itemTyped.Kind == "Target" {
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
