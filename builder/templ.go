package builder

const mainTmpl = `package main

##imports##

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/go-spirit/cli"

	"github.com/go-spirit/spirit/cmd"
)

var (
	config = ##config##
	revision = "##revision##"
)

func main() {

	app := cmd.App()
	app.Name = ##Name##

	app.Commands = append(
		app.Commands,
		&cli.Command{
			Name: "metadata",
			Usage: "the metadata of when build this app",
			Subcommands: cli.Commands{
				&cli.Command{
					Name:   "revision",
					Usage:  "show the packages revison",
					Action: showRevision,
				},
				&cli.Command{
					Name:   "config",
					Usage:  "show the configuration while build this app",
					Action: showConfig,
				},
			},
		},
	)

	err := cmd.Init()

	if err != nil {
		logrus.WithError(err).Errorln("run spirit failure")
		return
	}
}

func showConfig(ctx *cli.Context) (err error) {
	fmt.Println(config)
	return
}

func showRevision(ctx *cli.Context) (err error) {
	fmt.Println(revision)
	return
}

`
