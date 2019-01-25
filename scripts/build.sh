#!/usr/bin/env bash

PROJECT_REPO="github.com/backd-io/backd"

cd $GOPATH/src/${PROJECT_REPO}

for i in {admin,auth,backd,functions,objects,sessions} 
do 
    gox -osarch="linux/amd64" -output="bin/{{.Dir}}" github.com/backd-io/backd/cmd/$i
done

GIT_COMMIT=$(git rev-parse HEAD)

function build () {
    docker build --force-rm                                    \
        -f scripts/Dockerfile                                  \
        --build-arg ARTIFACT=$1                                \
        --build-arg API_PORT=$2                                \
        --build-arg METRICS_PORT=$3                            \
        --build-arg VCS_URL=${PROJECT_REPO}                    \
        --build-arg VCS_REF=${GIT_COMMIT}                      \
        --build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
        --build-arg VERSION=latest                             \
        -t backd/$1:latest .
    docker push backd/$1:latest
}

function buildcli () {
    docker build --force-rm                                    \
        -f scripts/Dockerfile.cli                              \
        --build-arg ARTIFACT=$1                                \
        --build-arg VCS_URL=${PROJECT_REPO}                    \
        --build-arg VCS_REF=${GIT_COMMIT}                      \
        --build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
        --build-arg VERSION=latest                             \
        -t backd/$1:latest .
    docker push backd/$1:latest
}

build "objects" 8081 8281
build "sessions" 8082 8282
build "auth" 8083 8283
build "admin" 8084 8284
build "functions" 8085 8285
buildcli "backd"
