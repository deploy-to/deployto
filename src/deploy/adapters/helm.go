package adapters

import (
	"bytes"
	"context"
	"deployto/src/deploy"
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
	deploy.DefaultAdapters["helm"] = &helmAdapter{}
}

type helmAdapter struct{}

func (h *helmAdapter) Apply(d *deploy.Deploy, script *types.Script, contex types.Values) (output types.Values, err error) {
	//TODO добавить логику использвания workdir/repositoryFS/helm repository
	target := types.DecodeTarget(contex["target"])

	var outputBuffer bytes.Buffer
	//set settings for helm
	options := &helmclient.KubeConfClientOptions{
		Options: &helmclient.Options{
			Namespace:        target.Spec.Kubeconfig.Namespace,
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
	repository := types.Get(contex, "", "repository")
	if repository == "" || filesystem.Supported(repository) {
		chartName = filepath.Join(d.FS.LocalPath, d.Workdir)
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
		chartName = chartRepo.Name + "/" + types.Get(contex, d.Aliases[len(d.Aliases)-1], "name")
	}

	valuesFile, err := yaml.Marshal(&contex)
	if err != nil {
		log.Error().Err(err).Str("path", "helm").Msg("Pasing yaml error")
		return nil, err
	}
	version := types.Get(contex, d.Aliases[len(d.Aliases)-1], "version")
	// put settings for chart and put values
	chartSpec := helmclient.ChartSpec{
		ReleaseName:     deploy.BuildAlias(d.Aliases),
		ChartName:       chartName,
		Version:         version,
		ValuesYaml:      string(valuesFile),
		CreateNamespace: true,
		Namespace:       target.Spec.Kubeconfig.Namespace,
		UpgradeCRDs:     true,
		//Wait:            true,
		Timeout: time.Duration(5 * float64(time.Minute)),
		//DryRun:          true,
	}
	d.Keeper.Push("helmChartSpec", chartSpec)

	release, err := helmClient.InstallOrUpgradeChart(context.TODO(), &chartSpec, nil)
	if err != nil {
		log.Error().Err(err).Msg("Install chart error")
		return nil, err
	}

	d.Keeper.Push("manifest", release.Manifest)

	if release.Info.Status.String() != "deployed" {
		log.Error().Err(err).Msg("Release chart not deployed")
		return nil, err
	}

	scriptOutput := make(types.Values)
	//	scriptOutput["manifest"] = release.Manifest
	scriptOutput["Values"] = release.Config
	scriptOutput["helmReleaseName"] = release.Name
	scriptOutput["helmReleaseVersion"] = release.Version
	return scriptOutput, nil
}

func (h *helmAdapter) Destroy(d *deploy.Deploy, script *types.Script, contex types.Values) error {
	panic("NOT IMPLIMENTED")
}
