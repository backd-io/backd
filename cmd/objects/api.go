package main

import (
	"github.com/fernandezvara/backd/internal/db"
	"github.com/fernandezvara/backd/internal/instrumentation"
	"google.golang.org/grpc"
)

type apiStruct struct {
	inst     *instrumentation.Instrumentation
	mongo    *db.Mongo
	sessions *grpc.ClientConn
}
