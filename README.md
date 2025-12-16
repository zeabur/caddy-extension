# Caddy with Zeabur extensions

It supports these Zeabur extensions:

- `_headers` file
- `_redirects` file

## Usage

```bash
docker build -t zeabur/caddy-static .
docker run -p 8080:8080 -v $(pwd)/examples/caddy:/usr/share/caddy -it zeabur/caddy-static
```

## Publish

```bash
MAJOR=1 MINOR=0 PATCH=0 bash build.sh
```

## Test

You should build the `zeabur/caddy-static` image first.

```bash
go test -v ./e2etest
```

It checks if the Zeabur extensions are working correctly.
