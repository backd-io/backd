package main

import (
	"context"
	"net/http"
	"time"

	"github.com/backd-io/backd/internal/pbsessions"

	"github.com/backd-io/backd/internal/db"
	"github.com/backd-io/backd/internal/instrumentation"
	"github.com/backd-io/backd/internal/rest"
	"github.com/julienschmidt/httprouter"
	"google.golang.org/grpc"
)

type apiStruct struct {
	inst     *instrumentation.Instrumentation
	mongo    *db.Mongo
	sessions *grpc.ClientConn
}

func (a *apiStruct) getSession(r *http.Request) (*pbsessions.Session, error) {

	var (
		cc pbsessions.SessionsClient
	)

	cc = pbsessions.NewSessionsClient(a.sessions)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return cc.GetSession(ctx, &pbsessions.GetSessionRequest{
		Id: r.Header.Get(rest.HeaderSessionID),
	})

}

func (a *apiStruct) getDataOne(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var (
		data map[string]interface{}
		err  error
	)

	// getSession & rbac

	rest.Response(w, data, err, nil, http.StatusOK, "")
}
