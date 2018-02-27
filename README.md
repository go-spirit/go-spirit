# spirit-builder
spirit build is a tools for build spirit component



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