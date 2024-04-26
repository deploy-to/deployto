package deploy

import (
	"bytes"
	"deployto/src/types"
	"net/url"
	"strings"
	"github.com/poncheg/go-helm-client"
	"gopkg.in/yaml.v3"
)

func init() {
	RunScripts["helm"] = HelmRunScript
}

func HelmRunScript(kubeconfig string, names []string, kind string, script *types.Script, target *types.Target, input map[string]any) (output map[string]any, err error) {
	// эта функци будет вызыватсья только для script.type = helm
	// для script.type == helm, атрибут kind можно игнорировать
	var outputBuffer bytes.Buffer

	opt := &Options{
		Namespace:        target.Namespace, // Change this to the namespace you wish the client to operate in.
		RepositoryCache:  "/tmp/.helmcache",
		RepositoryConfig: "/tmp/.helmrepo",
		Debug:            true,
		Linting:          true,
		DebugLog:         func(format string, v ...interface{}) {},
		Output:           &outputBuffer, // Not mandatory, leave open for default os.Stdout
	}

	helmClient, err := helmclient.New(opt)
	if err != nil {
		panic(err)
	}
	u, err := url.Parse(script.Repository)
	if err != nil {
		log.Debug(err)
	}
	ua := strings.Split(u, "/")
	chartRepo := repo.Entry{
		Name: ua[0],
		URL:  script.Repository,
	}

	// Add a chart-repository to the client.
	if err := helmClient.AddOrUpdateChartRepo(chartRepo); err != nil {
		panic(err)
	}
	valuesFile, err := yaml.Marshal(&input)
	if err != nil {
		panic(err)
	}

	chartSpec := ChartSpec{
		ReleaseName: kind,
		ChartName:   chartRepo + kind,
		//нужна версия чарта которую деплоим
		//Version: "",
		ValuesYaml:  string(valuesFile),
		Namespace:   target.Namespace,
		UpgradeCRDs: true,
		Wait:        true,
	}

	// Install a chart release.
	// Note that helmclient.Options.Namespace should ideally match the namespace in chartSpec.Namespace.
	release, err := helmClient.InstallOrUpgradeChart(&chartSpec, nil); err != nil {
	if err != nil {
		panic(err)
	}

	poutput, err := helmClient.GetReleaseValues(kind) 
	if err != nil {
		panic(err)
	}
	template, err := helmClient.TemplateChart(&chartSpec, nil)
	if err != nil {
		panic(err)
	}
	is := yaml.GetBytes(template)
	return nil, nil
}
