#!/usr/bin/env bash

set -eux

docker buildx build \
  -t zeabur/caddy-static:latest \
  -t zeabur/caddy-static:${MAJOR}.${MINOR}.${PATCH} \
  -t zeabur/caddy-static:${MAJOR}.${MINOR} \
  -t zeabur/caddy-static:${MAJOR} \
  --push \
  --platform linux/amd64,linux/arm64 \
  .
