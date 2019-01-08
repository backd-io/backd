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

// GET /domains/:domain/groups
func (a *apiStruct) getGroups(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var (
		query   map[string]interface{}
		sort    []string
		skip    int
		limit   int
		data    []structs.Group
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

	err = a.mongo.GetManyRBAC(session, true, backd.PermissionRead, ps.ByName("domain"), constants.ColGroups, query, sort, &data, skip, limit)
	rest.Response(w, data, err, nil, http.StatusOK, "")

}

// GET /domains/:domain/groups/:id
func (a *apiStruct) getGroupByID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var (
		group   structs.Group
		session *pbsessions.Session
		err     error
	)

	// getSession & rbac
	session, err = a.getSession(r)
	if err != nil {
		rest.Response(w, nil, err, nil, http.StatusOK, "")
		return
	}

	err = a.mongo.GetOneByIDRBACInterface(session, true, backd.PermissionRead, ps.ByName("domain"), constants.ColGroups, ps.ByName("id"), &group)
	rest.Response(w, group, err, nil, http.StatusOK, "")

}

// POST /domains/:domain/groups
func (a *apiStruct) postGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var (
		group   structs.Group
		session *pbsessions.Session
		err     error
	)

	// getSession & rbac
	session, err = a.getSession(r)
	if err != nil {
		rest.Response(w, nil, err, nil, http.StatusOK, "")
		return
	}

	err = rest.GetFromBody(r, &group)
	if err != nil {
		rest.BadRequest(w, r, constants.ReasonReadingBody)
		return
	}

	group.SetCreate(session.GetDomainId(), session.GetUserId())
	group.ID = db.NewXID().String()

	err = a.mongo.InsertRBACInterface(session, true, ps.ByName("domain"), constants.ColGroups, &group)
	rest.Response(w, group, err, nil, http.StatusCreated, "")

}

// PUT /domains/:domain/groups/:id
func (a *apiStruct) putGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var (
		group    structs.Group
		oldGroup structs.Group
		session  *pbsessions.Session
		err      error
	)

	// getSession & rbac
	session, err = a.getSession(r)
	if err != nil {
		rest.Response(w, nil, err, nil, http.StatusOK, "")
		return
	}

	err = rest.GetFromBody(r, &group)
	if err != nil {
		rest.BadRequest(w, r, constants.ReasonReadingBody)
		return
	}

	err = a.mongo.GetOneByIDRBACInterface(session, true, backd.PermissionRead, ps.ByName("domain"), constants.ColGroups, ps.ByName("id"), &oldGroup)
	if err != nil {
		rest.Response(w, nil, err, nil, http.StatusOK, "")
		return
	}

	// fix metadata
	group.CreatedAt = oldGroup.CreatedAt
	group.Owner = oldGroup.Owner
	group.ID = oldGroup.ID

	group.SetUpdate(session.GetDomainId(), session.GetUserId())

	err = a.mongo.UpdateByIDRBACInterface(session, true, ps.ByName("domain"), constants.ColGroups, ps.ByName("id"), &group)
	rest.Response(w, group, err, nil, http.StatusOK, "")

}

// DELETE /domains/:domain/groups/:id
func (a *apiStruct) deleteGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	a.delete(w, r, ps, ps.ByName("domain"), constants.ColGroups)

}
