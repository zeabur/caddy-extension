# Caddy with Zeabur extensions

It supports these Zeabur extensions:

- `_headers` file
- `_redirects` file

## Usage

```bash
docker build -t zeabur/caddy-static .
docker run -p 8080:8080 -v $(pwd)/examples:/usr/share/caddy -it zeabur/caddy-static
```
