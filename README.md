# onstatic

[![go](https://github.com/sters/go-project-boilerplate/workflows/Go/badge.svg)](https://github.com/sters/go-project-boilerplate/actions?query=workflow%3AGo)
[![codecov](https://codecov.io/gh/sters/go-project-boilerplate/branch/main/graph/badge.svg)](https://codecov.io/gh/sters/go-project-boilerplate)
[![go-report](https://goreportcard.com/badge/github.com/sters/go-project-boilerplate)](https://goreportcard.com/report/github.com/sters/go-project-boilerplate)


onstatic is static page hosting controller.

## Quick Start

Start application from [Releases](https://github.com/sters/onstatic/releases) or yourself.
```
go run cmd/server/main.go
```

Then, do register. Like this.
```
curl -X POST -H "X-ONSTATIC-KEY: onstaticonstaticonstatic" -H "X-ONSTATIC-REPONAME: git@github.com:sters/onstatic.git" localhost:18888/register
```
And you can get SSH Public Key that register to Your git repository's access authentication.

Finally, do pull. Like this:
```
curl -v -X POST -H "X-ONSTATIC-KEY: onstaticonstaticonstatic" -H "X-ONSTATIC-REPONAME: git@github.com:sters/onstatic.git" localhost:18888/pull
```
You can get hashed repository name. Try access `localhost:18888/{Hashed Repository Name}/{your file path}`.



## Other Informations

See [conf/conf.go](conf/conf.go), [onstatic/handler.go](onstatic/handler.go). You can do it.

