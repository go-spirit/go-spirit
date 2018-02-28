package git

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

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

func (p *GitFetcher) gitPorjectName(gitUrl string) string {
	found := regexp.MustCompile(`([^/]+)\.git$`).FindAllString(gitUrl, 1)

	if len(found) == 0 {
		return ""
	}

	return strings.TrimSuffix(found[0], ".git")
}

func (p *GitFetcher) Fetch(url, revision string, update bool, repoConf config.Configuration) (err error) {
	// git -C gopath options clone url

	repoName := p.gitPorjectName(url)
	if len(repoName) == 0 {
		err = fmt.Errorf("parse git url repository name failure")
		return
	}

	cloneDir := repoConf.GetString("dir")
	if len(cloneDir) == 0 {
		err = fmt.Errorf("unknown dir of repo: %s", url)
		return
	}

	gopath := p.conf.GetString("gopath", os.Getenv("GOPATH"))

	absCloneDir := filepath.Join(gopath, "src", cloneDir)

	needClone := false

	fi, errF := os.Stat(filepath.Join(absCloneDir, repoName))
	if errF != nil {
		if os.IsNotExist(errF) {
			needClone = true
		}
	} else if !fi.IsDir() {
		err = fmt.Errorf("%s is not a dir", absCloneDir)
		return
	}

	repoDir := filepath.Join(absCloneDir, repoName)

	if len(revision) == 0 {
		revision = "master"
	}

	if !needClone {
		err = p.checkout(repoDir, revision)
		if err != nil {
			return
		}

		var deteched bool
		deteched, err = p.isRepoDetached(repoDir, absCloneDir)
		if err != nil {
			return
		}

		// the deteched repo are not on branch, could not be pull
		if deteched {
			update = false
		}
	}

	if !update && !needClone {
		return
	}

	var cmdArgs []string
	var args []string

	if needClone {
		cmdArgs = []string{"-C", absCloneDir, "clone"}
		args = repoConf.GetStringList("args.clone")
		os.MkdirAll(absCloneDir, 0755)
	} else {
		cmdArgs = []string{"-C", filepath.Join(absCloneDir, repoName), "pull"}
		args = repoConf.GetStringList("args.pull")
	}

	cmdArgs = append(cmdArgs, args...)

	if needClone {
		cmdArgs = append(cmdArgs, url)
	}

	// execute command clone or pull
	result, err := utils.ExecCommand("git", cmdArgs...)

	if err != nil {
		err = fmt.Errorf("fetch repo failure, url: %s, error: %s\n%s\n", url, string(result), err)
		return
	}

	logrus.WithField("fetcher", "git").WithField("url", url).WithField("revision", revision).Infoln("fetched")

	if needClone {
		err = p.checkout(repoDir, revision)
		if err != nil {
			return
		}
	}

	return
}

func (p *GitFetcher) checkout(repoDir, revision string) (err error) {
	checkoutArgs := []string{"-C", repoDir, "checkout", revision}

	var result []byte
	result, err = utils.ExecCommand("git", checkoutArgs...)
	if err != nil {
		err = fmt.Errorf("checkout revision failure, dir: %s, error: %s\n%s\n", repoDir, string(result), err)
		return
	}

	logrus.WithField("fetcher", "git").WithField("dir", repoDir).WithField("revision", revision).Infoln("checkout")

	return
}

func (p *GitFetcher) isRepoDetached(wkdir, dir string) (bool, error) {
	result, err := utils.ExecCommand("git", "-C", wkdir, "status", "-b")

	if err != nil {
		err = fmt.Errorf("get git status failure: %s", err.Error())
		return false, err
	}

	return strings.Index(string(result), "detached") > 0, nil
}
