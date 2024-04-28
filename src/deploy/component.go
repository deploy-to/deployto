package deploy

import (
	"deployto/src/gitclient"
	"deployto/src/types"
	"deployto/src/yaml"
	"path/filepath"
	"strings"

	"github.com/lithammer/shortuuid/v3"
	"github.com/rs/zerolog/log"
)

func init() {
	RunScripts["component"] = Component
}
func Component(kubeconfig string, workdir string, name string, alias string, aliases []string, rootValues, input types.Values) (output types.Values, err error) {
	// COMPONENTS
	comps, err := yaml.GetComponent(workdir)
	if err != nil {
		log.Error().Err(err).Str("path", workdir).Msg("Components search error")
		return nil, err
	}
	output = make(types.Values)
	for _, c := range comps {
		l := log.With().Strs("aliases", aliases).Logger()
		output = make(types.Values)

		dependencies := types.Get(input, []types.Dependency(nil), "dependencies")
		for _, d := range dependencies {
			var dependencyAliases []string
			if !d.Root {
				dependencyAliases = aliases
			}
			if d.Name == "" {
				dependencyAliases = append(dependencyAliases, d.Kind)
			} else {
				dependencyAliases = append(dependencyAliases, d.Name)
			}

			if _, ok := rootValues[buildAlias(aliases)]; ok {
				l.Info().Strs("alias", dependencyAliases).Msg("Deployed earlier")
				continue
			}

			// dependencyValues := map[string]any{
			// 	"name":    d.Name,
			// 	"alias":   buildAlias(aliases),
			// 	"aliases": dependencyAliases,
			// }

			// for dependencies, if the script is not defined, I will try to get the default script by kind
			l.Debug().Strs("alias", dependencyAliases).Msg("Run dependency")

			//!!!o, e := RunScript(kubeconfig, dependencyAliases, "", nil, dependencyValues)
			// if e != nil {
			// 	l.Error().Err(e).Msg("RunScript error")
			// 	return output, err
			// }
			// if _, ok := output["dependencies"]; ok {
			// 	output["dependencies"].(types.Values)[buildAlias(aliases)] = o
			// } else {
			// 	dependencies := types.Values{}
			// 	output["dependencies"] = dependencies
			// }

		}

		var scriptWorkdir string
		repository := types.Get(input, "", "script.repository")
		path := types.Get(input, "", "script.path")
		if repository == "" || strings.HasPrefix(repository, "file://") {
			if repository == "" {
				scriptWorkdir = filepath.Join(workdir, path)
			} else {
				scriptWorkdir = filepath.Join(workdir, repository[len("file://"):], path)
			}
		} else {
			panic("not implimented")
		}

		o := make(types.Values)
		l.Debug().Msg("Get commit hash and tags")
		o["git"] = gitclient.GetValues(scriptWorkdir)
		l.Debug().Msg("TODO  BUILD && Push")
		o["build"], err = DoJob(scriptWorkdir, types.Get(input, "", "build.job"), input)

		c.Spec = types.MergeValues(output, c.Spec)

		o, e := RunScript(kubeconfig, scriptWorkdir,
			types.Get(c.Spec, "", "name"), types.Get(c.Spec, "", "alias"), append(aliases, alias),
			rootValues, c.Spec)
		if e != nil {
			l.Error().Err(e).Msg("RunScript error")
		}
		output[types.Get(c.Spec, "", "alias")] = o
	}

	return output, err
}

func buildAlias(names []string) string {
	return strings.Join(names, "-")
}

// TODO заглушка, для вызова job
// output где взять образ
func DoJob(workdir string, jobName string, values types.Values) (map[string]any, error) {
	uuid := shortuuid.New()
	return map[string]any{
		"image":      "ima-dummy" + uuid,
		"repository": "rep-dummy" + uuid,
		"tag":        "tag-dummy" + uuid,
	}, nil
}
