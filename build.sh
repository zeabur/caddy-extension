#!/usr/bin/env bash

set -eux

docker build -t zeabur/caddy-static --platform linux/amd64 .
docker push zeabur/caddy-static
