package adapters

import (
	"deployto/src/deploy"
	"deployto/src/filesystem"
	"deployto/src/types"
	"os"
	"path/filepath"
	"testing"
)

func TestTerraform(t *testing.T) {
	target := getTestTarget()
	input := types.Values{
		"hello": "Hello, from deploy to!",
	}
	// terraform yandex check
	fs := filesystem.Get("file://../../examples")
	deploy := deploy.NewDeploy(fs, "terraform-yandex", []string{"AAAAA"})
	output, err := Terraform(deploy, nil, types.Values(nil))
	if err != nil {
		_, err = TerraformDestroy(target, fs, "terraform-yandex", []string{"AAAAA"}, input, types.Values(nil))
		if err != nil {
			t.Fatalf("Terraform destroy error %v", err)
		}
		t.Fatalf("Terraform error %v", err)
	}
	t.Logf("Terraform output is = %v", output)
	output, err = TerraformDestroy(target, fs, "terraform-yandex", []string{"AAAAA"}, input, types.Values(nil))
	if err != nil {
		t.Fatalf("Terraform destroy error %v", err)
	}
	t.Logf("Terraform destroy output is = %v", output)
}

func TestTerraformHelloWorld(t *testing.T) {
	fs := filesystem.Get("temp")
	defer fs.Destroy()
	_ = os.WriteFile(filepath.Join(fs.LocalPath, "main.tf"), []byte(`
variable "deployto_context" {
	type = map(string)
}

output "hello_deployto" {	
	value = var.deployto_context["hello"]
}
	`), 0666)
	context := types.Values{
		"target": getTestTarget().AsValues(),
	}
	script := &types.Script{
		Values: types.Values{"hello": "Hello, from deploy to!"},
	}
	deploy := deploy.NewDeploy(fs, "/", nil)
	output, err := Terraform(deploy, script, context)
	if err != nil {
		t.Fatalf("Terraform error %v", err)
	}
	if output["hello_deployto"] != "Hello, from deploy to!" {
		t.Fatalf("Terraform error %v", err)
	}
}

func getTestTarget() *types.Target {
	return &types.Target{
		Base: types.Base{
			APIVersion: "deployto.dev/v1beta1",
			Kind:       "Target",
			Meta: types.MetaData{
				Name: "local-target-ya-cloud",
			},
		},
		Spec: types.TargetSpec{
			Terraform: types.Values{
				"provider": "yandex",
				"path":     "file://SERVER/target/local-target-ya-cloud/",
				"env": types.Values{
					"YC_ZONE":      "ru-central1-d",
					"YC_TOKEN":     "t1.9euelZqei5eUm8fOiZiYiZXKnY-Kju3rnpWayZeNnpuXjpuWjJWNkpOZjInl8_dNXAVO-e8aHSdN_N3z9w0LA0757xodJ038zef1656Vmp7GyZGayImMl8eWj46SyZGL7_zN5_XrnpWalYmQzpCczJ2PiZbJkZqWiorv_cXrnpWansbJkZrIiYyXx5aPjpLJkYs.vnTrvH1R705ZQWRpW0becelwv1eflb2OYlv6_qhgIzDXTi2VdITjK6KUOfuSmB3ffJ0sHK0IqOa6EKa86rhVAg",
					"YC_CLOUD_ID":  "b1g8o6pkacilvmta12t9",
					"YC_FOLDER_ID": "b1gdqo93r75d249hlicn",
				},
			},
		},
	}
}
