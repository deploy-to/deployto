package cmd

import (
	"errors"

	"github.com/charmbracelet/huh"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

func createApp(appName string) (string, error) {
	if len(appName) == 0 {
		log.Debug().Msg("Application name not specified")
		err := huh.NewInput().
			Title("Specify the application name").
			Value(&appName).
			Run()
		if err != nil {
			log.Err(err).Msg("huh.NewInput()...Run() error")
		}
		log.Debug().Msg(appName)
		return appName, nil
	}
	log.Debug().Str("appName", appName).Msg("get application name")
	return appName, nil
}

var Create = &cli.Command{
	Name:    "create",
	Aliases: []string{"c"},
	Usage:   "Create a new deployto application from a template",
	Action: func(cCtx *cli.Context) error {
		appName, err := createApp(cCtx.Args().First())
		if err != nil || len(appName) == 0 {
			return errors.Join(errors.New("APPLICATION NAME NOT SPECIFIED"), err)
		}

		for {
			var YN string
			err := huh.NewSelect[string]().
				Title("Add Component?").
				Options(
					huh.NewOption("Yes", "Y"),
					huh.NewOption("No", "N")).
				Value(&YN).
				Run()
			if err != nil {
				log.Err(err).Msg("huh.NewInput()...Run() error")
				return err
			}
			log.Debug().Str("YN", YN).Msg("You choose")

			if YN != "Y" || err != nil {
				return nil
			}
			addComponent("")
		}
	},
}
