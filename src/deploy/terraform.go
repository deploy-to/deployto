package deploy

import (
	"bufio"
	"context"
	"deployto/src/filesystem"
	"deployto/src/types"
	"os"

	"github.com/rs/zerolog/log"

	"github.com/hashicorp/terraform-exec/tfexec"
)

func init() {
	RunScriptFuncImplementations["terraform"] = Terraform
}

func Terraform(target *types.Target, repositoryFS *filesystem.Filesystem, workdir string, aliases []string, rootValues, input types.Values) (output types.Values, err error) {

	//TODO find terraform exec
	execPath := "/usr/local/bin/terraform"

	tf, err := tfexec.NewTerraform(workdir, execPath)
	if err != nil {
		log.Error().Err(err).Msg("error running NewTerraformr")
		return nil, err
	}
	f := os.Stdout
	writer := bufio.NewWriter(f)
	tf.SetStdout(writer)
	tf.SetStderr(writer)
	err = tf.Init(context.Background(), tfexec.Upgrade(true))
	if err != nil {
		log.Error().Err(err).Msg("error running Init")
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
	// read state file
	state, err := tf.Show(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("error running show state")
		return nil, err
	}

	scriptOutput := make(types.Values)
	// format https://developer.hashicorp.com/terraform/internals/json-format#state-representation
	scriptOutput["state"] = state.Values

	return scriptOutput, nil
}

func TerraformTest(target *types.Target, repositoryFS *filesystem.Filesystem, workdir string, aliases []string, rootValues, input types.Values) (output types.Values, err error) {

	//TODO find terraform exec
	execPath := "/usr/local/bin/terraform"

	tf, err := tfexec.NewTerraform(workdir, execPath)
	if err != nil {
		log.Error().Err(err).Msg("error running NewTerraformr")
		return nil, err
	}
	f := os.Stdout
	writer := bufio.NewWriter(f)
	tf.SetStdout(writer)
	tf.SetStderr(writer)
	err = tf.Init(context.Background(), tfexec.Upgrade(true))
	if err != nil {
		log.Error().Err(err).Msg("error running Init")
		return nil, err
	}

	err = tf.Test(context.Background(), writer)
	if err != nil {
		log.Error().Err(err).Msg("error running Plan")
		return nil, err
	}
	// read state file
	state, err := tf.Show(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("error running show state")
		return nil, err
	}
	scriptOutput := make(types.Values)
	// format https://developer.hashicorp.com/terraform/internals/json-format#state-representation
	scriptOutput["state"] = state.Values

	return scriptOutput, nil
}
