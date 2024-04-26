package cmd

import (
	"bytes"
	"context"
	"deployto/src/types"
	"deployto/src/yaml"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

type WriterWithReadFrom interface {
	io.Writer
	io.ReaderFrom
}

type ContextWrappedWriter struct {
	w WriterWithReadFrom
	c context.Context
}

type ReadFromResult struct {
	n   int64
	err error
}

func (cww *ContextWrappedWriter) Write(p []byte) (n int, err error) {
	log.Error().Msg("not used")
	var c int
	return c, errors.New("NOT IMPLIMENTED")
}

func (cww *ContextWrappedWriter) ReadFrom(r io.Reader) (n int64, err error) {
	if c, ok := r.(io.Closer); ok {
		ch := make(chan ReadFromResult, 1)
		go func() {
			n, err := cww.w.ReadFrom(r)
			ch <- ReadFromResult{n, err}
		}()

		closed := false
		for {
			select {
			case res := <-ch:
				return res.n, res.err
			case <-cww.c.Done():
				if !closed {
					closed = true
					err := c.Close()
					if err != nil {
						return 0, fmt.Errorf("error closing reader: %v", err)
					}
				}
				time.Sleep(time.Second * 1)
			}
		}

	} else {
		return cww.w.ReadFrom(r)
	}
}

func callJob(cCtx *cli.Context) error {
	path, err := os.Getwd()
	if err != nil {
		log.Error().Err(err).Msg("Get workdir error")
		return err
	}

	// Application
	apps := yaml.Get[types.Job](path)
	if len(apps) != 1 {
		log.Error().Int("len(app)", len(apps)).Str("path", path).Msg("wait one app")
		return errors.New("APP NOT FOUND")
	}
	cCtx.Args()
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Minute)
	defer cancel()
	for _, s := range apps {
		for _, c := range s.Spec.Steps {
			stringSlice := strings.Split(c.Run, "\n")
			for _, cc := range stringSlice {
				if len(cc) > 0 {
					ccc := strings.Split(cc, " ")
					var Stdout, Stderr bytes.Buffer

					c := exec.CommandContext(ctx, ccc[0], ccc[1:]...)
					c.Stderr = &ContextWrappedWriter{&Stderr, ctx}
					c.Stdout = &ContextWrappedWriter{&Stdout, ctx}
					err := c.Run()
					if err != nil {
						timeoutMsg := "Ci Stoped, timeout" + cCtx.String("citimeout")
						log.Error().Err(err).Msg(timeoutMsg)
					}
					fmt.Println(Stderr.String())
					fmt.Println(Stdout.String())
				}

			}

		}

	}

	return nil
}

var timeout int

var Job = &cli.Command{
	Name:    "job",
	Aliases: []string{"j"},
	Usage:   "Manipulate of Job for Current Application",
	Subcommands: []*cli.Command{
		{
			Name:  "run",
			Usage: "run a job",
			Action: func(cCtx *cli.Context) error {
				fmt.Println("new Job Run: ", cCtx.Args().First())
				err := callJob(cCtx)

				if err != nil {
					log.Err(err).Msg("huh.NewInput()...Run() error")
					return err
				}
				return nil
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
					Destination: &timeout,
				},
			},
		},
	},
}
