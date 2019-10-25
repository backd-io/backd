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

// GET /domains/:domain/users
func (a *apiStruct) getUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var (
		query   map[string]interface{}
		sort    []string
		skip    int
		limit   int
		data    []structs.User
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

	err = a.mongo.GetManyRBAC(session, true, backd.PermissionRead, ps.ByName("domain"), constants.ColUsers, query, sort, &data, skip, limit)
	rest.Response(w, data, err, http.StatusOK, "")

}

// GET /domains/:domain/users/:id
func (a *apiStruct) getUserByID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var (
		user    structs.User
		session *pbsessions.Session
		err     error
	)

	// getSession & rbac
	session, err = a.getSession(r)
	if err != nil {
		rest.ResponseErr(w, err)
		return
	}

	err = a.mongo.GetOneByIDRBACInterface(session, true, backd.PermissionRead, ps.ByName("domain"), constants.ColUsers, ps.ByName("id"), &user)
	rest.Response(w, user, err, http.StatusOK, "")

}

// POST /domains/:domain/users
func (a *apiStruct) postUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var (
		user    structs.User
		session *pbsessions.Session
		err     error
	)

	// getSession & rbac
	session, err = a.getSession(r)
	if err != nil {
		rest.ResponseErr(w, err)
		return
	}

	err = rest.GetFromBody(r, &user)
	if err != nil {
		rest.BadRequest(w, r, constants.ReasonReadingBody)
		return
	}

	err = user.SetPassword(user.Password)
	if err != nil {
		rest.BadRequest(w, r, constants.ReasonReadingBody)
		return
	}

	user.SetCreate(session.GetDomainId(), session.GetUserId())
	user.ID = db.NewXID().String()

	err = a.mongo.InsertRBACInterface(session, true, ps.ByName("domain"), constants.ColUsers, &user)
	rest.Response(w, user, err, http.StatusCreated, "")

}

// PUT /domains/:domain/users/:id
func (a *apiStruct) putUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var (
		user    structs.User
		oldUser structs.User
		session *pbsessions.Session
		err     error
	)

	// getSession & rbac
	session, err = a.getSession(r)
	if err != nil {
		rest.ResponseErr(w, err)
		return
	}

	err = rest.GetFromBody(r, &user)
	if err != nil {
		rest.BadRequest(w, r, constants.ReasonReadingBody)
		return
	}

	err = a.mongo.GetOneByIDRBACInterface(session, true, backd.PermissionRead, ps.ByName("domain"), constants.ColUsers, ps.ByName("id"), &oldUser)
	if err != nil {
		rest.ResponseErr(w, err)
		return
	}

	// update password if new has been passed
	if user.Password != "" {
		err = user.SetPassword(user.Password)
		if err != nil {
			rest.ResponseErr(w, err)
			return
		}
	}

	// fix metadata
	user.CreatedAt = oldUser.CreatedAt
	user.CreatedBy = oldUser.CreatedBy
	user.ID = oldUser.ID

	user.SetUpdate(session.GetDomainId(), session.GetUserId())

	err = a.mongo.UpdateByIDRBACInterface(session, true, ps.ByName("domain"), constants.ColUsers, ps.ByName("id"), &user)
	rest.Response(w, user, err, http.StatusOK, "")

}

// DELETE /domains/:domain/users/:id
func (a *apiStruct) deleteUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	a.delete(w, r, ps, ps.ByName("domain"), constants.ColUsers)

}

// GET /domains/:domain/users/:id/groups
func (a *apiStruct) getUserGroups(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var (
		query           map[string]interface{}
		membershipQuery map[string]interface{}
		sort            []string
		ids             []string
		data            []structs.Group
		session         *pbsessions.Session
		err             error
	)

	session, err = a.getSession(r)
	if err != nil {
		rest.ResponseErr(w, err)
		return
	}

	query, sort, _, _, err = rest.QueryStrings(r)
	if err != nil {
		rest.BadRequest(w, r, constants.ReasonBadQuery)
		return
	}

	// get all members of the group
	membershipQuery = map[string]interface{}{
		"u": ps.ByName("id"),
	}

	err = a.mongo.Session().DB(ps.ByName("domain")).C(constants.ColMembership).Find(membershipQuery).Distinct("g", &ids)
	if err != nil {
		rest.ResponseErr(w, err)
		return
	}

	if query == nil {
		query = make(map[string]interface{})
	}

	// query for the users with the id received on the last step
	query["_id"] = map[string]interface{}{
		"$in": ids,
	}

	err = a.mongo.GetManyRBAC(session, true, backd.PermissionRead, ps.ByName("domain"), constants.ColGroups, query, sort, &data, 0, 0)
	rest.Response(w, data, err, http.StatusOK, "")

}
