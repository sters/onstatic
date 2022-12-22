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

The plugin will run on another process. Even if your plugin has broken, it's no problem on onstatic itself. See [plugins](plugins) to understand implementations.

You need implement `plugin.OnstaticPluginServer` on [onstatic/plugin/plugin_grpc.pb.go](onstatic/plugin/plugin_grpc.pb.go) and build it like this:

```shell
go build -o your_plugin_name your_plugin_dir/main.go;
```

Then, you need set binary file into `.onstatic` dir on top of your repository. Like this:

```text
- foobar_repository
    - .onstatic
        - foo
        - bar
    - other_dirs
        - other_files
    - more_other_files
```

onstatic will automatically load plugins and reflect to endpoints.

## Other Information

See [conf/conf.go](conf/conf.go), [onstatic/handler.go](onstatic/handler.go). You can customize it.

