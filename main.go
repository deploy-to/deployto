package main

import (
	"deployto/src"
	"deployto/src/cmd"
	"deployto/src/types"
	"os"
	"sort"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

func main() {
	stg := types.Settings{}

	app := &cli.App{
		Name:  "deployto",
		Usage: "just deploy",
		Commands: []*cli.Command{
			cmd.Create,
			cmd.Add,
			cmd.Job,
		},
		Action: func(cCtx *cli.Context) error {
			src.LogSetting(stg)
			return cmd.Deployto(cCtx)
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "log-format",
				Aliases:     []string{"lg"},
				Value:       "pretty",
				Usage:       "Log Format: json, pretty",
				Destination: &stg.Logformat,
			},
			&cli.StringFlag{
				Name:        "log-level",
				Aliases:     []string{"ll"},
				Value:       "info",
				Usage:       "Log level: trace, debug, warn, info, fatal, panic, absent, disable",
				Destination: &stg.Loglevel,
			},
		},
	}
	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	if err := app.Run(os.Args); err != nil {
		log.Fatal().Err(err)
	}
}
