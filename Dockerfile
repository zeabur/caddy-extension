#syntax=docker/dockerfile:1

FROM caddy:2-builder AS builder

COPY --link . /opt/caddy-extensions
WORKDIR /opt/caddy-extensions

RUN xcaddy build \
    --with github.com/zeabur/caddy-extension=.

FROM caddy:2

WORKDIR /usr/share/caddy

COPY --from=builder /opt/caddy-extensions/caddy /usr/bin/caddy
COPY --link Caddyfile /etc/caddy/Caddyfile
