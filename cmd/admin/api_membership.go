package main

import (
	"net/http"

	"github.com/fernandezvara/backd/backd"
	"github.com/fernandezvara/backd/internal/constants"
	"github.com/fernandezvara/backd/internal/pbsessions"
	"github.com/fernandezvara/backd/internal/rest"
	"github.com/fernandezvara/backd/internal/structs"
	"github.com/julienschmidt/httprouter"
)

// GET /domains/:domain/groups/:id/members
func (a *apiStruct) getMembers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var (
		query           map[string]interface{}
		membershipQuery map[string]interface{}
		sort            []string
		ids             []string
		data            []structs.User
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
		"g": ps.ByName("id"),
	}

	err = a.mongo.Session().DB(ps.ByName("domain")).C(constants.ColMembership).Find(membershipQuery).Distinct("u", &ids)
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

	err = a.mongo.GetManyRBAC(session, true, backd.PermissionRead, ps.ByName("domain"), constants.ColUsers, query, sort, &data, 0, 0)
	rest.Response(w, data, err, http.StatusOK, "")

}

// PUT /domains/:domain/groups/:id/members/:user_id
func (a *apiStruct) putMembership(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var (
		membership structs.Membership
		session    *pbsessions.Session
		err        error
	)

	// getSession & rbac
	session, err = a.getSession(r)
	if err != nil {
		rest.ResponseErr(w, err)
		return
	}

	membership.GroupID = ps.ByName("id")
	membership.UserID = ps.ByName("user_id")

	err = a.mongo.InsertRBACInterface(session, true, ps.ByName("domain"), constants.ColMembership, &membership)
	rest.Response(w, nil, err, http.StatusNoContent, "")

}

// DELETE /domains/:domain/groups/:id/members/:user_id
func (a *apiStruct) deleteMembership(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

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

	err = a.mongo.DeleteByQueryRBAC(session, true, ps.ByName("domain"), constants.ColMembership, map[string]interface{}{
		"g": ps.ByName("id"),
		"u": ps.ByName("user_id"),
	})

	rest.Response(w, nil, err, http.StatusNoContent, "")

}
