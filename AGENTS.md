# AGENTS.md

The purpose of this project is to produce an nginx configuration for redirecting
to search providers based on the phrases used.

You can run the main executable like so:

```sh
go run ./cmd/sea/main.go
```

`./cmd/sea/nginx.conf.tmpl` is the nginx template, and you should keep
`nginx.conf` committed to the repository so people don't have to run the program
to get it. Commit `config.example.toml` as an example configuration file people
can use as per the instructions in `README.md`.

You can re-generate `nginx.go` after changing `nginx.conf.tmpl` or Go sources
like so:

```sh
go run ./cmd/sea/main.go > nginx.conf
```

You can run Go tests like so:

```sh
go test ./...
```

You may validate that the nginx configuration works with the docker compose file
like so:

```sh
# Leave an nginx docker image running on port 57321
docker compose up
# Check a redirect. The hostname comes from `server_name` in `config.toml`
curl -s -D - http://search.localhost:57321/?q=how+can+i -o /dev/null | grep -i '^Location:'
```
