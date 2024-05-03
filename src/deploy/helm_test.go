// for call this test 1) setup k8s environment  2) run $go test --tags=K8SIntegration  ./... -run K8SIntegration
package deploy

import (
	"deployto/src/types"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"testing"
)

func TestK8SIntegrationHelm(t *testing.T) {
	inputs := types.Values{
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
	output, err := Helm(getTarget(t), "", []string{"AAAAA"}, types.Values(nil), inputs)
	if err != nil {
		t.Fatalf("Helm error %v", err)
	}
	if output["values"].(types.Values)["auth"].(types.Values)["database"] != inputs["auth"].(types.Values)["database"] {
		t.Errorf("database name: GetValues()[database] = %v, want %v", output["values"].(types.Values)["auth"].(types.Values)["database"], inputs["auth"].(types.Values)["database"])
	}
	if output["values"].(types.Values)["auth"].(types.Values)["username"] != inputs["auth"].(types.Values)["username"] {
		t.Errorf("username name: GetValues()[username] = %v, want %v", output["values"].(types.Values)["auth"].(types.Values)["username"], inputs["auth"].(types.Values)["username"])
	}
	if output["values"].(types.Values)["auth"].(types.Values)["password"] != inputs["auth"].(types.Values)["password"] {
		t.Errorf("password name: GetValues()[password] = %v, want %v", output["values"].(types.Values)["auth"].(types.Values)["password"], inputs["auth"].(types.Values)["password"])
	}
	if output["values"].(types.Values)["auth"].(types.Values)["postgresPassword"] != inputs["auth"].(types.Values)["postgresPassword"] {
		t.Errorf("postgresPassword name: GetValues()[postgresPassword] = %v, want %v", output["values"].(types.Values)["auth"].(types.Values)["postgresPassword"], inputs["auth"].(types.Values)["postgresPassword"])
	}

	//Another Helm Chart
	inputs = types.Values{
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
	output, err = Helm(getTarget(t), "", []string{"AAAAA"}, types.Values(nil), inputs)
	if err != nil {
		t.Fatalf("Helm error %v", err)
	}
	if output["values"].(types.Values)["global"].(types.Values)["postgresql"].(types.Values)["database"] != inputs["global"].(types.Values)["postgresql"].(types.Values)["database"] {
		t.Errorf("database name: GetValues()[database] = %v, want %v", output["values"].(types.Values)["global"].(types.Values)["postgresql"].(types.Values)["database"], inputs["global"].(types.Values)["postgresql"].(types.Values)["database"])
	}
	if output["values"].(types.Values)["global"].(types.Values)["postgresql"].(types.Values)["username"] != inputs["global"].(types.Values)["postgresql"].(types.Values)["username"] {
		t.Errorf("username name: GetValues()[username] = %v, want %v", output["values"].(types.Values)["global"].(types.Values)["postgresql"].(types.Values)["username"], inputs["global"].(types.Values)["postgresql"].(types.Values)["username"])
	}
	if output["values"].(types.Values)["global"].(types.Values)["postgresql"].(types.Values)["password"] != inputs["global"].(types.Values)["postgresql"].(types.Values)["password"] {
		t.Errorf("password name: GetValues()[password] = %v, want %v", output["values"].(types.Values)["global"].(types.Values)["postgresql"].(types.Values)["password"], inputs["global"].(types.Values)["postgresql"].(types.Values)["password"])
	}
	if output["values"].(types.Values)["global"].(types.Values)["postgresql"].(types.Values)["postgresqlPassword"] != inputs["global"].(types.Values)["postgresql"].(types.Values)["postgresqlPassword"] {
		t.Errorf("postgresqlPassword name: GetValues()[postgresqlPassword] = %v, want %v", output["values"].(types.Values)["global"].(types.Values)["postgresql"].(types.Values)["postgresqlPassword"], inputs["global"].(types.Values)["postgresql"].(types.Values)["postgresqlPassword"])
	}

	//negative test Chart doesnt exist
	//Another Helm Chart
	inputs = types.Values{
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
	output, err = Helm(getTarget(t), "", []string{"AAAAA"}, types.Values(nil), inputs)
	if err != nil {
		t.Fatalf("Helm error %v", err)
	}
	if !strings.HasPrefix(err.Error(), "no chart version found for postgresql-ha-99.99.99") {
		t.Errorf("Need chart error = %v, want %v", err.Error(), "no chart version found for postgresql-ha-99.99.99")
	}
}

func getTarget(t *testing.T) *types.Target {
	usr, err := user.Current()
	if err != nil {
		t.Fatalf("i'm a groot")
	}
	kubeConfigFilename := filepath.Join(usr.HomeDir, ".kube/config")
	KubeConfigEnv := os.Getenv("KUBECONFIG")
	if KubeConfigEnv != "" {
		kubeConfigFilename = KubeConfigEnv
	}

	kubeconfig, err := os.ReadFile(kubeConfigFilename)
	if err != nil {
		t.Fatalf("for call this test 1) setup k8s environment  2) run $go test --tags=K8SIntegration  ./... -run K8SIntegration")
	}

	return &types.Target{
		Namespace:  "test",
		Kubeconfig: kubeconfig,
	}
}
