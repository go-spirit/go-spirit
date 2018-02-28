package git

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"

	"github.com/go-spirit/spirit-builder/builder/fetcher"
	"github.com/go-spirit/spirit-builder/utils"
	"github.com/gogap/config"
)

type GitFetcher struct {
	conf config.Configuration
}

func init() {
	fetcher.RegisterFetcher("git", NewGitFetcher)
}

func NewGitFetcher(conf config.Configuration) (f fetcher.Fetcher, err error) {
	if conf == nil {
		conf = config.NewConfig()
	}

	ft := &GitFetcher{
		conf: conf,
	}

	f = ft
	return
}

func (p *GitFetcher) Fetch(url, revision string, update bool, repoConf config.Configuration) (err error) {

	repoName := utils.GitRepoName(url)
	if len(repoName) == 0 {
		err = fmt.Errorf("parse git url repository name failure")
		return
	}

	dir := repoConf.GetString("dir")
	if len(dir) == 0 {
		err = fmt.Errorf("unknown dir of repo: %s", url)
		return
	}

	gopath := p.conf.GetString("gopath", os.Getenv("GOPATH"))

	absWorkDir := filepath.Join(gopath, "src", dir)
	absRepoDir := filepath.Join(absWorkDir, repoName)

	err = os.MkdirAll(absWorkDir, 0755)
	if err != nil {
		return
	}

	exist, err := utils.DirExist(absRepoDir)
	if err != nil {
		return
	}

	needClone := !exist

	if len(revision) == 0 {
		revision = "master"
	}

	if needClone {
		err = utils.GitClone(absWorkDir, url, repoConf.GetStringList("args.clone")...)
		if err != nil {
			return
		}
		update = false
		logrus.WithField("fetcher", "git").WithField("url", url).WithField("revision", revision).Infoln("fetched")
	}

	err = utils.GitCheckout(absRepoDir, revision)
	if err != nil {
		return
	}

	logrus.WithField("fetcher", "git").WithField("url", url).WithField("revision", revision).Infoln("checked out")

	if update {
		var deteched bool
		deteched, err = utils.GitDetached(absRepoDir)
		if err != nil {
			return
		}

		if !deteched {
			err = utils.GitPull(absRepoDir, repoConf.GetStringList("args.pull")...)
			if err != nil {
				return
			}
			logrus.WithField("fetcher", "git").WithField("url", url).WithField("revision", revision).Infoln("updated")
		} else {
			logrus.WithField("fetcher", "git").WithField("url", url).WithField("revision", revision).Warnln("git detetched, update skipped")
		}
	}

	return
}
