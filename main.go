package main

import (
	"errors"
	"math/rand"
	"os"
	"time"

	"github.com/go-spirit/spirit-builder/builder"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

import (
	_ "github.com/go-spirit/spirit-builder/builder/fetcher/git"
	_ "github.com/go-spirit/spirit-builder/builder/fetcher/goget"
)

func main() {
	app := cli.NewApp()

	app.Commands = cli.Commands{
		cli.Command{
			Name:   "pull",
			Action: pull,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "config",
					Usage: "config file",
				},
				cli.StringSliceFlag{
					Name:  "name",
					Usage: "project name",
				},
			},
		},
		cli.Command{
			Name:   "build",
			Action: build,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "config",
					Usage: "config file",
				},
				cli.StringSliceFlag{
					Name:  "name",
					Usage: "project name",
				},
			},
		},
	}

	rand.Seed(time.Now().UnixNano())

	err := app.Run(os.Args)

	if err != nil {
		logrus.Errorln(err)
		return
	}

	return
}

func build(ctx *cli.Context) (err error) {

	configfile := ctx.String("config")
	if len(configfile) == 0 {
		err = errors.New("config file not specified")
		return
	}

	builder, err := builder.NewBuilder(
		builder.ConfigFile(configfile),
	)

	if err != nil {
		return
	}

	buildNames := ctx.StringSlice("name")

	if len(buildNames) == 0 {
		buildNames = builder.ListProject()
	}

	if len(buildNames) == 0 {
		return
	}

	err = builder.Build(buildNames...)

	if err != nil {
		return
	}

	return
}

func pull(ctx *cli.Context) (err error) {

	configfile := ctx.String("config")
	if len(configfile) == 0 {
		err = errors.New("config file not specified")
		return
	}

	builder, err := builder.NewBuilder(
		builder.ConfigFile(configfile),
	)

	if err != nil {
		return
	}

	buildNames := ctx.StringSlice("name")

	if len(buildNames) == 0 {
		buildNames = builder.ListProject()
	}

	if len(buildNames) == 0 {
		return
	}

	err = builder.Pull(buildNames...)

	if err != nil {
		return
	}

	return
}
