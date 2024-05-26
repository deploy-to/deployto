package adapters

import (
	"bufio"
	"context"
	"deployto/src/deploy"
	"deployto/src/filesystem"
	"deployto/src/types"
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"

	"github.com/hashicorp/terraform-exec/tfexec"
	tfjson "github.com/hashicorp/terraform-json"
)

func init() {
	deploy.DefaultAdapters["terraform"] = &terraformAdapter{}
}

type terraformAdapter struct{}

var TerraformBinary string //setup in cli.StringFlag{ Name: "terraform-binary", Destination: &deploy.TerraformBinary,

func (t *terraformAdapter) Apply(d *deploy.Deploy, script *types.Script, compContext types.Values) (output types.Values, err error) {
	fullpath := filepath.Join(d.FS.LocalPath, d.Workdir)

	//TODO find terraform exec OR use docker instalation OR import github.com/hashicorp/terraform/command ("Business Source License")
	if _, err := os.Stat(TerraformBinary); err != nil {
		if _, err := os.Stat("/usr/local/bin/terraform"); err == nil {
			TerraformBinary = "/usr/local/bin/terraform"
		} else {
			log.Error().Err(err).Str("path", TerraformBinary).Msg("Terraform not found. Set terraform-binary config")
			return nil, err
		}
	}

	//create tf client
	tf, err := tfexec.NewTerraform(fullpath, TerraformBinary)
	if err != nil {
		log.Error().Err(err).Msg("Terraform - error running NewTerraform")
		return nil, err
	}

	//set output for tf client
	f := os.Stdout
	writer := bufio.NewWriter(f)
	tf.SetStdout(writer)
	tf.SetStderr(writer)

	if _, err := os.Stat(filepath.Join(os.Getenv("HOME"), ".terraformrc")); err != nil {
		log.Warn().Msg("Yandex provider need config https://yandex.cloud/ru/docs/tutorials/infrastructure-management/terraform-quickstart#configure-provider")
	}

	jsonData, err := json.Marshal(compContext)
	if err != nil {
		log.Error().Err(err).Msg("terraform - json.Marshal error")
		return
	}

	// init
	err = tf.Init(context.TODO(), tfexec.Upgrade(true))
	if err != nil {
		log.Error().Err(err).Msg("error running Init")
		return nil, err
	}

	// plan
	plan, err := tf.Plan(context.TODO(), tfexec.Var("deployto_context="+string(jsonData)))
	if err != nil {
		log.Error().Err(err).Msg("error running Plan")
		return nil, err
	}
	if !plan {
		log.Debug().Msg("The plan does not exist, but I call the apply to get the output")
	}

	err = tf.Apply(context.TODO(), tfexec.Var("deployto_context="+string(jsonData)))
	if err != nil {
		log.Error().Err(err).Msg("error running Apply")
		return nil, err
	}

	// read state file
	state, err := tf.Show(context.TODO())
	if err != nil {
		log.Error().Err(err).Msg("error running show state")
		return nil, err
	}

	return tfJsonStateOutput2Values(state.Values.Outputs, &d.Secrets), nil
}

func tfJsonStateOutput2Values(tf map[string]*tfjson.StateOutput, secrets *types.Secrets) types.Values {
	result := make(types.Values, len(tf))
	if deploytoOutputRaw, deploytoOutputExists := tf["deployto_output"]; deploytoOutputExists {
		deploytoOutput, ok := deploytoOutputRaw.Value.(types.Values)
		if !ok {
			log.Error().Type("type", deploytoOutputRaw.Value).Msg("tfJsonStateOutput2Values - not implemented type")
			return nil
		}
		for k, valueRaw := range deploytoOutput {
			switch v := valueRaw.(type) {
			case string:
				result[k] = v
				if deploytoOutputRaw.Sensitive {
					secrets.Add(v)
				}
			//TODO для массива и мапы надо пробижаьться по всем элементам, и если значение строка, то проследить за Sensitive
			default:
				result[k] = valueRaw
			}
		}
	} else {
		for k, v := range tf {
			switch stateOutputValue := v.Value.(type) {
			case string:
				result[k] = stateOutputValue
				if v.Sensitive {
					secrets.Add(stateOutputValue)
				}
			default:
				result[k] = v.Value
			}
		}
	}
	return result
}

func (t *terraformAdapter) Destroy(d *deploy.Deploy, script *types.Script, compContext types.Values) error {
	panic("NOT IMPLIMENTED")
}

func TerraformDestroy(target *types.Target, repositoryFS *filesystem.Filesystem, workdir string, aliases []string, appContext, input types.Values) (output types.Values, err error) {
	fullpath := filepath.Join(repositoryFS.LocalPath, workdir)
	//TODO find terraform exec
	execPath := "/usr/local/bin/terraform"

	tf, err := tfexec.NewTerraform(fullpath, execPath)
	if err != nil {
		log.Error().Err(err).Msg("TerraformDestroy - error running NewTerraform")
		return nil, err
	}
	f := os.Stdout
	writer := bufio.NewWriter(f)
	tf.SetStdout(writer)
	tf.SetStderr(writer)

	err = tf.Destroy(context.TODO())
	if err != nil {
		log.Error().Err(err).Msg("error running destory")
		return nil, err
	}
	return nil, nil
}
