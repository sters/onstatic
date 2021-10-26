# onstatic

[![go](https://github.com/sters/onstatic/workflows/Go/badge.svg)](https://github.com/sters/onstatic/actions?query=workflow%3AGo)
[![codecov](https://codecov.io/gh/sters/onstatic/branch/master/graph/badge.svg)](https://codecov.io/gh/sters/onstatic)
[![go-report](https://goreportcard.com/badge/github.com/sters/onstatic)](https://goreportcard.com/report/github.com/sters/onstatic)


The onstatic is static page hosting controller.

## Quick Start

Start application from [Releases](https://github.com/sters/onstatic/releases) or yourself.

```shell
go run cmd/server/main.go
```

Then, do register. Like this.

```shell
curl -X POST -H "X-ONSTATIC-KEY: onstaticonstaticonstatic" -H "X-ONSTATIC-REPONAME: git@github.com:sters/onstatic.git" localhost:18888/register
```

And you can get SSH Public Key that register to Your git repository's access authentication.

Finally, do pull. Like this:

```shell
curl -v -X POST -H "X-ONSTATIC-KEY: onstaticonstaticonstatic" -H "X-ONSTATIC-REPONAME: git@github.com:sters/onstatic.git" localhost:18888/pull
```

You can get hashed repository name. Try access `localhost:18888/{Hashed Repository Name}/{your file path}`.


## Plugins

See [plugins/example](plugins/example).

You need implement `EntryPoint` on [onstatic/plugin/api.go](onstatic/plugin/api.go).

Also you need set *.so file into `.onstatic` dir on top of your repository. Like this:

```text
- foobar_repository
    - .onstatic
        - foo.so
        - bar.so
    - other_dirs
        - other_files
    - more_other_files
```

## Other Informations

See [conf/conf.go](conf/conf.go), [onstatic/handler.go](onstatic/handler.go). You can do it.

