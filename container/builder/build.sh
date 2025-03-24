#!/bin/bash

set -e
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
IMAGE_NAME=${IMAGE_NAME:-localhost/ay-dev/slurm-builder:latest}
OS_BASE_IMAGE=${OS_BASE_IMAGE:-docker.io/library/debian:bookworm}
TLS_VERIFY=${TLS_VERIFY:-false}
docker build \
    --build-arg OS_BASE_IMAGE=${OS_BASE_IMAGE} \
    -f $SCRIPT_DIR/Dockerfile \
    -t $IMAGE_NAME $SCRIPT_DIR
