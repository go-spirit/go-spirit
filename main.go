package main

import (
	"errors"
	"math/rand"
	"os"
	"time"

	"github.com/go-spirit/go-spirit/builder"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

import (
	_ "github.com/go-spirit/go-spirit/builder/fetcher/git"
	_ "github.com/go-spirit/go-spirit/builder/fetcher/goget"
)

func main() {
	app := cli.NewApp()
	app.Usage = "spirit project builder"

	app.Commands = cli.Commands{
		cli.Command{
			Name:   "pull",
			Usage:  "pull project repositories",
			Action: pull,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "config",
					Usage: "config file",
				},
				cli.BoolFlag{
					Name:  "update",
					Usage: "update repo if exist",
				},
				cli.StringSliceFlag{
					Name:  "name",
					Usage: "project name",
				},
			},
		},
		cli.Command{
			Name:   "build",
			Usage:  "build project",
			Action: build,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "config, c",
					Usage: "config file",
				},
				cli.StringSliceFlag{
					Name:  "name, n",
					Usage: "project name",
				},
			},
		},
	}

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "log-level",
			Usage: "debug, info, warn, error, fatal, panic",
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

func initLogLevel(ctx *cli.Context) (err error) {
	strlvl := ctx.Parent().String("log-level")

	if len(strlvl) == 0 {
		return
	}

	lvl, err := logrus.ParseLevel(strlvl)
	if err != nil {
		return
	}

	logrus.SetLevel(lvl)

	return
}

func build(ctx *cli.Context) (err error) {

	err = initLogLevel(ctx)
	if err != nil {
		return
	}

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

	err = initLogLevel(ctx)
	if err != nil {
		return
	}

	configfile := ctx.String("config")
	if len(configfile) == 0 {
		err = errors.New("config file not specified")
		return
	}

	update := ctx.Bool("update")

	builder, err := builder.NewBuilder(
		builder.ConfigFile(configfile),
		builder.NeedUpdate(update),
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
