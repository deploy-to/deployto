package deploy

import (
	"bufio"
	"context"
	"deployto/src"
	"deployto/src/filesystem"
	"deployto/src/types"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/hashicorp/terraform-exec/tfexec"
)

func init() {
	RunScriptFuncImplementations["terraform"] = Terraform
}

func Terraform(target *types.Target, repositoryFS *filesystem.Filesystem, workdir string, aliases []string, rootValues, input types.Values, dump *src.ContextDump) (output types.Values, err error) {
	fullpath := filepath.Join(repositoryFS.LocalPath, workdir)

	//TODO find terraform exec
	execPath := "/usr/local/bin/terraform"

	//create tf client
	tf, err := tfexec.NewTerraform(fullpath, execPath)
	if err != nil {
		log.Error().Err(err).Msg("error running NewTerraformr")
		return nil, err
	}

	//Set map for input variable
	envs := make(map[string]string)
	if target != nil && target.Spec != nil {
		_, ok := target.Spec.Terraform["env"]
		if ok {
			targets := convertMapToString(target.Spec.Terraform["env"].(types.Values))
			for k, v := range targets {
				envs[k] = v
			}
		}
	}
	if input != nil {
		inputs := convertMapToString(input)
		for k, v := range inputs {
			envs[k] = v
		}
	}

	jsonData, err := json.Marshal(envs)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Convert JSON bytes to string
	jsonString := string(jsonData)

	//set output for tf client
	f := os.Stdout
	writer := bufio.NewWriter(f)
	tf.SetStdout(writer)
	tf.SetStderr(writer)

	// init
	err = tf.Init(context.Background(), tfexec.Upgrade(true))
	if err != nil {
		log.Error().Err(err).Msg("error running Init")
		return nil, err
	}

	// plan
	plan, err := tf.Plan(context.Background(), tfexec.Var("deployto_context="+jsonString))
	if err != nil {
		log.Error().Err(err).Msg("error running Plan")
		return nil, err
	}
	if plan {
		err := tf.Apply(context.Background(), tfexec.Var("deployto_context="+jsonString))
		if err != nil {
			log.Error().Err(err).Msg("error running Apply")
			return nil, err
		}
	} else {
		log.Error().Msg("please check Plan output")
		return nil, err
	}

	// read state file
	state, err := tf.Show(context.TODO())
	if err != nil {
		log.Error().Err(err).Msg("error running show state")
		return nil, err
	}

	// set output
	scriptOutput := make(types.Values)
	// format https://developer.hashicorp.com/terraform/internals/json-format#state-representation
	scriptOutput["state"] = state.Values

	return scriptOutput, nil
}

func TerraformDestroy(target *types.Target, repositoryFS *filesystem.Filesystem, workdir string, aliases []string, rootValues, input types.Values) (output types.Values, err error) {

	fullpath := filepath.Join(repositoryFS.LocalPath, workdir)
	//TODO find terraform exec
	execPath := "/usr/local/bin/terraform"

	tf, err := tfexec.NewTerraform(fullpath, execPath)
	if err != nil {
		log.Error().Err(err).Msg("error running NewTerraformr")
		return nil, err
	}
	f := os.Stdout
	writer := bufio.NewWriter(f)
	tf.SetStdout(writer)
	tf.SetStderr(writer)

	err = tf.Destroy(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("error running destory")
		return nil, err
	}
	return nil, nil
}

func fillmap(s map[string]any, d map[string]string, prefix string) {
	for key, value := range s {
		strKey := fmt.Sprintf("%v%v", prefix, key)
		strValue := fmt.Sprintf("%v", value)

		d[strKey] = strValue
	}
}

func convertMapToString(m map[string]interface{}) map[string]string {
	result := make(map[string]string)
	for key, value := range m {
		switch v := value.(type) {
		case int:
			result[strings.ToLower(key)] = strconv.Itoa(v)
		case string:
			result[strings.ToLower(key)] = v
		case bool:
			result[strings.ToLower(key)] = strconv.FormatBool(v)
		case map[string]any:
			// If the value is another map, recursively convert it to map[string]string
			innerMap := convertMapToString(v)
			for k, v := range innerMap {
				result[strings.ToLower(key)+"_"+strings.ToLower(k)] = v // Add a prefix to keys of nested maps
			}
		// Add more cases for other types if needed
		default:
			// Handle other types or skip them
			log.Printf("Skipping key '%s' with unsupported type %T\n", key, value)
		}
	}
	return result
}
