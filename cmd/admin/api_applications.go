package main

import (
	"net/http"

	"github.com/backd-io/backd/backd"
	"github.com/backd-io/backd/internal/constants"
	"github.com/backd-io/backd/internal/db"
	"github.com/backd-io/backd/internal/pbsessions"
	"github.com/backd-io/backd/internal/rest"
	"github.com/backd-io/backd/internal/structs"
	"github.com/julienschmidt/httprouter"
)

// GET /applications
func (a *apiStruct) getApplications(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var (
		query   map[string]interface{}
		sort    []string
		skip    int
		limit   int
		data    []structs.Application
		session *pbsessions.Session
		err     error
	)

	session, err = a.getSession(r)
	if err != nil {
		rest.Response(w, nil, err, nil, http.StatusOK, "")
		return
	}

	query, sort, skip, limit, err = rest.QueryStrings(r)
	if err != nil {
		rest.BadRequest(w, r, constants.ReasonBadQuery)
		return
	}

	err = a.mongo.GetManyRBAC(session, true, backd.PermissionRead, constants.DBBackdApp, constants.ColApplications, query, sort, &data, skip, limit)
	rest.Response(w, data, err, nil, http.StatusOK, "")

}

// GET /applications/:id
func (a *apiStruct) getApplicationByID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var (
		application structs.Application
		session     *pbsessions.Session
		err         error
	)

	// getSession & rbac
	session, err = a.getSession(r)
	if err != nil {
		rest.Response(w, nil, err, nil, http.StatusOK, "")
		return
	}

	// applications reside on backd application database
	err = a.mongo.GetOneByIDRBACInterface(session, false, backd.PermissionRead, constants.DBBackdApp, constants.ColApplications, ps.ByName("id"), &application)
	rest.Response(w, application, err, nil, http.StatusOK, "")

}

// POST /applications
func (a *apiStruct) postApplication(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var (
		application structs.Application
		session     *pbsessions.Session
		err         error
	)

	// getSession & rbac
	session, err = a.getSession(r)
	if err != nil {
		rest.Response(w, nil, err, nil, http.StatusOK, "")
		return
	}

	err = rest.GetFromBody(r, &application)
	if err != nil {
		rest.BadRequest(w, r, constants.ReasonReadingBody)
		return
	}

	application.SetCreate(session.GetDomainId(), session.GetUserId())
	application.ID = db.NewXID().String()

	err = a.mongo.InsertRBACInterface(session, true, constants.DBBackdApp, constants.ColApplications, &application)
	rest.Response(w, application, err, nil, http.StatusCreated, "")

}

// PUT /applications/:id
func (a *apiStruct) putApplication(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var (
		application    structs.Application
		oldApplication structs.Application
		session        *pbsessions.Session
		err            error
	)

	// getSession & rbac
	session, err = a.getSession(r)
	if err != nil {
		rest.Response(w, nil, err, nil, http.StatusOK, "")
		return
	}

	err = rest.GetFromBody(r, &application)
	if err != nil {
		rest.BadRequest(w, r, constants.ReasonReadingBody)
		return
	}

	err = a.mongo.GetOneByIDRBACInterface(session, true, backd.PermissionRead, constants.DBBackdApp, constants.ColApplications, ps.ByName("id"), &oldApplication)
	if err != nil {
		rest.Response(w, nil, err, nil, http.StatusOK, "")
		return
	}

	// fix metadata
	application.CreatedAt = oldApplication.CreatedAt
	application.Owner = oldApplication.Owner
	application.ID = oldApplication.ID

	application.SetUpdate(session.GetDomainId(), session.GetUserId())

	err = a.mongo.UpdateByIDRBACInterface(session, true, constants.DBBackdApp, constants.ColApplications, ps.ByName("id"), &application)
	rest.Response(w, application, err, nil, http.StatusOK, "")

}

// DELETE /applications/:id
func (a *apiStruct) deleteApplication(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	a.delete(w, r, ps, constants.DBBackdApp, constants.ColApplications)

}
