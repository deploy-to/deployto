//go:build K8SIntegration

// for call this test 1) setup k8s environment  2) run $go test --tags=K8SIntegration  ./... -run K8SIntegration
package adapters

import (
	"deployto/src/deploy"
	"deployto/src/types"
	"strings"
	"testing"
)

func TestK8SIntegrationHelm(t *testing.T) {
	script := types.Values{
		"repository":       "https://charts.bitnami.com/bitnami",
		"name":             "postgresql",
		"version":          "15.1.0",
		"fullnameOverride": "postgresql-standalone",
		"auth": types.Values{
			"database":         "servicea-db",
			"username":         "servicea-user",
			"password":         "HOrFk14CyX",
			"postgresPassword": "xxdsdsddsxxxx",
		},
	}
	helm := &helmAdapter{}
	deploy := deploy.NewDeploy(nil, "", []string{"rootAlias"})
	output, err := helm.Apply(deploy, types.DecodeScript("", script), getTarget(t).AsValues())
	if err != nil {
		t.Fatalf("Helm error %v", err)
	}
	if output["values"].(types.Values)["auth"].(types.Values)["database"] != script["auth"].(types.Values)["database"] {
		t.Errorf("database name: GetValues()[database] = %v, want %v", output["values"].(types.Values)["auth"].(types.Values)["database"], script["auth"].(types.Values)["database"])
	}
	if output["values"].(types.Values)["auth"].(types.Values)["username"] != script["auth"].(types.Values)["username"] {
		t.Errorf("username name: GetValues()[username] = %v, want %v", output["values"].(types.Values)["auth"].(types.Values)["username"], script["auth"].(types.Values)["username"])
	}
	if output["values"].(types.Values)["auth"].(types.Values)["password"] != script["auth"].(types.Values)["password"] {
		t.Errorf("password name: GetValues()[password] = %v, want %v", output["values"].(types.Values)["auth"].(types.Values)["password"], script["auth"].(types.Values)["password"])
	}
	if output["values"].(types.Values)["auth"].(types.Values)["postgresPassword"] != script["auth"].(types.Values)["postgresPassword"] {
		t.Errorf("postgresPassword name: GetValues()[postgresPassword] = %v, want %v", output["values"].(types.Values)["auth"].(types.Values)["postgresPassword"], script["auth"].(types.Values)["postgresPassword"])
	}

	//Another Helm Chart
	script = types.Values{
		"repository":       "https://charts.bitnami.com/bitnami",
		"name":             "postgresql-ha",
		"version":          "14.0.10",
		"fullnameOverride": "postgresql-ha",
		"maxConnections":   "1000",
		"global": types.Values{
			"postgresql": types.Values{
				"database":         "service-db",
				"username":         "service-user",
				"password":         "878787878",
				"postgresPassword": "qweqweqweqweq",
				"repmgrDatabase":   "repmgr",
				"repmgrUsername":   "repmgr",
				"repmgrPassword":   "repmgrpss",
			},
		},
		"postgresql": types.Values{
			"database":         "service-db",
			"username":         "service-user",
			"password":         "878787878",
			"postgresPassword": "qweqweqweqweq",
			"repmgrDatabase":   "repmgr",
			"repmgrUsername":   "repmgr",
			"repmgrPassword":   "repmgrpss",
		},
		"auth": types.Values{
			"database":         "service-db",
			"username":         "service-user",
			"password":         "878787878",
			"postgresPassword": "qweqweqweqweq",
			"repmgrDatabase":   "repmgr",
			"repmgrUsername":   "repmgr",
			"repmgrPassword":   "repmgrpss",
		},
	}

	output, err = helm.Apply(deploy, types.DecodeScript("", script), types.Values(nil))

	if err != nil {
		t.Fatalf("Helm error %v", err)
	}
	if output["values"].(types.Values)["global"].(types.Values)["postgresql"].(types.Values)["database"] != script["global"].(types.Values)["postgresql"].(types.Values)["database"] {
		t.Errorf("database name: GetValues()[database] = %v, want %v", output["values"].(types.Values)["global"].(types.Values)["postgresql"].(types.Values)["database"], script["global"].(types.Values)["postgresql"].(types.Values)["database"])
	}
	if output["values"].(types.Values)["global"].(types.Values)["postgresql"].(types.Values)["username"] != script["global"].(types.Values)["postgresql"].(types.Values)["username"] {
		t.Errorf("username name: GetValues()[username] = %v, want %v", output["values"].(types.Values)["global"].(types.Values)["postgresql"].(types.Values)["username"], script["global"].(types.Values)["postgresql"].(types.Values)["username"])
	}
	if output["values"].(types.Values)["global"].(types.Values)["postgresql"].(types.Values)["password"] != script["global"].(types.Values)["postgresql"].(types.Values)["password"] {
		t.Errorf("password name: GetValues()[password] = %v, want %v", output["values"].(types.Values)["global"].(types.Values)["postgresql"].(types.Values)["password"], script["global"].(types.Values)["postgresql"].(types.Values)["password"])
	}
	if output["values"].(types.Values)["global"].(types.Values)["postgresql"].(types.Values)["postgresqlPassword"] != script["global"].(types.Values)["postgresql"].(types.Values)["postgresqlPassword"] {
		t.Errorf("postgresqlPassword name: GetValues()[postgresqlPassword] = %v, want %v", output["values"].(types.Values)["global"].(types.Values)["postgresql"].(types.Values)["postgresqlPassword"], script["global"].(types.Values)["postgresql"].(types.Values)["postgresqlPassword"])
	}

	//negative test Chart doesnt exist
	//Another Helm Chart
	script = types.Values{
		"repository":       "https://charts.bitnami.com/bitnami",
		"name":             "postgresql-ha",
		"version":          "99.99.99",
		"fullnameOverride": "postgresql-ha",
		"global": types.Values{
			"postgresql": types.Values{
				"database":           "service-db",
				"username":           "service-user",
				"password":           "878787878",
				"postgresqlPassword": "qweqweqweqweq",
			},
		},
	}
	_, err = helm.Apply(deploy, types.DecodeScript("", script), types.Values(nil))
	if err == nil {
		t.Fatalf("wait helm error")
	}
	errorText := "no chart version found for postgresql-ha-99.99.99"
	if !strings.Contains(err.Error(), errorText) {
		t.Errorf("Need chart error = %v, want %v", err.Error(), errorText)
	}
}

func getTarget(_ *testing.T) *types.Target {
	return &types.Target{
		Spec: types.TargetSpec{
			Kubeconfig: types.Kubeconfig{
				Namespace:  "test",
				UseDefault: true,
			},
		},
	}
}
