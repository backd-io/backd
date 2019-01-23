#!/usr/bin/env bash

PROJECT_REPO="github.com/backd-io/backd"

cd $GOPATH/src/${PROJECT_REPO}

GIT_COMMIT=$(git rev-parse HEAD)

for i in {admin,auth,backd,functions,objects,sessions} ; do
    docker build --force-rm                                    \
        -f scripts/Dockerfile                                  \
        --build-arg ARTIFACT=${i}                              \
        --build-arg VCS_URL=${PROJECT_REPO}                    \
        --build-arg VCS_REF=${GIT_COMMIT}                      \
        --build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
        --build-arg VERSION=latest                             \
        -t backd/$i:latest .
    docker push backd/${i}:latest
done