package fetcher

import (
	"errors"
	"fmt"
	"github.com/gogap/config"
)

type Fetcher interface {
	Fetch(url, revision string, args ...string) (err error)
}

type NewFetcherFunc func(conf config.Configuration) (Fetcher, error)

var (
	fetcherFuncs = make(map[string]NewFetcherFunc)
)

func RegisterFetcher(name string, fn NewFetcherFunc) (err error) {
	if len(name) == 0 {
		err = errors.New("fetcher name is empty")
		return
	}

	if fn == nil {
		err = errors.New("new fetcher func is nil")
		return
	}

	_, exist := fetcherFuncs[name]

	if exist {
		err = fmt.Errorf("fetcher of %s already exist", name)
		return
	}

	fetcherFuncs[name] = fn

	return
}

func NewFetcher(name string, conf config.Configuration) (f Fetcher, err error) {
	if len(name) == 0 {
		err = errors.New("fetcher name is empty")
		return
	}

	fn, exist := fetcherFuncs[name]

	if !exist {
		err = fmt.Errorf("fetcher of %s not registerd", name)
		return
	}

	f, err = fn(conf)
	return
}
