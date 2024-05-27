package deploy

import (
	"deployto/src/filesystem"
	"deployto/src/types"
	"deployto/src/yaml"
	"errors"
	"path/filepath"
	"slices"

	"github.com/rs/zerolog/log"
)

const (
	NEW        string = "NEW"
	DEPLOYING  string = "DEPLOYING"
	DEPLOYED   string = "DEPLOYED"
	DESTROYING string = "DESTROYING"
	DESTROYED  string = "UNDEPLOYED"
	ERROR      string = "ERROR"
)

type DeployState struct {
	Root      *DeployState
	FS        *filesystem.Filesystem
	Workdir   string
	Aliases   []string
	Keeper    *DeployKeeper
	Templater *Templater
	Secrets   types.Secrets
	Values    types.Values
	State     string
}

func NewDeploy(fs *filesystem.Filesystem, workdir string, aliases []string) *DeployState {
	if fs.Type != filesystem.LOCAL {
		log.Error().Msg("NewDeploy not implemented for not LOCAL filesystem")
		return nil
	}
	localPath := filepath.Join(fs.LocalPath, workdir, filesystem.DeploytoDirName, "deployed")

	d := &DeployState{
		FS:      fs,
		Workdir: workdir,
		Aliases: aliases,
		Keeper:  GetDeployKeeper(localPath),
		Secrets: types.NewSecrets(),
		State:   NEW,
	}
	d.Root = d
	d.Templater = NewTemplater(d)
	return d
}

func (d *DeployState) Child(fs *filesystem.Filesystem, workdir string, aliases []string) *DeployState {
	newDeploy := &DeployState{
		Root:    d.Root,
		FS:      fs,
		Workdir: workdir,
		Aliases: aliases,
		Keeper:  d.Keeper.Next(BuildAlias(aliases)),
		Secrets: d.Secrets,
	}
	newDeploy.Templater = NewTemplater(d)
	return newDeploy
}

func (d *DeployState) Apply(envName string) (err error) {
	// Envirement
	environments := yaml.Get[types.Environment](d.FS, filesystem.DeploytoDirName)
	var environment *types.Environment
	for _, e := range environments {
		if e.Base.Meta.Name == envName {
			environment = e
		}
	}
	if environment == nil {
		log.Error().Int("len(environments)", len(environments)).Str("path", d.FS.URI).Str("waitEnvironment", envName).Msg("environment not found")
		return errors.New("ENVIRONMENT NOT FOUND")
	}
	log.Debug().Str("name", environment.Base.Meta.Name).Msg("Environment found")
	// Targets
	var targets []*types.Target
	for _, t := range yaml.Get[types.Target](d.FS, filesystem.DeploytoDirName) {
		if slices.Contains(environment.Spec.Targets, t.Base.Meta.Name) {
			targets = append(targets, t)
		}
	}
	if len(targets) != len(environment.Spec.Targets) {
		log.Error().Int("len(targets)", len(targets)).Int("len(environment.Spec.Targets)", len(environment.Spec.Targets)).Msg("Target not found")
		return errors.New("TARGET NOT FOUND")
	}
	log.Debug().Int("len(targets)", len(targets)).Msg("Targets found")

	log.Info().Str("file", environment.Base.Status.FileName).Str("name", environment.Base.Meta.Name).Msg("Deploy environment")
	for _, target := range targets {
		log.Info().Str("file", target.Base.Status.FileName).Str("name", target.Base.Meta.Name).Msg("Deploy target")

		context := types.Values{
			"environment": environment.Base.Meta.Name,
			"target":      target.AsValues(),
		}
		componentAdapter := DefaultAdapters["component"]
		release, e := componentAdapter.Apply(d, nil, context)
		if e != nil {
			log.Error().Err(e).Msg("Component deploy error")
			err = errors.Join(err, e)
		}
		//TODO Save release
		log.Debug().Any("release", release).Msg("Component deploy result")

		//TODO Run target script (move this logic to src/deploy/)
	}
	//TODO Run Env script (move this logic to src/deploy/?   Or stop move on target?)

	return err
}

func (d *DeployState) Destroy(env string) error {
	return nil
}
