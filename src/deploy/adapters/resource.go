package adapters

import (
	"deployto/src/deploy"
	"deployto/src/filesystem"
	"deployto/src/types"
	"deployto/src/yaml"
	"errors"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

func init() {
	deploy.DefaultAdapters["resource"] = (*resource)(nil)
	deploy.DefaultAdapters[""] = (*resource)(nil) //default script type
}

type resource struct{}

func (r *resource) Apply(d *deploy.DeployState, script *types.Script, input types.Values) (output types.Values, err error) {
	log.Debug().Strs("aliases", d.Aliases).Msg("Search resource")
	resource := searchResourceInRepositories(script /*, selector*/)
	if resource == nil {
		log.Error().Any("script", script).Msg("Resource not found")
		return nil, errors.New("Resource not found")
	}
	log.Debug().Str("templateDir", resource.Status.FileName).Msg("found template")
	return ApplySingleComponent(d, resource, input)
}

func (r *resource) Destroy(d *deploy.DeployState, script *types.Script, input types.Values) error {
	panic("NOT IMPLIMENTED")
}

func searchResourceInRepositories(script *types.Script /*, selector *types.ResourceArg*/) *types.Component {
	var repositories []*filesystem.Filesystem
	for _, r := range GetTemplateRepositories() {
		repositories = append(repositories, filesystem.Get(r))
	}
	if len(repositories) == 0 {
		log.Warn().Msg("Template repositories not found")
	}

	// var path string
	// if selector.Version == "" && selector.Name == "" {
	// 	path = repositories[0].FS.Join("resources", script.Name, "default.yaml")
	// } else {
	// 	if selector.Version != "" && selector.Name != "" {
	// 		path = repositories[0].FS.Join("resources", script.Name, selector.Version, selector.Name+".yaml")
	// 	} else {
	// 		log.Error().Str("Resource", script.Name).Str("Version", selector.Version).Str("Name", selector.Name).Msg("Resource.Version && Resource.Name - Set both or nothing")
	// 		return nil
	// 	}
	// }

	path := filepath.Join("resources", script.Name, "default.yaml")
	if script.Path != "" {
		path = filepath.Join("resources", script.Name, script.Path)
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
		return nil
	}
	if len(comps) > 1 {
		log.Error().Str("path", path).Msg("only one component allowed")
		return nil
	}
	return comps[0]
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
