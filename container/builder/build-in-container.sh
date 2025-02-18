#! /bin/bash

set -e
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
mkdir -p /tmp/build/containers
podman run --rm \
    -v /tmp/build/containers:/var/lib/containers \
    -v $SCRIPT_DIR/../../..:/code \
    --privileged \
    -e OS_BASE_IMAGE=docker.io/library/debian:bookworm \
    -e IMAGE_NAME=localhost/ay-dev/slurm-builder:latest \
    -it quay.io/containers/buildah:v1.35.4 \
    bash /code/slurm/container/builder/build.sh
