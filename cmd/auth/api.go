package main

import (
	"context"
	"net/http"
	"time"

	"github.com/backd-io/backd/backd"
	"github.com/backd-io/backd/internal/constants"
	"github.com/backd-io/backd/internal/db"
	"github.com/backd-io/backd/internal/instrumentation"
	"github.com/backd-io/backd/internal/pbsessions"
	"github.com/backd-io/backd/internal/rest"
	"github.com/backd-io/backd/internal/structs"
	"github.com/julienschmidt/httprouter"
	"google.golang.org/grpc"
)

type apiStruct struct {
	inst     *instrumentation.Instrumentation
	mongo    *db.Mongo
	sessions *grpc.ClientConn
}

func (a *apiStruct) postSession(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var (
		sessionRequest  structs.SessionRequest
		sessionResponse structs.SessionResponse
		success         bool
		err             error
	)

	err = rest.GetFromBody(r, &sessionRequest)
	if err != nil {
		rest.BadRequest(w, r, constants.ReasonReadingBody)
		return
	}

	success, sessionResponse, err = a.createSession(sessionRequest)
	if err != nil || success == false {
		rest.Unauthorized(w, r)
		return
	}

	rest.Response(w, sessionResponse, err, nil, http.StatusOK, "")
}

func (a *apiStruct) deleteSession(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var (
		result *pbsessions.Result
		cc     pbsessions.SessionsClient
		err    error
	)

	cc = pbsessions.NewSessionsClient(a.sessions)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err = cc.DeleteSession(ctx, &pbsessions.GetSessionRequest{
		Id: r.Header.Get(backd.HeaderSessionID),
	})

	rest.Response(w, result, err, nil, http.StatusOK, "")
}
