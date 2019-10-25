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

// GET /objects/:collection/:id/:relation/:direction
func (a *apiStruct) getObjectIDRelations(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var (
		relations     []structs.Relation       // database relations
		objects       []map[string]interface{} // objects related the requester can access to
		query         map[string]interface{}
		session       *pbsessions.Session
		applicationID string
		err           error
	)

	session, applicationID, err = a.getSession(r)
	if err != nil {
		rest.ResponseErr(w, err)
		return
	}

	switch ps.ByName("direction") {
	case "in": // incoming relations */* -> relation_name -> collection:/id
		// if requester can not read the item stop here, it's unauthorized
		if a.mongo.Can(session, false, applicationID, ps.ByName("collection"), ps.ByName("id"), backd.PermissionRead) == false {
			rest.Unauthorized(w, r)
			return
		}

		query = map[string]interface{}{
			"dst": ps.ByName("collection"),
			"did": ps.ByName("id"),
			"rel": ps.ByName("relation"),
		}

		err = a.mongo.GetAll(applicationID, constants.ColRelations, query, []string{}, &relations)
		if err != nil {
			rest.ResponseErr(w, err)
			return
		}

		// add only those relations the user can manage
		for _, relation := range relations {
			var obj map[string]interface{}
			obj, err = a.mongo.GetOneByIDRBAC(session, false, backd.PermissionUpdate, applicationID, relation.Source, relation.SourceID)
			if err == nil {
				objects = append(objects, obj)
			}
		}
	case "out": // outcoming relations collection:/id -> relation_name -> */*
		// if requester can not read the item stop here, it's unauthorized
		if a.mongo.Can(session, false, applicationID, ps.ByName("collection"), ps.ByName("id"), backd.PermissionRead) == false {
			rest.Unauthorized(w, r)
			return
		}

		query = map[string]interface{}{
			"src": ps.ByName("collection"),
			"sid": ps.ByName("id"),
			"rel": ps.ByName("relation"),
		}

		err = a.mongo.GetAll(applicationID, constants.ColRelations, query, []string{}, &relations)
		if err != nil {
			rest.ResponseErr(w, err)
			return
		}

		// add only those relations the user can manage
		for _, relation := range relations {
			var obj map[string]interface{}
			obj, err = a.mongo.GetOneByIDRBAC(session, false, backd.PermissionUpdate, applicationID, relation.Destination, relation.DestinationID)
			if err == nil {
				objects = append(objects, obj)
			}
		}
	default:
		rest.BadRequest(w, r, "wrong direction parameter")
		return
	}

	rest.Response(w, objects, nil, http.StatusOK, "")

}

// GET /relations/:id
func (a *apiStruct) getRelationID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var (
		relation      structs.Relation
		session       *pbsessions.Session
		applicationID string
		err           error
	)

	session, applicationID, err = a.getSession(r)
	if err != nil {
		rest.ResponseErr(w, err)
		return
	}

	err = rest.GetFromBody(r, &relation)
	if err != nil {
		rest.BadRequest(w, r, constants.ReasonBadQuery)
		return
	}

	err = a.mongo.GetOneByID(applicationID, constants.ColRelations, ps.ByName("id"), &relation)
	if err != nil {
		rest.NotFound(w, r)
		return
	}

	// in order to make a relation the requester must be able to write the item.
	// adding a relation can be considered like update the item data itself.
	if a.mongo.Can(session, false, applicationID, relation.Source, relation.SourceID, backd.PermissionUpdate) == false {
		rest.Unauthorized(w, r)
		return
	}

	// also... to be able to make a relation the requester must be able to read the destination item.
	if a.mongo.Can(session, false, applicationID, relation.Destination, relation.DestinationID, backd.PermissionRead) == false {
		rest.Unauthorized(w, r)
		return
	}

	rest.Response(w, relation, err, http.StatusCreated, rest.Location("relations", relation.ID))

}

