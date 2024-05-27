package adapters

import (
	"bytes"
	"context"
	"deployto/src/deploy"
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
	deploy.DefaultAdapters["job"] = (*job)(nil)
}

type job struct{}

func (j *job) Apply(d *deploy.DeployState, script *types.Script, scriptContext types.Values) (output types.Values, err error) {
	//TODO определить, как искать path у job. Т.к. с одной стороны он должен указывать на место, где искать описание job, а с другой, на место, где будет выполняться
	jobs := yaml.Get[types.Job](d.FS, d.FS.FS.Join(d.Workdir, ".deployto"))
	for _, job := range jobs {
		if job.Meta.Name == script.Name {
			return runJob(d, job, scriptContext)
		}
	}
	log.Error().Str("jobName", script.Name).Msg("job not found")
	return nil, errors.New("job not found")
}

func runJob(d *deploy.DeployState, job *types.Job, scriptContext types.Values) (types.Values, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10)*time.Minute) //TODO timeout to job file? or to cmg.Flags?
	defer cancel()

	tmpPath := filepath.Join(os.TempDir(), "deployto-wf", deploy.BuildAlias(d.Aliases), shortuuid.New())
	err := os.MkdirAll(tmpPath, 0766)
	if err != nil {
		log.Error().Err(err).Str("tmpPath", tmpPath).Strs("aliases", d.Aliases).Msg("Job create tmp dir - error")
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
			// TODO  унифицировать с подготовкой контекста для скриптов
			envTemplatedValue, err := d.Templater.TemplatingString(envValue, types.MergeValues(
				d.Root.Values,
				scriptContext,
				types.Values{
					"aliases": d.Aliases,
					"alias":   deploy.BuildAlias(d.Aliases),
				},
			))
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
				command.Dir = filepath.Join(d.FS.LocalPath, d.Workdir)
				command.Env = append(env, "DEPLOYTO_OUTPUT="+stepOutputFile)
				command.Stderr = stderr
				command.Stdout = stdout

				lineContext := d.Child(d.FS, d.Workdir, []string{""})
				lineContext.Keeper.Push("command", map[string]any{
					"path": command.Path,
					"args": command.Args,
					"dir":  command.Dir,
					"env":  command.Env,
				})
				err = command.Run()
				log.Debug().Str("stdout", stdout.String()).Msg("command.Run() - output")
				lineContext.Keeper.Push("stdout", stdout.String())
				if stderr.Len() > 0 {
					log.Warn().Str("stderr", stderr.String()).Msg("command.Run() - return error")
					lineContext.Keeper.Push("stderr", stderr.String())
				}
				if !step.ContinueOnError && err != nil {
					log.Error().Err(err).Str("stepLine", stepLine).Msg("command run error")
					return nil, err
				}

				if step.Id != "" { //read output only if step.Id defined
					outputRaw, err := os.ReadFile(stepOutputFile)
					if err != nil || len(outputRaw) == 0 {
						log.Debug().Err(err).Str("stepLine", stepLine).Str("stepOutputFile", stepOutputFile).Msg("no output file or empty")
					} else {
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

func (j *job) Destroy(d *deploy.DeployState, script *types.Script, scriptContext types.Values) error {
	panic("NOT IMPLIMENTED")
}
