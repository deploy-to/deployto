package deploy

import (
	"deployto/src/filesystem"
	"deployto/src/types"
	"testing"
)

func TestTerraform(t *testing.T) {
	// check Hello World from Terraform
	// 	apiVersion: deployto.dev/v1beta1
	// kind: Target
	// metadata:
	//   name: local-target-ya-cloud
	// spec:
	//   terraform:
	//     provider: yandex
	//     path: "file://SERVER/target/local-target-ya-cloud/"
	//     env:
	//      # // TODO настроить получение секретов из hasicorp vault
	//      #  example     YC_TOKEN: ${{vault.host/storage-name/local-target-ya-cloud/token}}
	//       YC_ZONE: "ru-central1-d"
	//       YC_TOKEN: "t1.9euelZqei5eUm8fOiZiYiZXKnY-Kju3rnpWayZeNnpuXjpuWjJWNkpOZjInl8_dNXAVO-e8aHSdN_N3z9w0LA0757xodJ038zef1656Vmp7GyZGayImMl8eWj46SyZGL7_zN5_XrnpWalYmQzpCczJ2PiZbJkZqWiorv_cXrnpWansbJkZrIiYyXx5aPjpLJkYs.vnTrvH1R705ZQWRpW0becelwv1eflb2OYlv6_qhgIzDXTi2VdITjK6KUOfuSmB3ffJ0sHK0IqOa6EKa86rhVAg"
	//       YC_CLOUD_ID: "b1g8o6pkacilvmta12t9"
	//       YC_FOLDER_ID: "b1gdqo93r75d249hlicn"
	target := types.Target{
		Base: types.Base{
			APIVersion: "deployto.dev/v1beta1",
			Kind:       "Target",
			Meta: &types.MetaData{
				Name: "local-target-ya-cloud",
			},
		},
		Spec: &types.TargetSpec{
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
	path := "../../examples/terraform-hello-world"
	fs := filesystem.GetDeploytoRootFilesystem(filesystem.Get("file://"+path), "/")
	input := types.Values{
		"hello": "Hello, from deploy to!",
	}
	output, err := Terraform(&target, fs, path, []string{"AAAAA"}, types.Values(nil), input, nil)
	if err != nil {
		t.Fatalf("Terraform error %v", err)
	}
	t.Logf("Terraform output is = %v", output)
	// terraform yandex check
	path2 := "../../examples/terraform-yandex"
	fs = filesystem.GetDeploytoRootFilesystem(filesystem.Get("file://"+path), "/")
	output, err = Terraform(&target, fs, path2, []string{"AAAAA"}, types.Values(nil), input, nil)
	if err != nil {
		_, err = TerraformDestroy(&target, fs, path2, []string{"AAAAA"}, input, types.Values(nil))
		if err != nil {
			t.Fatalf("Terraform destroy error %v", err)
		}
		t.Fatalf("Terraform error %v", err)
	}
	t.Logf("Terraform output is = %v", output)
	output, err = TerraformDestroy(&target, fs, path2, []string{"AAAAA"}, input, types.Values(nil))
	if err != nil {
		t.Fatalf("Terraform destroy error %v", err)
	}
	t.Logf("Terraform destroy output is = %v", output)
}
