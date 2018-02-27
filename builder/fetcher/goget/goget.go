package goget

import (
	"fmt"
	"os"
	"path"

	"github.com/sirupsen/logrus"

	"github.com/go-spirit/spirit-builder/builder/fetcher"
	"github.com/go-spirit/spirit-builder/utils"
	"github.com/gogap/config"
)

type GoGetFetcher struct {
	conf config.Configuration
}

func init() {
	fetcher.RegisterFetcher("goget", NewGoGetFetcher)
}

func NewGoGetFetcher(conf config.Configuration) (f fetcher.Fetcher, err error) {

	if conf == nil {
		conf = config.NewConfig()
	}

	ft := &GoGetFetcher{
		conf: conf,
	}

	f = ft

	return
}

func (p *GoGetFetcher) Fetch(url, revision string, args ...string) (err error) {
	cmdArgs := []string{"get"}

	cmdArgs = append(cmdArgs, args...)
	cmdArgs = append(cmdArgs, url)

	result, err := utils.ExecCommand("go", cmdArgs...)

	if err != nil {
		err = fmt.Errorf("fetch repo failure, url: %s, error: %s\n%s\n", url, string(result), err)
		return
	}

	logrus.WithField("fetcher", "goget").WithField("url", url).WithField("revision", revision).Infoln("fetched")

	if len(revision) == 0 {
		return
	}

	gopath := p.conf.GetString("gopath", os.Getenv("GOPATH"))

	checkoutWD := path.Join(gopath, "src", url)
	checkoutArgs := []string{"-C", checkoutWD, "checkout", revision}

	result, err = utils.ExecCommand("git", checkoutArgs...)
	if err != nil {
		err = fmt.Errorf("checkout revision failure, url: %s, error: %s\n%s\n", url, string(result), err)
		return
	}

	logrus.WithField("fetcher", "goget").WithField("url", url).WithField("revision", revision).Infoln("checkout")

	return
}
