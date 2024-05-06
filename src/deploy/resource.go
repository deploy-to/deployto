package deploy

import (
	"deployto/src/filesystem"
	"deployto/src/types"
	"deployto/src/yaml"
	"errors"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

func init() {
	RunScriptFuncImplementations["resource"] = Resource
	RunScriptFuncImplementations[""] = Resource //default script type
}

func Resource(target *types.Target, fs *filesystem.Filesystem, workDir string, aliases []string, rootValues, input types.Values) (output types.Values, err error) {
	selector := types.DecodeResourceArg(input)
	return runResourceTyped(target, fs, aliases, rootValues, selector)
}

func runResourceTyped(target *types.Target, _ *filesystem.Filesystem, aliases []string, rootValues types.Values, selector *types.ResourceArg) (output types.Values, err error) {
	log.Debug().Strs("aliases", aliases).Msg("Search template")
	template := searchResourceInRepositories(selector)
	log.Debug().Str("templateDir", template.Status.FileName).Msg("found template")
	//var similars []*types.Component
	//TODO cache
	// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	return RunSingleComponent(target, aliases, rootValues, selector.Values, template)
}

func searchResourceInRepositories(selector *types.ResourceArg) *types.Component {
	if selector.Resource == "" {
		log.Debug().Msg("selector.Resource not set")
		return nil
	}
	var repositories []*filesystem.Filesystem
	for _, r := range GetTemplateRepositories() {
		repositories = append(repositories, filesystem.Get(r))
	}
	if len(repositories) == 0 {
		log.Warn().Msg("Template repositories not found")
	}

	var path string
	if selector.Version == "" && selector.Name == "" {
		path = repositories[0].FS.Join("resources", selector.Resource, "default.yaml")
	} else {
		if selector.Version != "" && selector.Name != "" {
			path = repositories[0].FS.Join("resources", selector.Resource, selector.Version, selector.Name+".yaml")
		} else {
			log.Error().Str("Resource", selector.Resource).Str("Version", selector.Version).Str("Name", selector.Name).Msg("Resource.Version && Resource.Name - Set both or nothing")
			return nil
		}
	}

	for _, r := range repositories {
		if comp := tryGet(r, path); comp != nil {
			return comp
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

var TemplateRepositories cli.StringSlice //setup in cli.StringSliceFlag{ Name: "templateRepositories", Destination: &deploy.TemplateRepositories,

func GetTemplateRepositories() (result []string) {
	for _, r := range TemplateRepositories.Value() {
		if filesystem.Supported(r) {
			result = append(result, r)
		} else {
			log.Error().Str("repository", r).Msg("Unsupport repository")
		}
	}
	result = append(result, "git@github.com:deploy-to/deployto.git")
	return result
}
