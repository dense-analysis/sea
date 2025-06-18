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

## Development Requirements

- Go >= 1.22
- `nginx` >= 1.14

## Contributing

```sh
git clone git@github.com:dense-analysis/sea.git
cd sea
go run ./cmd/sea/main.go
```
