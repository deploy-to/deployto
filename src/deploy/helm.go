package deploy

import (
	"bytes"
	"context"
	"deployto/src/types"
	localyaml "deployto/src/yaml"
	"net/url"
	"strings"

	helmclient "github.com/poncheg/go-helm-client"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
	"helm.sh/helm/v3/pkg/repo"
)

func init() {
	RunScriptFuncImplementations["helm"] = Helm
}

func Helm(target *types.Target, workdir string, aliases []string, rootValues, input types.Values) (output types.Values, err error) {
	var outputBuffer bytes.Buffer
	// log.Error().Err(err).Str("path", path).Msg("Application/Components search error")
	// log.Debug().Str("environment", environmentArg).Msg("start deployto")
	//set settings for helm
	opt := &helmclient.KubeConfClientOptions{
		Options: &helmclient.Options{
			Namespace:        target.Namespace, // Change this to the namespace you wish to install the chart in.
			RepositoryCache:  "/tmp/.helmcache",
			RepositoryConfig: "/tmp/.helmrepo",
			Debug:            true,
			Linting:          true, // Change this to false if you don't want linting.
			DebugLog:         func(format string, v ...interface{}) {},
			Output:           &outputBuffer, // Not mandatory, leave open for default os.Stdout
		},
		KubeContext: "",
		KubeConfig:  target.Kubeconfig,
	}

	helmClient, err := helmclient.NewClientFromKubeConf(opt)

	if err != nil {
		log.Error().Err(err).Msg("Create helm client error")
		return nil, err
	}
	// get repository url
	repository := types.Get(input, "", "repository")
	u, err := url.Parse(repository)
	if err != nil {
		log.Error().Err(err).Msg("Url parsing  error")
		return nil, err
	}
	// get name for repository from url path
	ua := strings.Split(u.Path, "/")
	chartRepo := repo.Entry{
		Name: ua[0],
		URL:  repository,
	}

	// Add a chart-repository to the client.
	if err := helmClient.AddOrUpdateChartRepo(chartRepo); err != nil {
		log.Error().Err(err).Str("path", "helm").Msg("Add a chart-repository to the client error")
	}
	valuesFile, err := yaml.Marshal(&input)
	if err != nil {
		log.Error().Err(err).Str("path", "helm").Msg("Pasing yaml error")
	}
	kind := types.Get(input, aliases[len(aliases)-1], "repository")
	// put settings for chart and put values
	chartSpec := helmclient.ChartSpec{
		ReleaseName: kind,
		ChartName:   chartRepo.Name + kind,
		//нужна версия чарта которую деплоим
		//Version: "",
		ValuesYaml:  string(valuesFile),
		Namespace:   target.Namespace,
		UpgradeCRDs: true,
		Wait:        true,
	}

	// Install a chart release.
	// Note that helmclient.Options.Namespace should ideally match the namespace in chartSpec.Namespace.
	_, err = helmClient.InstallOrUpgradeChart(context.TODO(), &chartSpec, nil)
	if err != nil {
		log.Error().Err(err).Str("path", "helm").Msg("Install chart error")
	}

	poutput, err := helmClient.GetReleaseValues(kind, true)
	if err != nil {
		log.Error().Err(err).Str("path", "helm").Msg("Get Release chart error")
		return nil, err
	}
	template, err := helmClient.TemplateChart(&chartSpec, nil)
	if err != nil {
		log.Error().Err(err).Str("path", "helm").Msg("Template chart error")
	}
	services, ingresses := localyaml.GetBytes2(template)
	hostList := searchHost(services, ingresses, target.Namespace)

	poutput["host"] = hostList
	return poutput, nil
}

func searchHost(services []types.Service, ingresses []types.Ingress, namespace string) (host []map[string]interface{}) {
	svc := make(map[string]interface{})
	svcList := []string{}
	for _, s := range services {
		for _, p := range s.Spec.Ports {
			svcList = append(svcList, s.Metadata.Name+namespace+"svc.cluster.local"+":"+string(p.Port))
		}

	}
	svc["service"] = svcList

	ing := make(map[string]interface{})
	ingList := []string{}
	for _, s := range ingresses {
		for _, p := range s.Spec.Rules {
			ingList = append(ingList, p.Host)
		}

	}
	ing["ingress"] = ingList
	host = append(host, svc, ing)
	return
}
