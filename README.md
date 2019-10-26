# backd

[![Codacy Badge](https://api.codacy.com/project/badge/Grade/a1acac71eff84a05964b0c6f1b395f05)](https://www.codacy.com/manual/fernandezvara/backd_2?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=fernandezvara/backd&amp;utm_campaign=Badge_Grade)
[![Go Report Card](https://goreportcard.com/badge/github.com/fernandezvara/backd)](https://goreportcard.com/report/github.com/fernandezvara/backd)
[![GoDoc](https://godoc.org/github.com/fernandezvara/backd?status.svg)](https://godoc.org/github.com/fernandezvara/backd)
[![CircleCI](https://circleci.com/gh/fernandezvara/backd.svg?style=svg)](https://circleci.com/gh/fernandezvara/backd)

Platform for rapid application development.

```go
println("Work in Progress.")
```

## quick-start

```bash
  docker stack deploy backd --compose-file docker-compose.yml
```

```bash
# retrieve logged in token
kubectl -n kube-system describe secret $(kubectl -n kube-system get secret | grep kubernetes-dashboard-token | awk '{print $1}')

# proxy to dashboard, visit: https://127.0.0.1:8443/
export POD_NAME=$(kubectl get pods -n kube-system -l "app=kubernetes-dashboard,release=dashboard" -o jsonpath="{.items[0].metadata.name}")
kubectl -n kube-system port-forward $POD_NAME 8443:8443
```

## backd-cli

The CLI allows to make most of the actions doable by the using of the API. CLI helps bootstrapping the cluster.

Usage:

```bash

```

## required tools for development

### gox

```bash
go get github.com/mitchellh/gox
```

### govendor

```bash
go get -u github.com/kardianos/govendor
```

## backd go client

[Client Documentation](https://gowalker.org/github.com/fernandezvara/backd/backd)