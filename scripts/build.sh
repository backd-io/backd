#!/usr/bin/env bash

PROJECT_REPO="github.com/fernandezvara/backd"
GIT_COMMIT=$(git rev-parse HEAD)

function build () {
    docker build --force-rm                                     \
        -f scripts/Dockerfile                                   \
        --build-arg ARTIFACT="${1}"                             \
        --build-arg ARCH="${2}"                                 \
        --build-arg API_PORT=$3                                 \
        --build-arg METRICS_PORT=$4                             \
        --build-arg VCS_URL=${PROJECT_REPO}                     \
        --build-arg VCS_REF=${GIT_COMMIT}                       \
        --build-arg BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ") \
        --build-arg VERSION=latest                              \
        -t backd/$1:$2 .
    docker push backd/$1:$2
}

function buildcli () {
    docker build --force-rm                                     \
        -f scripts/Dockerfile.cli                               \
        --build-arg ARTIFACT="${1}"                             \
        --build-arg ARCH="${2}"                                 \
        --build-arg VCS_URL=${PROJECT_REPO}                     \
        --build-arg VCS_REF=${GIT_COMMIT}                       \
        --build-arg BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ") \
        --build-arg VERSION=latest                              \
        -t backd/$1:$2 .
    docker push backd/$1:$2
}

cd $GOPATH/src/${PROJECT_REPO}

for ARCH in {amd64,arm}
do
    for i in {admin,auth,backd,functions,objects,sessions} 
    do 
        gox -os="linux" -arch="${ARCH}" -output="bin/{{.Dir}}_{{.OS}}_{{.Arch}}" github.com/fernandezvara/backd/cmd/$i
        gox -os="darwin" -arch="amd64" -output="bin/{{.Dir}}_{{.OS}}_{{.Arch}}" github.com/fernandezvara/backd/cmd/$i
    done
 
    build "objects" ${ARCH} 8081 8181
    build "sessions" ${ARCH} 8082 8182
    build "auth" ${ARCH} 8083 8183
    build "admin" ${ARCH} 8084 8184
    build "functions" ${ARCH} 8085 8185
    buildcli "backd" ${ARCH} 
done
