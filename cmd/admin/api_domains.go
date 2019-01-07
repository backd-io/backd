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

// GET /domains

// GET /domains/:id
func (a *apiStruct) getDomainByID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var (
		domain  structs.Domain
		session *pbsessions.Session
		err     error
	)

	// getSession & rbac
	session, err = a.getSession(r)
	if err != nil {
		rest.Response(w, nil, err, nil, http.StatusOK, "")
		return
	}

	// domains reside on backd application database
	err = a.mongo.GetOneByIDRBACInterface(session, false, backd.PermissionRead, constants.DBBackdApp, constants.ColDomains, ps.ByName("domain"), &domain)
	rest.Response(w, domain, err, nil, http.StatusOK, "")

}

// POST /domains
func (a *apiStruct) postDomain(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var (
		domain  structs.Domain
		session *pbsessions.Session
		err     error
	)

	// getSession & rbac
	session, err = a.getSession(r)
	if err != nil {
		rest.Response(w, nil, err, nil, http.StatusOK, "")
		return
	}

	err = rest.GetFromBody(r, &domain)
	if err != nil {
		rest.BadRequest(w, r, constants.ReasonReadingBody)
		return
	}

	domain.SetCreate(session.GetDomainId(), session.GetUserId())

	err = a.mongo.InsertRBACInterface(session, true, constants.DBBackdApp, constants.ColDomains, &domain)
	rest.Response(w, domain, err, nil, http.StatusCreated, "")

}

// PUT /domains/:id
func (a *apiStruct) putDomain(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var (
		domain    structs.Domain
		oldDomain structs.Domain
		session   *pbsessions.Session
		err       error
	)

	// getSession & rbac
	session, err = a.getSession(r)
	if err != nil {
		rest.Response(w, nil, err, nil, http.StatusOK, "")
		return
	}

	err = rest.GetFromBody(r, &domain)
	if err != nil {
		rest.BadRequest(w, r, constants.ReasonReadingBody)
		return
	}

	err = a.mongo.GetOneByIDRBACInterface(session, true, backd.PermissionRead, constants.DBBackdApp, constants.ColDomains, ps.ByName("id"), &oldDomain)
	if err != nil {
		rest.Response(w, nil, err, nil, http.StatusOK, "")
		return
	}

	// fix metadata
	domain.CreatedAt = oldDomain.CreatedAt
	domain.Owner = oldDomain.Owner
	domain.ID = oldDomain.ID

	domain.SetUpdate(session.GetDomainId(), session.GetUserId())

	err = a.mongo.UpdateByIDRBACInterface(session, true, constants.DBBackdApp, constants.ColDomains, ps.ByName("id"), &domain)
	rest.Response(w, domain, err, nil, http.StatusOK, "")

}

// DELETE /domains/:id
func (a *apiStruct) deleteDomain(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	a.delete(w, r, ps, constants.DBBackdApp, constants.ColDomains)

}
