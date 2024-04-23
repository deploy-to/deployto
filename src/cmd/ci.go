package cmd

import (
	"errors"

	"github.com/charmbracelet/huh"
	"github.com/rs/zerolog/log"

	"github.com/urfave/cli/v2"
)

func callCi(componentName string) error {
	componentName, err := getComponentName(componentName)
	if err != nil || len(componentName) == 0 {
		return errors.Join(errors.New("COMPONENT NAME NOT SPECIFIED"), err)
	}

	var componentType string
	err = huh.NewSelect[string]().
		Title("Choose component type").
		Options(
			huh.NewOption("Web service", "WS"),
			huh.NewOption("Cron job", "CJ"),
			huh.NewOption("Other", "O")).
		Value(&componentType).
		Run()
	if err != nil {
		log.Err(err).Msg("huh.NewSelect()...Run() error")
	}
	log.Debug().Str("componentType", componentType).Msg("You choose")

	return nil
}

var Ci = &cli.Command{
	Name:    "CI",
	Aliases: []string{"a"},
	Usage:   "Build Current Application",
	Action: func(cCtx *cli.Context) error {
		return callCi(cCtx.Args().First())
	},
}
