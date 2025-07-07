# Search Engine Adapter

SEA (Search Engine Adapter) is a lightweight compiled `nginx` proxy
configuration that dynamically routes incoming search queries to Google,
ChatGPT, or Wikipedia based on keyword matching. The configuration is compiled
via a Go program, and stored in this repository so it can be easily copied and
applied to `nginx`.

## Features

- **Dynamic Routing**: Routes queries to Google Search, ChatGPT API, or
  Wikipedia based on the semantic meanings of words and phrases.
- **Single HTTP Endpoint**: One URL for all search needs, find everything
  quickly online, and get the best answers from the best sources.
- **Simple Configuration**: Ships as a ready-to-use NGINX `.conf` file.
- **Lightweight & Fast**: Minimal dependencies for high throughput. Very simple
  security profile, as there is no web app to exploit.
- **Optimized Matching**: Regex rules are evaluated once per request using
  NGINX's `map` directive. Enable `pcre_jit` for the best performance.

## Installation

To install into `nginx`, you may include the `nginx.conf` provided in either
`/etc/nginx/conf.d/*.conf` or in `/etc/nginx/sites-enabled/`.

If you run the Go program to generate a more heavily customised configuration,
you can do the following:

```sh
# Copy the config
cp config-example.toml config.toml
# Edit the config
# You might want to change `listen` or `server_name`
vim config.toml
go run cmd/sea/main.go > sea.conf
```

### TLS / Let's Encrypt

If you want to run SEA behind HTTPS using Certbot, set `listen_ssl`,
`ssl_certificate`, and `ssl_certificate_key` in `config.toml`.  Enable the
`letsencrypt` and `redirect_http` flags to generate a second server block for
port 80 that redirects to HTTPS.  The resulting configuration can be safely
included in your `nginx` setup, and Certbot may modify only the certificate
paths when certificates renew.

## Development Requirements

- Go >= 1.22
- `nginx` >= 1.14

## Contributing

```sh
git clone git@github.com:dense-analysis/sea.git
cd sea
go run ./cmd/sea/main.go
```
