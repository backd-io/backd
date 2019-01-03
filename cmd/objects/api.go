package main

import (
	"context"
	"net/http"
	"time"

	"github.com/backd-io/backd/backd"

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

func (a *apiStruct) getSession(r *http.Request) (session *pbsessions.Session, applicationID string, err error) {

	var (
		cc pbsessions.SessionsClient
	)

	if r.Header.Get(rest.HeaderSessionID) == "" || r.Header.Get(rest.HeaderApplicationID) == "" {
		err = rest.ErrUnauthorized
		return
	}

	cc = pbsessions.NewSessionsClient(a.sessions)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	session, err = cc.GetSession(ctx, &pbsessions.GetSessionRequest{
		Id: r.Header.Get(rest.HeaderSessionID),
	})

	applicationID = r.Header.Get(rest.HeaderApplicationID)
	return

}

func (a *apiStruct) getDataID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var (
		data          map[string]interface{}
		session       *pbsessions.Session
		applicationID string
		err           error
	)

	// getSession & rbac
	session, applicationID, err = a.getSession(r)
	if err != nil {
		rest.Response(w, nil, err, nil, http.StatusOK, "")
		return
	}

	data, err = a.mongo.GetOneByIDRBAC(session, backd.PermissionRead, applicationID, ps.ByName("collection"), ps.ByName("id"))
	rest.Response(w, data, err, nil, http.StatusOK, "")
}

func (a *apiStruct) getData(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var (
		data          []map[string]interface{}
		session       *pbsessions.Session
		applicationID string
		err           error
	)

	session, applicationID, err = a.getSession(r)
	if err != nil {
		rest.Response(w, nil, err, nil, http.StatusOK, "")
		return
	}

	err = a.mongo.GetManyRBAC(session, backd.PermissionRead, applicationID, ps.ByName("collection"), map[string]interface{}{}, []string{}, &data, 0, 0)
	rest.Response(w, data, err, nil, http.StatusOK, "")
}

func (a *apiStruct) postData(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var (
		data          map[string]interface{}
		inserted      map[string]interface{}
		session       *pbsessions.Session
		applicationID string
		err           error
	)

	session, applicationID, err = a.getSession(r)
	if err != nil {
		rest.Response(w, nil, err, nil, http.StatusOK, "")
		return
	}

	err = rest.GetFromBody(r, &data)
	if err != nil {
		rest.Response(w, data, err, nil, http.StatusOK, "")
		return
	}

	inserted, err = a.mongo.InsertRBAC(session, applicationID, ps.ByName("collection"), data)
	if err != nil {
		rest.Response(w, nil, err, nil, http.StatusOK, "")
		return
	}

	rest.Response(w, inserted, err, nil, http.StatusOK, rest.Location(ps.ByName("collection"), inserted["_id"].(string)))

}

func (a *apiStruct) putDataID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var (
		data          map[string]interface{}
		updated       map[string]interface{}
		session       *pbsessions.Session
		applicationID string
		err           error
	)

	session, applicationID, err = a.getSession(r)
	if err != nil {
		rest.Response(w, nil, err, nil, http.StatusOK, "")
		return
	}

	err = rest.GetFromBody(r, &data)
	if err != nil {
		rest.Response(w, data, err, nil, http.StatusOK, "")
		return
	}

	updated, err = a.mongo.UpdateByIDRBAC(session, applicationID, ps.ByName("collection"), ps.ByName("id"), data)
	if err != nil {
		rest.Response(w, nil, err, nil, http.StatusOK, "")
		return
	}

	rest.Response(w, updated, err, nil, http.StatusOK, "")

}

func (a *apiStruct) deleteDataID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var (
		session       *pbsessions.Session
		applicationID string
		err           error
	)

	session, applicationID, err = a.getSession(r)
	if err != nil {
		rest.Unauthorized(w, r)
		return
	}

	err = a.mongo.DeleteByIDRBAC(session, applicationID, ps.ByName("collection"), ps.ByName("id"))
	rest.Response(w, nil, err, nil, http.StatusNoContent, "")

}
