package main

import (
	"net/http"

	"github.com/backd-io/backd/backd"
	"github.com/backd-io/backd/internal/constants"
	"github.com/backd-io/backd/internal/pbsessions"
	"github.com/backd-io/backd/internal/rest"
	"github.com/backd-io/backd/internal/structs"
	"github.com/julienschmidt/httprouter"
)

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
		rest.Response(w, nil, err, nil, http.StatusOK, "")
		return
	}

	err = a.mongo.GetOneByIDRBACInterface(session, true, backd.PermissionRead, ps.ByName("domain"), constants.ColUsers, ps.ByName("id"), &user)
	rest.Response(w, user, err, nil, http.StatusOK, "")

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
		rest.Response(w, nil, err, nil, http.StatusOK, "")
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

	err = a.mongo.InsertRBACInterface(session, true, ps.ByName("domain"), constants.ColUsers, &user)
	rest.Response(w, user, err, nil, http.StatusCreated, "")

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
		rest.Response(w, nil, err, nil, http.StatusOK, "")
		return
	}

	err = rest.GetFromBody(r, &user)
	if err != nil {
		rest.BadRequest(w, r, constants.ReasonReadingBody)
		return
	}

	err = a.mongo.GetOneByIDRBACInterface(session, true, backd.PermissionRead, ps.ByName("domain"), constants.ColUsers, ps.ByName("id"), &oldUser)
	if err != nil {
		rest.Response(w, nil, err, nil, http.StatusOK, "")
		return
	}

	// update password if new has been passed
	if user.Password != "" {
		err = user.SetPassword(user.Password)
		if err != nil {
			rest.BadRequest(w, r, constants.ReasonReadingBody)
			return
		}
	}

	// fix metadata
	user.CreatedAt = oldUser.CreatedAt
	user.Owner = oldUser.Owner
	user.ID = oldUser.ID

	user.SetUpdate(session.GetDomainId(), session.GetUserId())

	err = a.mongo.UpdateByIDRBACInterface(session, true, ps.ByName("domain"), constants.ColUsers, ps.ByName("id"), &user)
	rest.Response(w, user, err, nil, http.StatusOK, "")

}

// DELETE /domains/:domain/users/:id
func (a *apiStruct) deleteUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	a.delete(w, r, ps, ps.ByName("domain"), constants.ColUsers)

}
