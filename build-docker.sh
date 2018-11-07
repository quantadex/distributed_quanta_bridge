#!/usr/bin/env bash

# Builds and pushes the docker image.
#
# Usage:
#   ./build-docker.sh [VER=0.0]
#
# The docker image will be tagged with the specified VER followed by the current git hash.

set -e

REPO="691216021071.dkr.ecr.us-east-1.amazonaws.com"
name="quanta-bridge"

if [[ "${1}" == "" ]]; then
    VER="0.0"
else
    VER="${1}"
fi

# if [[ $(git diff --stat) != '' ]]; then
#    echo "The current git directory is dirty. Please stash, commit or remove your changes. (Hint: git diff --stat)"
#    exit 1
# fi

VER=${VER}-$(git rev-parse --short HEAD)

image="$REPO/$name:$VER"
dockerfile="Dockerfile"

echo "Build: $name"
echo "Image: $image"


if [[ $(aws ecr describe-repositories --region us-east-1 | grep $name | wc -l) = "0" ]]; then
   aws ecr create-repository --repository-name $name --region us-east-1
fi

# eval $(minikube docker-env)
$(aws ecr get-login --no-include-email --region us-east-1)

# ./build-linux-binary.sh

cd node && go build && cd ..

docker build -t $image -f $dockerfile .
docker tag $image $REPO/$name:latest
docker tag $image $name:latest
docker push $image
docker push $REPO/$name:latest

echo "... done: $image"
echo $image



