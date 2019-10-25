package main

import (
	"net/http"

	"github.com/fernandezvara/backd/backd"
	"github.com/fernandezvara/backd/internal/constants"
	"github.com/fernandezvara/backd/internal/db"
	"github.com/fernandezvara/backd/internal/pbsessions"
	"github.com/fernandezvara/backd/internal/rest"
	"github.com/fernandezvara/backd/internal/structs"
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
		rest.ResponseErr(w, err)
		return
	}

	query, sort, skip, limit, err = rest.QueryStrings(r)
	if err != nil {
		rest.BadRequest(w, r, constants.ReasonBadQuery)
		return
	}

	err = a.mongo.GetManyRBAC(session, true, backd.PermissionRead, constants.DBBackdApp, constants.ColApplications, query, sort, &data, skip, limit)
	rest.Response(w, data, err, http.StatusOK, "")

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
		rest.ResponseErr(w, err)
		return
	}

	// applications reside on backd application database
	err = a.mongo.GetOneByIDRBACInterface(session, false, backd.PermissionRead, constants.DBBackdApp, constants.ColApplications, ps.ByName("id"), &application)
	rest.Response(w, application, err, http.StatusOK, "")

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
		rest.ResponseErr(w, err)
		return
	}

	err = rest.GetFromBody(r, &application)
	if err != nil {
		rest.BadRequest(w, r, constants.ReasonReadingBody)
		return
	}

	application.SetCreate(session.GetDomainId(), session.GetUserId())
	application.ID = db.NewXID().String()

	// Create application skeleton
	err = a.mongo.CreateApplicationDatabase(application.ID)
	if err != nil {
		rest.ResponseErr(w, err)
		return
	}

	err = a.mongo.InsertRBACInterface(session, false, constants.DBBackdApp, constants.ColApplications, &application)
	rest.Response(w, application, err, http.StatusCreated, "")

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
		rest.ResponseErr(w, err)
		return
	}

	err = rest.GetFromBody(r, &application)
	if err != nil {
		rest.BadRequest(w, r, constants.ReasonReadingBody)
		return
	}

	err = a.mongo.GetOneByIDRBACInterface(session, false, backd.PermissionRead, constants.DBBackdApp, constants.ColApplications, ps.ByName("id"), &oldApplication)
	if err != nil {
		rest.ResponseErr(w, err)
		return
	}

	// fix metadata
	application.CreatedAt = oldApplication.CreatedAt
	application.CreatedBy = oldApplication.CreatedBy
	application.ID = oldApplication.ID

	application.SetUpdate(session.GetDomainId(), session.GetUserId())

	err = a.mongo.UpdateByIDRBACInterface(session, false, constants.DBBackdApp, constants.ColApplications, ps.ByName("id"), &application)
	rest.Response(w, application, err, http.StatusOK, "")

}

// DELETE /applications/:id
func (a *apiStruct) deleteApplication(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	a.delete(w, r, ps, constants.DBBackdApp, constants.ColApplications)

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

	err = a.mongo.DeleteByIDRBAC(session, false, constants.DBBackdApp, constants.ColApplications, ps.ByName("id"))
	rest.Response(w, nil, err, http.StatusNoContent, "")

}
