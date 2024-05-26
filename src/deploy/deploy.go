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

type Deploy struct {
	Root      *Deploy
	FS        *filesystem.Filesystem
	Workdir   string
	Aliases   []string
	Keeper    *DeployKeeper
	Templater *Templater
	Secrets   types.Secrets
	Values    types.Values
}

func NewDeploy(fs *filesystem.Filesystem, workdir string, aliases []string) *Deploy {
	if fs.Type != filesystem.LOCAL {
		return nil
	}
	localPath := filepath.Join(fs.LocalPath, workdir, filesystem.DeploytoDirName, "deployed")

	d := &Deploy{
		FS:      fs,
		Workdir: workdir,
		Aliases: aliases,
		Keeper:  GetDeployKeeper(localPath),
		Secrets: types.NewSecrets(),
	}
	d.Root = d
	d.Templater = NewTemplater(d)
	return d
}

func (d *Deploy) Child(fs *filesystem.Filesystem, workdir string, aliases []string) *Deploy {
	newDeploy := &Deploy{
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

func (d *Deploy) Apply(envName string) (err error) {
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

func (d *Deploy) Destroy(env string) error {
	return nil
}
