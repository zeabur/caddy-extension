#syntax=docker/dockerfile:1

FROM caddy:2-builder AS builder

COPY --link . /opt/caddy-extensions
WORKDIR /opt/caddy-extensions

RUN --mount=type=cache,target=/go/pkg/mod \
    xcaddy build \
        --with github.com/zeabur/caddy-extension=.

FROM busybox

WORKDIR /usr/share/caddy

COPY --from=builder /opt/caddy-extensions/caddy /usr/bin/caddy
COPY --link Caddyfile /etc/caddy/Caddyfile

EXPOSE 8080
EXPOSE 2019

ENTRYPOINT [ "caddy" ]
CMD [ "run", "--config", "/etc/caddy/Caddyfile", "--adapter", "caddyfile"]
