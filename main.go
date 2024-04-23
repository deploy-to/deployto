package main

import (
	"deployto/src"
	"deployto/src/cmd"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

func main() {
	src.LogSetting()

	app := &cli.App{
		Name:  "deployto",
		Usage: "just deploy",
		Commands: []*cli.Command{
			cmd.Create,
			cmd.Add,
			cmd.Ci,
		},
		Action: func(cCtx *cli.Context) error {
			return cmd.Deployto(cCtx)
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal().Err(err)
	}
}
