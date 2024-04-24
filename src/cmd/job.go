package cmd

import (
	"context"
	"deployto/src/types"
	"deployto/src/yaml"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

const (
	citimeout = "10" //munutes
)

func callJob(cCtx *cli.Context) error {
	path, err := os.Getwd()
	if err != nil {
		log.Error().Err(err).Msg("Get workdir error")
		return err
	}

	// Application
	appPath := yaml.GetAppPath(path)
	apps := yaml.Get[types.Job](appPath)
	if len(apps) != 1 {
		log.Error().Int("len(app)", len(apps)).Str("path", appPath).Msg("wait one app")
		return errors.New("APP NOT FOUND")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(1)*time.Millisecond)
	defer cancel()
	for _, s := range apps {
		for _, c := range s.Spec.Steps {
			stringSlice := strings.Split(c.Run, "\n")
			for _, cc := range stringSlice {
				if err := exec.CommandContext(ctx, cc).Run(); err != nil {
					timeoutMsg := "Ci Stoped, timeout" + cCtx.String("citimeout")
					log.Error().Err(err).Msg(timeoutMsg)
				}
			}

		}

	}

	return nil
}

var Job = &cli.Command{
	Name:    "job",
	Aliases: []string{"j"},
	Usage:   "Manipulate of Job for Current Application",
	Subcommands: []*cli.Command{
		{
			Name:  "run",
			Usage: "run a job",
			Action: func(cCtx *cli.Context) error {
				fmt.Println("new task template: ", cCtx.Args().First())
				callJob(cCtx)
				return nil
			},
		},
	},
	Flags: []cli.Flag{
		&cli.IntFlag{
			Name:        "timeout",
			Aliases:     []string{"t"},
			Value:       10, //minutes
			Usage:       "timeout for ci process",
			DefaultText: "10 minutes",
			Action: func(ctx *cli.Context, v int) error {
				if v >= 100 {
					return fmt.Errorf("Flag timeout value %v out of range[1-100]", v)
				}
				return nil
			},
		},
	},
}
