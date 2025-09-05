#!/usr/bin/env bash

set -eux

docker build -t zeabur/caddy-static --push --platform linux/amd64 .
