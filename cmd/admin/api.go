package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/fernandezvara/backd/backd"
	"github.com/fernandezvara/backd/internal/db"
	"github.com/fernandezvara/backd/internal/instrumentation"
	"github.com/fernandezvara/backd/internal/pbsessions"
	"github.com/fernandezvara/backd/internal/rest"
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

	fmt.Printf("session.admin: %+v\n", session)

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

	err = a.mongo.DeleteByIDRBAC(session, true, database, collection, ps.ByName("id"))
	rest.Response(w, nil, err, http.StatusNoContent, "")

}
