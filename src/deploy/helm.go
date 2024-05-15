package deploy

import (
	"bytes"
	"context"
	"deployto/src"
	"deployto/src/filesystem"
	"deployto/src/types"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	helmclient "github.com/poncheg/go-helm-client"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
	"helm.sh/helm/v3/pkg/repo"
)

func init() {
	RunScriptFuncImplementations["helm"] = Helm
}

func Helm(target *types.Target, repositoryFS *filesystem.Filesystem, workdir string, aliases []string, rootValues, input types.Values, dump *src.ContextDump) (output types.Values, err error) {
	var outputBuffer bytes.Buffer
	//set settings for helm
	options := &helmclient.KubeConfClientOptions{
		Options: &helmclient.Options{
			Namespace:        target.Spec.Namespace, // Change this to the namespace you wish to install the chart in.
			RepositoryCache:  filepath.Join(os.TempDir(), ".helmcache"),
			RepositoryConfig: filepath.Join(os.TempDir(), ".helmrepo"),
			Debug:            true,
			Linting:          true, // Change this to false if you don't want linting.
			DebugLog: func(format string, v ...interface{}) {
				log.Debug().Str("ctx", "helm").Msgf(format, v...)
			},
			Output: &outputBuffer, // Not mandatory, leave open for default os.Stdout
		},
		KubeContext: "",
		KubeConfig:  target.LoadKubeconfig(),
	}
	helmClient, err := helmclient.NewClientFromKubeConf(options)

	var chartName string
	if err != nil {
		log.Error().Err(err).Msg("Create helm client error")
		return nil, err
	}
	// get repository url
	repository := types.Get(input, "", "repository")
	if repository == "" || filesystem.Supported(repository) {
		chartName = filepath.Join(repositoryFS.LocalPath, workdir)
	} else {
		if repository[len(repository)-1] != '/' {
			repository += "/"
		}
		u, err := url.Parse(repository)
		if err != nil {
			log.Error().Err(err).Msg("Url parsing  error")
			return nil, err
		}
		// get name for repository from url path
		ua := strings.Split(u.Path, "/")
		chartRepo := repo.Entry{
			Name: ua[1],
			URL:  repository,
		}
		// Add a chart-repository to the client.
		if err := helmClient.AddOrUpdateChartRepo(chartRepo); err != nil {
			log.Error().Err(err).Str("path", "helm").Msg("Add a chart-repository to the client error")
			return nil, err
		}
		chartName = chartRepo.Name + "/" + types.Get(input, aliases[len(aliases)-1], "name")
	}

	valuesFile, err := yaml.Marshal(&input)
	if err != nil {
		log.Error().Err(err).Str("path", "helm").Msg("Pasing yaml error")
		return nil, err
	}
	version := types.Get(input, aliases[len(aliases)-1], "version")
	// put settings for chart and put values
	chartSpec := helmclient.ChartSpec{
		ReleaseName:     buildAlias(aliases),
		ChartName:       chartName,
		Version:         version,
		ValuesYaml:      string(valuesFile),
		CreateNamespace: true,
		Namespace:       target.Spec.Namespace,
		UpgradeCRDs:     true,
		//Wait:            true,
		Timeout: time.Duration(5 * float64(time.Minute)),
		//DryRun:          true,
	}
	dump.Push("helmChartSpec", chartSpec)

	//helmClient.De

	// Install a chart release.
	// Note that helmclient.Options.Namespace should ideally match the namespace in chartSpec.Namespace.
	release, err := helmClient.InstallOrUpgradeChart(context.TODO(), &chartSpec, nil)
	if err != nil {
		log.Error().Err(err).Msg("Install chart error")
		return nil, err
	}

	dump.Push("manifest", release.Manifest)

	if release.Info.Status.String() != "deployed" {
		log.Error().Err(err).Msg("Release chart not deployed")
		return nil, err
	}

	scriptOutput := make(types.Values)
	//	scriptOutput["manifest"] = release.Manifest
	scriptOutput["values"] = release.Config
	scriptOutput["name"] = release.Name
	scriptOutput["version"] = release.Version
	return scriptOutput, nil
}
