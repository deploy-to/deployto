package deploy

import (
	"context"
	"deployto/src/filesystem"
	"deployto/src/types"

	"github.com/rs/zerolog/log"

	"github.com/hashicorp/terraform-exec/tfexec"
)

func init() {
	RunScriptFuncImplementations["terraform"] = Terraform
}

func Terraform(target *types.Target, repositoryFS *filesystem.Filesystem, workdir string, aliases []string, rootValues, input types.Values) (output types.Values, err error) {

	execPath := "/usr/local/bin/terraform"

	tf, err := tfexec.NewTerraform(workdir, execPath)
	if err != nil {
		log.Error().Err(err).Msg("error running NewTerraformr")
		return nil, err
	}

	plan, err := tf.Plan(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("error running Plan")
		return nil, err
	}
	if plan {
		err := tf.Apply(context.Background())
		if err != nil {
			log.Error().Err(err).Msg("error running Apply")
			return nil, err
		}
	}
	//state, err := tf.ShowStateFile(context.Background())
	output["terraform"] = ""

	return
}
