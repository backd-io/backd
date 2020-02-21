#!/usr/bin/env bash

REPO="github.com/backd-io/backd"
PROJECT_REPO="${GOPATH}/src/${REPO}"
GIT_COMMIT=$(git rev-parse HEAD)
BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

function build () {
    # builder & pusher
    cd $PROJECT_REPO
    docker buildx build                                               \
        -f scripts/Dockerfile                                         \
        --platform linux/amd64,linux/arm64,linux/arm/v7               \
        --build-arg ARTIFACT=${1}                                     \
        --build-arg API_PORT=${2}                                     \
        --build-arg METRICS_PORT=${3}                                 \
        --build-arg VCS_URL=${REPO}                                   \
        --build-arg VCS_REF=${GIT_COMMIT}                             \
        --build-arg BUILD_DATE=${BUILD_DATE}                          \
        --build-arg VERSION=latest                                    \
        -t backd/${1}:latest --push .
}

# --platform linux/amd64,linux/arm64,linux/arm/v6,linux/arm/v7  \

function buildcli () {
    # builder & pusher
    cd $PROJECT_REPO
    docker buildx build                                               \
        -f scripts/Dockerfile.cli                                     \
        --platform linux/amd64,linux/arm64,linux/arm/v7               \
        --build-arg ARTIFACT=${1}                                     \
        --build-arg VCS_URL=${REPO}                                   \
        --build-arg VCS_REF=${GIT_COMMIT}                             \
        --build-arg BUILD_DATE=${BUILD_DATE}                          \
        --build-arg VERSION=latest                                    \
        -t backd/${1}:latest --push .
}

build "objects"   8081 8181
build "sessions"  8082 8182
build "auth"      8083 8183
build "admin"     8084 8184
build "functions" 8085 8185
buildcli "backd"

# build for current os/arch
cd $PROJECT_REPO
mkdir -p bin 
rm -rf bin/*

cd $PROJECT_REPO/bin
for i in {admin,auth,backd,functions,objects,sessions}
do 
    go build -o $i github.com/backd-io/backd/cmd/$i
done