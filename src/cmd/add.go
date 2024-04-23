package cmd

import (
	"errors"

	"github.com/charmbracelet/huh"
	"github.com/rs/zerolog/log"

	"github.com/urfave/cli/v2"
)

func getComponentName(componentName string) (string, error) {
	if len(componentName) == 0 {
		log.Debug().Msg("Component name name not specified")
		err := huh.NewInput().
			Title("Specify the component name:").
			Value(&componentName).
			Run()
		if err != nil {
			log.Err(err).Msg("huh.NewInput()...Run() error")
		}
		log.Debug().Msg(componentName)
		return componentName, nil
	}
	log.Debug().Str("componentName", componentName).Msg("get component name")
	return componentName, nil
}

func addComponent(componentName string) error {
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

var Add = &cli.Command{
	Name:    "Add",
	Aliases: []string{"a"},
	Usage:   "Add a new deployto component from a template",
	Action: func(cCtx *cli.Context) error {
		return addComponent(cCtx.Args().First())
	},
}
