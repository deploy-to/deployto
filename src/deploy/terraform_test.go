package deploy

import (
	"deployto/src/types"
	"testing"
)

func TestTerraform(t *testing.T) {
	// check Hello World from Terraform
	path := "../../examples/terraform-hello-world/"
	output, err := Terraform(nil, nil, path, []string{"AAAAA"}, types.Values(nil), types.Values(nil), nil)
	if err != nil {
		t.Fatalf("Terraform error %v", err)
	}
	t.Logf("Terraform output is = %v", output)
	// terraform yandex check
	path2 := "../../examples/terraform-yandex"
	output, err = Terraform(nil, nil, path2, []string{"AAAAA"}, types.Values(nil), types.Values(nil), nil)
	if err != nil {
		_, err = TerraformDestroy(nil, nil, path2, []string{"AAsAAA"}, types.Values(nil), types.Values(nil))
		if err != nil {
			t.Fatalf("Terraform destroy error %v", err)
		}
		t.Fatalf("Terraform error %v", err)
	}
	t.Logf("Terraform output is = %v", output)
	output, err = TerraformDestroy(nil, nil, path2, []string{"AAAAA"}, types.Values(nil), types.Values(nil))
	if err != nil {
		t.Fatalf("Terraform destroy error %v", err)
	}
	t.Logf("Terraform destroy output is = %v", output)
}
