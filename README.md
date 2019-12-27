# onstatic

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

See [conf/conf.go](conf/conf.go), [onstatic/handler.go](onstatic/handler.go).

