package main

import (
	"context"
	"net/http"
	"time"

	"github.com/backd-io/backd/backd"
	"github.com/backd-io/backd/internal/db"
	"github.com/backd-io/backd/internal/instrumentation"
	"github.com/backd-io/backd/internal/pbsessions"
	"github.com/backd-io/backd/internal/rest"
	"github.com/julienschmidt/httprouter"
	"google.golang.org/grpc"
)

type apiStruct struct {
	inst          *instrumentation.Instrumentation
	mongo         *db.Mongo
	sessions      *grpc.ClientConn
	bootstrapCode string
}

func (a *apiStruct) getSession(r *http.Request) (session *pbsessions.Session, err error) {

	var (
		cc pbsessions.SessionsClient
	)

	if r.Header.Get(backd.HeaderSessionID) == "" {
		err = rest.ErrUnauthorized
		return
	}

	cc = pbsessions.NewSessionsClient(a.sessions)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	session, err = cc.GetSession(ctx, &pbsessions.GetSessionRequest{
		Id: r.Header.Get(backd.HeaderSessionID),
	})

	return

}

func (a *apiStruct) delete(w http.ResponseWriter, r *http.Request, ps httprouter.Params, database, collection string) {

	var (
		session *pbsessions.Session
		err     error
	)

	// getSession & rbac
	session, err = a.getSession(r)
	if err != nil {
		rest.ResponseErr(w, err)
		return
	}

	_, err = a.mongo.DeleteByIDRBAC(r.Context(), session, true, database, collection, ps.ByName("id"))
	rest.Response(w, nil, err, http.StatusNoContent, "")

}
