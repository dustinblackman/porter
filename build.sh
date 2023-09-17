#!/usr/bin/env bash

set -e

TAG="v0.51.9"

DOCKER_DEFAULT_PLATFORM=linux/amd64 docker build -t ghcr.io/dustinblackman/porter:"$TAG" -f ./docker/Dockerfile .
docker push ghcr.io/dustinblackman/porter:"$TAG"
