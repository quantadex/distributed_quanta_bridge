#!/usr/bin/env bash

set -e

REPO="691216021071.dkr.ecr.us-east-1.amazonaws.com"
name="quanta-bridge"

# eval $(minikube docker-env)
$(aws ecr get-login --no-include-email --region us-east-1)

VER=git-$(git rev-parse --short HEAD)
#VER=$(date +%Y%m%d)
image="$name:$VER"
dockerfile="Dockerfile"

echo "Build: $name"
if [[ $(aws ecr describe-repositories | grep $name | wc -l) = "0" ]]; then
   aws ecr create-repository --repository-name $name
fi

cd node && go build && cd ..

docker build -t $image -f $dockerfile .
docker tag $image $REPO/$name:latest
docker tag $image $name:latest
#docker push $image
docker push $REPO/$name:latest

echo "... done: $image"
echo

