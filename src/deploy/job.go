package deploy

import (
	"bytes"
	"context"
	"deployto/src/filesystem"
	"deployto/src/types"
	"deployto/src/yaml"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/lithammer/shortuuid/v3"
	"github.com/rs/zerolog/log"
)

func init() {
	RunScriptFuncImplementations["job"] = JobScript
}

func JobScript(target *types.Target, fs *filesystem.Filesystem, workdir string, aliases []string, rootContext, context types.Values) (output types.Values, err error) {
	resource := types.Get(context, "", "resource")
	if resource == "" {
		log.Error().Msg("job name not found")
		return nil, errors.New("job name not found")
	}

	jobs := yaml.Get[types.Job](fs, fs.FS.Join(workdir, ""))
	for _, job := range jobs {
		if job.Meta.Name == resource {
			return runJob(fs, workdir, job, aliases, context)
		}
	}
	log.Error().Str("jobName", resource).Msg("job not found")
	return nil, errors.New("job not found")
}

func runJob(fs *filesystem.Filesystem, workdir string, job *types.Job, aliases []string, jobContext types.Values) (types.Values, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10)*time.Minute) //TODO timeout to job file? or to cmg.Flags?
	defer cancel()

	tmpPath := filepath.Join(os.TempDir(), "deployto-wf", buildAlias(aliases), shortuuid.New())
	err := os.MkdirAll(tmpPath, 0766)
	if err != nil {
		log.Error().Err(err).Str("tmpPath", tmpPath).Strs("aliases", aliases).Msg("Job create tmp dir - error")
		return nil, err
	}
	defer func() {
		err := os.RemoveAll(tmpPath)
		if err != nil {
			log.Error().Str("path", tmpPath).Msg("remove tmp dir error")
		}
	}()

	output := make(types.Values)

	for _, step := range job.Spec.Steps {
		var env []string
		for envKey, envValue := range step.Env {
			envTemplatedValue, err := templatingString(envValue, mergeContext(nil, jobContext, aliases))
			if err != nil {
				log.Error().Err(err).Str("key", envKey).Str("template", envValue).Msg("Templating error")
				return nil, err
			}
			env = append(env, fmt.Sprintf("%s=%s", envKey, envTemplatedValue))
		}

		var stepOutputFile string
		if step.Id != "" {
			stepOutputFile = filepath.Join(tmpPath, step.Id)
		} else {
			stepOutputFile = filepath.Join(tmpPath, shortuuid.New())
		}

		for _, stepLine := range strings.Split(step.Run, "\n") {
			if len(strings.Trim(stepLine, " \t")) > 0 {
				stdout := new(bytes.Buffer)
				stderr := new(bytes.Buffer)

				command := exec.CommandContext(ctx, "bash", "-c", stepLine)
				command.Dir = filepath.Join(fs.LocalPath, workdir)
				command.Env = append(env, "DEPLOYTO_OUTPUT="+stepOutputFile)
				command.Stderr = stderr
				command.Stdout = stdout
				err = command.Run()
				if step.ContinueOnError || err != nil {
					log.Error().Err(err).Str("stepLine", stepLine).Msg("command run error")
					return nil, err
				}
				log.Debug().Str("stdout", stdout.String()).Msg("command.Run() - output")
				if stderr.Len() > 0 {
					log.Warn().Str("stderr", stderr.String()).Msg("command.Run() - return error")
				}

				if step.Id != "" { //read output only if step.Id defined
					outputRaw, err := os.ReadFile(stepOutputFile)
					if err != nil {
						log.Error().Err(err).Str("stepLine", stepLine).Str("stepOutputFile", stepOutputFile).Msg("read output error")
						return nil, err
					}
					if len(outputRaw) > 0 {
						stepOutput := make(map[string]string)
						for _, outputLine := range strings.Split(string(outputRaw), "\n") {
							if len(outputLine) > 0 {
								kvvv := strings.Split(outputLine, "=")
								stepOutput[kvvv[0]] = strings.Join(kvvv[1:], "=")
							}
						}
						output[step.Id] = stepOutput
					}
				}
			}
		}
	}
	return output, nil
}
