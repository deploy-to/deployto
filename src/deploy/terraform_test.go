package deploy

import (
	"deployto/src/types"
	"testing"
)

func TestTerraform(t *testing.T) {
	// check Hello World from Terraform
	path := "/home/dagonch9/deployto/examples/terraform-hello-world/"
	output, err := Terraform(nil, nil, path, []string{"AAAAA"}, types.Values(nil), types.Values(nil))
	if err != nil {
		t.Fatalf("Terraform error %v", err)
	}
	t.Logf("Terraform output is = %v", output)
}
