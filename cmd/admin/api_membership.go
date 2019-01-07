package main

import (
	"net/http"

	"github.com/backd-io/backd/internal/constants"
	"github.com/backd-io/backd/internal/pbsessions"
	"github.com/backd-io/backd/internal/rest"
	"github.com/backd-io/backd/internal/structs"
	"github.com/julienschmidt/httprouter"
)

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
		rest.Response(w, nil, err, nil, http.StatusOK, "")
		return
	}

	membership.GroupID = ps.ByName("id")
	membership.UserID = ps.ByName("user_id")

	err = a.mongo.InsertRBACInterface(session, true, ps.ByName("domain"), constants.ColMembership, &membership)
	rest.Response(w, nil, err, nil, http.StatusNoContent, "")

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
		rest.Response(w, nil, err, nil, http.StatusOK, "")
		return
	}

	err = a.mongo.DeleteByQueryRBAC(session, true, ps.ByName("domain"), constants.ColMembership, map[string]interface{}{
		"g": ps.ByName("id"),
		"u": ps.ByName("user_id"),
	})

	rest.Response(w, nil, err, nil, http.StatusNoContent, "")

}
