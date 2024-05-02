//go:build K8SIntegration

// for call this test 1) setup k8s environment  2) run $go test --tags=K8SIntegration  ./... -run K8SIntegration
package deploy

import (
	"deployto/src/types"
	"os"
	"os/user"
	"path/filepath"
	"testing"
)

func TestK8SIntegrationHelm(t *testing.T) {
	output, err := Helm(getTarget(t), "", []string{"AAAAA"}, types.Values(nil), types.Values{
		"repository": "https://charts.bitnami.com/bitnami",
		"name":       "postgresql",
	})
	if err != nil {
		t.Fatalf("Helm error %v", err)
	}
	t.Logf("Output: %v", output)
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

func checkIfError(t *testing.T, err error) {
	if err != nil {
		t.FailNow()
	}
}
