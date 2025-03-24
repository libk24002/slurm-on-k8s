#!/bin/bash

set -e
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
IMAGE=${IMAGE:-localhost/ay-dev/slurm-slurmdbd:latest}
OS_BASE_IMAGE=${OS_BASE_IMAGE:-localhost/ay-dev/slurm-base:latest}
BUILDER_IMAGE=${BUILDER_IMAGE:-localhost/ay-dev/slurm-builder:latest}
TLS_VERIFY=${TLS_VERIFY:-false}
docker build \
    --build-arg OS_BASE_IMAGE=${OS_BASE_IMAGE} \
    --build-arg BUILDER_IMAGE=${BUILDER_IMAGE} \
    -f $SCRIPT_DIR/Dockerfile \
    -t $IMAGE $SCRIPT_DIR
