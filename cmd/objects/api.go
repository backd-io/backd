package main

import (
	"github.com/backd-io/backd/internal/db"
	"github.com/backd-io/backd/internal/instrumentation"
	"google.golang.org/grpc"
)

type apiStruct struct {
	inst     *instrumentation.Instrumentation
	mongo    *db.Mongo
	sessions *grpc.ClientConn
}
