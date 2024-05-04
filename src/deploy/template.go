package deploy

import (
	"deployto/src"
	"deployto/src/filesystem"
	"deployto/src/types"
	"deployto/src/yaml"
	"errors"
	"os"

	"github.com/rs/zerolog/log"
)

func init() {
	RunScriptFuncImplementations["template"] = Template
	RunScriptFuncImplementations[""] = Template //default script type
}

func Template(target *types.Target, fs *filesystem.Filesystem, aliases []string, rootValues, input types.Values) (output types.Values, err error) {
	selector := types.DecodeTemplateArg(input)
	return runTemplateTyped(target, fs, aliases, rootValues, selector)
}

func runTemplateTyped(target *types.Target, fs *filesystem.Filesystem, aliases []string, rootValues types.Values, selector *types.TemplateArg) (output types.Values, err error) {
	log.Debug().Strs("aliases", aliases).Msg("Search template")
	template := searchTemplate(selector)
	log.Debug().Str("templateDir", template.Status.FileName).Msg("found template")
	//var similars []*types.Component
	//TODO cache
	// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	return RunSingleComponent(target, fs, aliases, rootValues, selector.Values, template)
}

func searchTemplate(selector *types.TemplateArg) *types.Component {
	var repositories []*filesystem.Filesystem
	for _, r := range src.GetTemplateRepositories() {
		repositories = append(repositories, filesystem.GetFilesystem(r))
	}
	if len(repositories) == 0 {
		log.Warn().Msg("Template repositories not found")
	}
	paths := []string{
		repositories[0].FS.Join("/resources", selector.Resource, selector.Version, selector.Name+".yaml"),
		repositories[0].FS.Join("/resources", selector.Resource, selector.Version, "defaul.yaml"),
		repositories[0].FS.Join("/resources", selector.Resource, "defaul.yaml"),
	}
	for _, path := range paths {
		for _, r := range repositories {
			if c := tryGet(r, path); c != nil {
				return c
			}
		}
	}
	return nil
}

func tryGet(fs *filesystem.Filesystem, path string) *types.Component {
	_, err := fs.FS.Stat(path)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			log.Error().Err(err).Str("path", path).Msg("filesystem.Stat error")
		}
		return nil
	}
	comps := yaml.GetFromFile[types.Component](fs, path)
	if len(comps) == 0 {
		log.Error().Str("path", path).Msg("file exists, but component not found")
	}
	if len(comps) > 1 {
		log.Error().Str("path", path).Msg("only one component allowed")
	}
	return nil
}
