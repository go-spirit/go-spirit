# spirit-builder
spirit build is a tools for build spirit component

## Install spirit-builder

#### install

```bash
go get github.com/go-spirit/spirit-builder
```

> if do not install, please run 

```bash
go install github.com/go-spirit/spirit-builder
```


#### update

```bash
go get -u github.com/go-spirit/spirit-builder
```



## Run todo project

#### pull project source

```bash
> spirit-builder pull --config build.conf
INFO[0000] fetched    fetcher=goget revision=master url=github.com/spirit-component/examples/todo
INFO[0000] checkout   fetcher=goget revision=master url=github.com/spirit-component/examples/todo
INFO[0000] fetched    fetcher=goget revision=master url=github.com/spirit-component/postapi
INFO[0000] checkout   fetcher=goget revision=master url=github.com/spirit-component/postapi
```

#### build project

```bash
spirit-builder build --config build.conf
INFO[0000] building                                      project=todo
```


`build.conf`

> use goget fetcher

```hocon
# project
todo  {

	# import packages
	packages = ["github.com/spirit-component/examples/todo", "github.com/spirit-component/postapi"]

	build-args = []

	fetchers {
		git {
			gopath = ${GOPATH}
		}
		goget {
			gopath = ${GOPATH}
		}
	}


	# the dependencies
	repos = {
		todo {
			fetcher = goget
			args = ["-v"]
			url = "github.com/spirit-component/examples/todo"
			revision = master
		}

		postapi {
			fetcher = goget
			args = ["-v"]
			url = "github.com/spirit-component/postapi"
			revision = master
		}
	}
}
```


> use git fetcher

```hocon
# project
todo  {

	# import packages
	packages = ["github.com/spirit-component/examples/todo", "github.com/spirit-component/postapi"]

	build-args = []

	fetchers {
		git {
			gopath = ${GOPATH}
		}
		goget {
			gopath = ${GOPATH}
		}
	}


	# the dependencies
	repos = {
		todo {
			fetcher = git
			args = {
				clone = []
				pull = []
			}
			url = "git@github.com:spirit-component/examples.git"
			revision = master
			dir = "github.com/spirit-component"
		}

		postapi {
			fetcher = git
			args = {
				clone = []
				pull = []
			}
			url = "git@github.com:spirit-component/postapi.git"
			revision = master
			dir = "github.com/spirit-component"
		}
	}
}
```