#! /bin/bash

set -e
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
mkdir -p /tmp/build/containers
podman run --rm \
    -v /tmp/build/containers:/var/lib/containers \
    -v $SCRIPT_DIR/../../..:/code \
    --privileged \
    -e BUILDER_IMAGE=localhost/ay-dev/slurm-builder:latest \
    -e OS_BASE_IMAGE=localhost/ay-dev/slurm-base:latest \
    -it m.zjvis.net/quay.io/containers/buildah:v1.35.4 \
    bash /code/slurm/container/login/build-openmpi.sh

