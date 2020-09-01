#!/bin/bash
DOCKER_USER=admin
DOCKER_PASS=XgE+CC5Vyo2n
DOCKER_HOST=hub.htres.cn

BUILD_TARGET=$1

BUILD_TARGET=${BUILD_TARGET:-dev}
SRC_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

BUILD_MODULE=$2
BUILD_MODULE=${BUILD_MODULE:-all}

if [[ ! $BUILD_TARGET =~ ^(dev|test|release)$ ]]; then
  echo "usage: package [dev | test | release] [albcp | lbmc]"
  exit 1
fi

pushd $SRC_DIR/..
VERSION=$(make version)
docker login --username=${DOCKER_USER} --password=${DOCKER_PASS} ${DOCKER_HOST}

case "${BUILD_MODULE}" in
"all")
docker build -f $SRC_DIR/Dockerfile.alb-cp --build-arg BUILD_TARGET=${BUILD_TARGET} -t ${DOCKER_HOST}/pub/albcp:${VERSION} .
docker push ${DOCKER_HOST}/pub/albcp:${VERSION}
docker build -f $SRC_DIR/Dockerfile.lbmc --build-arg BUILD_TARGET=${BUILD_TARGET} -t ${DOCKER_HOST}/pub/lbmc:${VERSION} .
docker push ${DOCKER_HOST}/pub/lbmc:${VERSION}
;;

"albcp")
docker build -f $SRC_DIR/Dockerfile.alb-cp --build-arg BUILD_TARGET=${BUILD_TARGET} -t ${DOCKER_HOST}/pub/albcp:${VERSION} .
docker push ${DOCKER_HOST}/pub/albcp:${VERSION}
;;

"lbmc")
docker build -f $SRC_DIR/Dockerfile.lbmc --build-arg BUILD_TARGET=${BUILD_TARGET} -t ${DOCKER_HOST}/pub/lbmc:${VERSION} .
docker push ${DOCKER_HOST}/pub/lbmc:${VERSION}
esac

popd