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

func (a *apiStruct) internalGetSession(r *http.Request) (session *pbsessions.Session, err error) {

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

func (a *apiStruct) getSession(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	session, err := a.internalGetSession(r)
	rest.Response(w, session, err, http.StatusOK, "")

}

func (a *apiStruct) getMe(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var (
		session *pbsessions.Session
		user    structs.User
		err     error
	)

	session, err = a.internalGetSession(r)
	if err != nil {
		rest.ResponseErr(w, err)
		return
	}

	err = a.mongo.GetOneByIDRBACInterface(r.Context(), session, true, backd.PermissionRead, session.GetDomainId(), constants.ColUsers, session.GetUserId(), &user)
	rest.Response(w, user, err, http.StatusOK, "")

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

	rest.Response(w, sessionResponse, err, http.StatusOK, "")
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

	rest.Response(w, result, err, http.StatusOK, "")
}
