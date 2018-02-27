package git

import (
	"github.com/go-spirit/spirit-builder/builder/fetcher"
	"github.com/gogap/config"
)

type GitFetcher struct {
	conf config.Configuration
}

func init() {
	fetcher.RegisterFetcher("git", NewGitFetcher)
}

func NewGitFetcher(conf config.Configuration) (f fetcher.Fetcher, err error) {

	return
}

func (p *GitFetcher) Fetch(url, revision string, args ...string) (err error) {

	return
}
