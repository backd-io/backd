package main

import (
	"net/http"

	"github.com/backd-io/backd/backd"
	"github.com/backd-io/backd/internal/constants"
	"github.com/backd-io/backd/internal/pbsessions"
	"github.com/backd-io/backd/internal/rest"
	"github.com/julienschmidt/httprouter"
)

// GET /objects/:collection/:id
func (a *apiStruct) getObjectID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var (
		data          map[string]interface{}
		session       *pbsessions.Session
		applicationID string
		err           error
	)

	// getSession & rbac
	session, applicationID, err = a.getSession(r)
	if err != nil {
		rest.ResponseErr(w, err)
		return
	}

	data, err = a.mongo.GetOneByIDRBAC(session, false, backd.PermissionRead, applicationID, ps.ByName("collection"), ps.ByName("id"))
	rest.Response(w, data, err, http.StatusOK, "")
}

// GET /objects/:collection
func (a *apiStruct) getObject(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

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
		rest.ResponseErr(w, err)
		return
	}

	query, sort, skip, limit, err = rest.QueryStrings(r)
	if err != nil {
		rest.BadRequest(w, r, constants.ReasonBadQuery)
		return
	}

	err = a.mongo.GetManyRBAC(session, false, backd.PermissionRead, applicationID, ps.ByName("collection"), query, sort, &data, skip, limit)
	rest.Response(w, data, err, http.StatusOK, "")

}

// POST /objects/:collection
func (a *apiStruct) postObject(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var (
		data          map[string]interface{}
		inserted      map[string]interface{}
		session       *pbsessions.Session
		applicationID string
		err           error
	)

	session, applicationID, err = a.getSession(r)
	if err != nil {
		rest.ResponseErr(w, err)
		return
	}

	err = rest.GetFromBody(r, &data)
	if err != nil {
		rest.BadRequest(w, r, constants.ReasonBadQuery)
		return
	}

	inserted, err = a.mongo.InsertRBAC(session, false, applicationID, ps.ByName("collection"), data)
	if err != nil {
		rest.ResponseErr(w, err)
		return
	}

	rest.Response(w, inserted, err, http.StatusCreated, rest.Location(ps.ByName("collection"), inserted["_id"].(string)))

}

// PUT /objects/:collection/:id
func (a *apiStruct) putObjectID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var (
		data          map[string]interface{}
		updated       map[string]interface{}
		session       *pbsessions.Session
		applicationID string
		err           error
	)

	session, applicationID, err = a.getSession(r)
	if err != nil {
		rest.ResponseErr(w, err)
		return
	}

	err = rest.GetFromBody(r, &data)
	if err != nil {
		rest.BadRequest(w, r, constants.ReasonBadQuery)
		return
	}

	updated, err = a.mongo.UpdateByIDRBAC(session, false, applicationID, ps.ByName("collection"), ps.ByName("id"), data)
	if err != nil {
		rest.ResponseErr(w, err)
		return
	}

	rest.Response(w, updated, err, http.StatusOK, "")

}

// DELETE /objects/:collection/:id
func (a *apiStruct) deleteObjectID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

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
	rest.Response(w, nil, err, http.StatusNoContent, "")

}