// GET /relations/:collection/:id/:direction
// This endpoint returns the 'relation' object, not the object itself
func (a *apiStruct) getRelations(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var (
		relations         []structs.Relation // database relations
		relationsToReturn []structs.Relation // relations the requester can access to
		query             map[string]interface{}
		session           *pbsessions.Session
		applicationID     string
		err               error
	)

	session, applicationID, err = a.getSession(r)
	if err != nil {
		rest.ResponseErr(w, err)
		return
	}

	switch ps.ByName("direction") {
	case "in": // incoming relations */* -> relation_name -> collection:/id
		// if requester can not read the item stop here, it's unauthorized
		if a.mongo.Can(session, false, applicationID, ps.ByName("collection"), ps.ByName("id"), backd.PermissionRead) == false {
			rest.Unauthorized(w, r)
			return
		}

		query = map[string]interface{}{
			"dst": ps.ByName("collection"),
			"did": ps.ByName("id"),
		}

		err = a.mongo.GetAll(applicationID, constants.ColRelations, query, []string{}, &relations)
		if err != nil {
			rest.ResponseErr(w, err)
			return
		}

		// add only those relations the user can manage
		for _, relation := range relations {
			if a.mongo.Can(session, false, applicationID, relation.Source, relation.SourceID, backd.PermissionUpdate) {
				relationsToReturn = append(relationsToReturn, relation)
			}
		}
	case "out": // outcoming relations collection:/id -> relation_name -> */*
		// if requester can not read the item stop here, it's unauthorized
		if a.mongo.Can(session, false, applicationID, ps.ByName("collection"), ps.ByName("id"), backd.PermissionRead) == false {
			rest.Unauthorized(w, r)
			return
		}

		query = map[string]interface{}{
			"src": ps.ByName("collection"),
			"sid": ps.ByName("id"),
		}

		err = a.mongo.GetAll(applicationID, constants.ColRelations, query, []string{}, &relations)
		if err != nil {
			rest.ResponseErr(w, err)
			return
		}

		// add only those relations the user can manage
		for _, relation := range relations {
			if a.mongo.Can(session, false, applicationID, relation.Destination, relation.DestinationID, backd.PermissionRead) {
				relationsToReturn = append(relationsToReturn, relation)
			}
		}
	default:
		rest.BadRequest(w, r, "wrong direction parameter")
		return
	}

	rest.Response(w, relationsToReturn, err, http.StatusOK, "")

}

// POST /relations
func (a *apiStruct) postRelation(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var (
		relation      structs.Relation
		session       *pbsessions.Session
		applicationID string
		err           error
	)

	session, applicationID, err = a.getSession(r)
	if err != nil {
		rest.ResponseErr(w, err)
		return
	}

	err = rest.GetFromBody(r, &relation)
	if err != nil {
		rest.BadRequest(w, r, constants.ReasonBadQuery)
		return
	}

	// in order to make a relation the requester must be able to write the item.
	// adding a relation can be considered like update the item data itself.
	if a.mongo.Can(session, false, applicationID, relation.Source, relation.SourceID, backd.PermissionUpdate) == false {
		rest.Unauthorized(w, r)
		return
	}

	// also... to be able to make a relation the requester must be able to read the destination item.
	if a.mongo.Can(session, false, applicationID, relation.Destination, relation.DestinationID, backd.PermissionRead) == false {
		rest.Unauthorized(w, r)
		return
	}

	relation.ID = db.NewXID().String()
	relation.SetCreate(session.GetDomainId(), session.GetUserId())

	err = a.mongo.Insert(applicationID, constants.ColRelations, &relation)
	if err != nil {
		rest.ResponseErr(w, err)
		return
	}

	rest.Response(w, relation, err, http.StatusCreated, rest.Location("relations", relation.ID))

}

// DELETE /relations/:id
func (a *apiStruct) deleteRelationID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var (
		session       *pbsessions.Session
		applicationID string
		relation      structs.Relation
		err           error
	)

	session, applicationID, err = a.getSession(r)
	if err != nil {
		rest.Unauthorized(w, r)
		return
	}

	err = a.mongo.GetOneByID(applicationID, constants.ColRelations, ps.ByName("id"), &relation)
	if err != nil {
		rest.NotFound(w, r)
		return
	}

	// in order to delete a relation the requester must be able to write the item.
	if a.mongo.Can(session, false, applicationID, relation.Source, relation.SourceID, backd.PermissionUpdate) == false {
		rest.Unauthorized(w, r)
		return
	}

	err = a.mongo.DeleteByID(applicationID, constants.ColRelations, ps.ByName("id"))
	rest.Response(w, nil, err, http.StatusNoContent, "")

}
