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

// GET /applications/:id/functions
func (a *apiStruct) getFunctions(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var (
		query   map[string]interface{}
		sort    []string
		skip    int
		limit   int
		data    []structs.Function
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

	err = a.mongo.GetManyRBAC(session, true, backd.PermissionRead, ps.ByName("id"), constants.ColFunctions, query, sort, &data, skip, limit)
	rest.Response(w, data, err, http.StatusOK, "")

}

// GET /applications/:id/functions/:name
func (a *apiStruct) getFunctionByID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var (
		query    map[string]interface{}
		function map[string]interface{}
		session  *pbsessions.Session
		err      error
	)

	// getSession & rbac
	session, err = a.getSession(r)
	if err != nil {
		rest.ResponseErr(w, err)
		return
	}

	query = map[string]interface{}{
		"name": ps.ByName("name"),
	}

	// applications reside on backd application database
	function, err = a.mongo.GetOneRBAC(session, false, backd.PermissionRead, ps.ByName("id"), constants.ColFunctions, query)
	rest.Response(w, function, err, http.StatusOK, "")

}

// POST /applications/:id/functions
func (a *apiStruct) postFunction(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var (
		function structs.Function
		session  *pbsessions.Session
		err      error
	)

	// getSession & rbac
	session, err = a.getSession(r)
	if err != nil {
		rest.ResponseErr(w, err)
		return
	}

	err = rest.GetFromBody(r, &function)
	if err != nil {
		rest.BadRequest(w, r, constants.ReasonReadingBody)
		return
	}

	function.SetCreate(session.GetDomainId(), session.GetUserId())
	function.ID = db.NewXID().String()

	err = a.mongo.InsertRBACInterface(session, false, ps.ByName("id"), constants.ColFunctions, &function)
	rest.Response(w, function, err, http.StatusCreated, "")

}

// PUT /applications/:id/functions/:name
func (a *apiStruct) putFunction(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var (
		function    structs.Function
		query       map[string]interface{}
		oldFunction map[string]interface{}
		session     *pbsessions.Session
		err         error
	)

	// getSession & rbac
	session, err = a.getSession(r)
	if err != nil {
		rest.ResponseErr(w, err)
		return
	}

	err = rest.GetFromBody(r, &function)
	if err != nil {
		rest.BadRequest(w, r, constants.ReasonReadingBody)
		return
	}

	query = map[string]interface{}{
		"name": ps.ByName("name"),
	}

	oldFunction, err = a.mongo.GetOneRBAC(session, false, backd.PermissionRead, ps.ByName("id"), constants.ColFunctions, query)
	// err = a.mongo.GetOneByIDRBACInterface(session, false, backd.PermissionRead, ps.ByName("id"), constants.ColFunctions, ps.ByName("id"), &oldApplication)
	if err != nil {
		rest.ResponseErr(w, err)
		return
	}

	// fix metadata
	function.FromInterface(oldFunction["meta"].(map[string]interface{}))
	function.ID = oldFunction["_id"].(string)

	function.SetUpdate(session.GetDomainId(), session.GetUserId())

	err = a.mongo.UpdateByIDRBACInterface(session, true, constants.DBBackdApp, constants.ColApplications, ps.ByName("id"), &function)
	rest.Response(w, function, err, http.StatusOK, "")

}

// DELETE /applications/:id/functions/:name
func (a *apiStruct) deleteFunction(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var (
		query   map[string]interface{}
		session *pbsessions.Session
		err     error
	)

	// getSession & rbac
	session, err = a.getSession(r)
	if err != nil {
		rest.ResponseErr(w, err)
		return
	}

	err = a.mongo.DeleteByQueryRBAC(session, false, ps.ByName("id"), constants.ColFunctions, query)
	rest.Response(w, nil, err, http.StatusNoContent, "")

}
