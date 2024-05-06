package main

import (
	"deployto/src/cmd"
	"os"
	"sort"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "deployto",
		Usage: "just deploy",
		Commands: []*cli.Command{
			cmd.Create,
			cmd.Add,
			cmd.Job,
		},
		Before: func(ctx *cli.Context) error {
			cmd.LogSetting(ctx.String("log-format"), ctx.String("log-level"))
			return cmd.LoadYamlConfig(ctx)
		},
		Action: cmd.Deployto,
		Flags:  cmd.Flags,
	}
	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	if err := app.Run(os.Args); err != nil {
		log.Fatal().Err(err).Msg("deployto error")
	}
}
