package cmd

import (
	"deployto/src/deploy"
	"deployto/src/filesystem"
	"errors"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

func Deployto(cCtx *cli.Context) error {
	environmentArg := cCtx.Args().First()
	if len(environmentArg) == 0 {
		environmentArg = "local"
	}
	log.Debug().Str("environment", environmentArg).Msg("start deployto")

	path, err := os.Getwd()
	if err != nil {
		log.Error().Err(err).Msg("Get workdir error")
		return err
	}

	fs := filesystem.GetDeploytoRootFilesystem(filesystem.Get("file://"+path), "/")
	if fs == nil {
		log.Error().Msg("components dir (.deployto) not found")
		return errors.New("components dir (.deployto) not found")
	}

	deploy := deploy.NewDeploy(fs, "/", nil)

	return deploy.Apply(environmentArg)
}
