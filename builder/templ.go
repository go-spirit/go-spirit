package builder

const mainTmpl = `package main

##imports##

import (
	"github.com/sirupsen/logrus"

	"github.com/go-spirit/spirit/cmd"
)

func main() {

	cmd.App()

	err := cmd.Init()

	if err != nil {
		logrus.WithError(err).Errorln("run spirit failure")
		return
	}
}
`
