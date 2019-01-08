package main

import (
	"context"
	"net/http"
	"time"

	"github.com/backd-io/backd/internal/constants"

	"github.com/backd-io/backd/backd"
	"github.com/backd-io/backd/internal/db"
	"github.com/backd-io/backd/internal/instrumentation"
	"github.com/backd-io/backd/internal/pbsessions"
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

	if r.Header.Get(backd.HeaderSessionID) == "" || r.Header.Get(backd.HeaderApplicationID) == "" {
		err = rest.ErrUnauthorized
		return
	}

	cc = pbsessions.NewSessionsClient(a.sessions)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	session, err = cc.GetSession(ctx, &pbsessions.GetSessionRequest{
		Id: r.Header.Get(backd.HeaderSessionID),
	})

	applicationID = r.Header.Get(backd.HeaderApplicationID)
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

	data, err = a.mongo.GetOneByIDRBAC(session, false, backd.PermissionRead, applicationID, ps.ByName("collection"), ps.ByName("id"))
	rest.Response(w, data, err, nil, http.StatusOK, "")
}

func (a *apiStruct) getData(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var (
		query         map[string]interface{}
		sort          []string
		skip          int
		limit         int
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

	query, sort, skip, limit, err = rest.QueryStrings(r)
	if err != nil {
		rest.BadRequest(w, r, constants.ReasonBadQuery)
		return
	}

	err = a.mongo.GetManyRBAC(session, false, backd.PermissionRead, applicationID, ps.ByName("collection"), query, sort, &data, skip, limit)
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

	inserted, err = a.mongo.InsertRBAC(session, false, applicationID, ps.ByName("collection"), data)
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

	updated, err = a.mongo.UpdateByIDRBAC(session, false, applicationID, ps.ByName("collection"), ps.ByName("id"), data)
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

	err = a.mongo.DeleteByIDRBAC(session, false, applicationID, ps.ByName("collection"), ps.ByName("id"))
	rest.Response(w, nil, err, nil, http.StatusNoContent, "")

}
