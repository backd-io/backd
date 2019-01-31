package main

import (
	"net/http"

	"github.com/backd-io/backd/internal/db"

	"github.com/backd-io/backd/backd"
	"github.com/backd-io/backd/internal/constants"
	"github.com/backd-io/backd/internal/pbsessions"
	"github.com/backd-io/backd/internal/rest"
	"github.com/backd-io/backd/internal/structs"
	"github.com/julienschmidt/httprouter"
)

// GET /domains
func (a *apiStruct) getDomains(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var (
		query   map[string]interface{}
		sort    []string
		skip    int
		limit   int
		data    []structs.Domain
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

	err = a.mongo.GetManyRBAC(session, true, backd.PermissionRead, constants.DBBackdApp, constants.ColDomains, query, sort, &data, skip, limit)
	rest.Response(w, data, err, http.StatusOK, "")

}

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
		rest.ResponseErr(w, err)
		return
	}

	// domains reside on backd application database
	err = a.mongo.GetOneByIDRBACInterface(session, false, backd.PermissionRead, constants.DBBackdApp, constants.ColDomains, ps.ByName("domain"), &domain)
	rest.Response(w, domain, err, http.StatusOK, "")

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
		rest.ResponseErr(w, err)
		return
	}

	err = rest.GetFromBody(r, &domain)
	if err != nil {
		rest.BadRequest(w, r, constants.ReasonReadingBody)
		return
	}

	domain.SetCreate(session.GetDomainId(), session.GetUserId())
	domain.ID = db.NewXID().String()

	err = a.mongo.CreateDomainDatabase(domain.ID)
	if err != nil {
		rest.ResponseErr(w, err)
		return
	}

	err = a.mongo.InsertRBACInterface(session, true, constants.DBBackdApp, constants.ColDomains, &domain)
	rest.Response(w, domain, err, http.StatusCreated, "")

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
		rest.ResponseErr(w, err)
		return
	}

	err = rest.GetFromBody(r, &domain)
	if err != nil {
		rest.BadRequest(w, r, constants.ReasonReadingBody)
		return
	}

	err = a.mongo.GetOneByIDRBACInterface(session, true, backd.PermissionRead, constants.DBBackdApp, constants.ColDomains, ps.ByName("id"), &oldDomain)
	if err != nil {
		rest.ResponseErr(w, err)
		return
	}

	// fix metadata
	domain.CreatedAt = oldDomain.CreatedAt
	domain.CreatedBy = oldDomain.CreatedBy
	domain.ID = oldDomain.ID

	domain.SetUpdate(session.GetDomainId(), session.GetUserId())

	err = a.mongo.UpdateByIDRBACInterface(session, true, constants.DBBackdApp, constants.ColDomains, ps.ByName("id"), &domain)
	rest.Response(w, domain, err, http.StatusOK, "")

}

// DELETE /domains/:id
func (a *apiStruct) deleteDomain(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	a.delete(w, r, ps, constants.DBBackdApp, constants.ColDomains)

}
