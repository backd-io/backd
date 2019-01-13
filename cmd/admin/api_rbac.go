package main

import (
	"net/http"

	"github.com/backd-io/backd/backd"
	"github.com/backd-io/backd/internal/constants"
	"github.com/backd-io/backd/internal/pbsessions"
	"github.com/backd-io/backd/internal/rest"
	"github.com/backd-io/backd/internal/structs"
	"github.com/globalsign/mgo/bson"
	"github.com/julienschmidt/httprouter"
)

// POST /applications/:id/rbac
func (a *apiStruct) rbacApplications(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	a.rbac(w, r, ps, false)
}

// POST /domains/:id/rbac
func (a *apiStruct) rbacDomains(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	a.rbac(w, r, ps, true)
}

// NOTE: These endpoints do not act as natural REST, since work more like a RPC call. Access to the record on the database won't be done from the API.
//       As return codes:
//         - Bad Request, self explained
//         - OK, if the action sets the state requested (also if no action was necessary)
//         - Unauthorized
// Actions allowed:
//   - add.     Adds a new set of permissions(1-n) to the entity. Again, try to don't set permissions on users.
//   - remove.  Removes permissions(1-n) to the entity on the resource
//   - set.     Set will set only the new permissions requested overwritting the old ones.
//                To remove all permissions just set an empty array.
//   - get.     Returns all permissions the entity has granted for the resource. (not need to be explicit, it can be inherit by *)
//
func (a *apiStruct) rbac(w http.ResponseWriter, r *http.Request, ps httprouter.Params, isDomain bool) {

	var (
		rbac    backd.RBAC
		session *pbsessions.Session
		err     error
	)

	// getSession & rbac
	session, err = a.getSession(r)
	if err != nil {
		rest.ResponseErr(w, err)
		return
	}

	err = rest.GetFromBody(r, &rbac)
	if err != nil {
		rest.BadRequest(w, r, constants.ReasonReadingBody)
		return
	}

	if isDomain {
		rbac.Collection = constants.ColDomains
	}

	if rbac.Action == "" || rbac.Collection == "" || rbac.CollectionID == "" || rbac.DomainID == "" || rbac.IdentityID == "" {
		rest.BadRequest(w, r, constants.ReasonReadingBody)
		return
	}

	// ensure the current user can administer the resource
	if !a.mongo.Can(session, isDomain, ps.ByName("id"), rbac.Collection, rbac.CollectionID, backd.PermissionAdmin) {
		rest.Unauthorized(w, r)
		return
	}

	switch rbac.Action {
	case "add":
		for _, perm := range rbac.Permissions {
			query := map[string]interface{}{
				"did": rbac.DomainID,
				"iid": rbac.IdentityID,
				"c":   rbac.Collection,
				"cid": rbac.CollectionID,
				"p":   perm,
			}
			var count int
			count, err = a.mongo.Count(ps.ByName("id"), constants.ColRBAC, query)
			if count == 0 {
				// insert
				err = a.mongo.Insert(ps.ByName("id"), constants.ColRBAC, query)
				if err != nil {
					rest.ResponseErr(w, err)
					return
				}
			}
		}
		rest.Response(w, nil, err, http.StatusOK, "")
	case "remove":
		for _, perm := range rbac.Permissions {
			query := map[string]interface{}{
				"did": rbac.DomainID,
				"iid": rbac.IdentityID,
				"c":   rbac.Collection,
				"cid": rbac.CollectionID,
				"p":   perm,
			}
			var count int
			count, err = a.mongo.Count(ps.ByName("id"), constants.ColRBAC, query)
			if count == 1 {
				// insert
				err = a.mongo.Delete(ps.ByName("id"), constants.ColRBAC, query)
				if err != nil {
					rest.ResponseErr(w, err)
					return
				}
			}
		}
		rest.Response(w, nil, err, http.StatusOK, "")

	case "set":
		query := map[string]interface{}{
			"did": rbac.DomainID,
			"iid": rbac.IdentityID,
			"c":   rbac.Collection,
			"cid": rbac.CollectionID,
		}
		// remove all
		err = a.mongo.Delete(ps.ByName("id"), constants.ColRBAC, query)
		if err != nil {
			rest.ResponseErr(w, err)
			return
		}
		for _, perm := range rbac.Permissions {
			query := map[string]interface{}{
				"did": rbac.DomainID,
				"iid": rbac.IdentityID,
				"c":   rbac.Collection,
				"cid": rbac.CollectionID,
				"p":   perm,
			}
			err = a.mongo.Insert(ps.ByName("id"), constants.ColRBAC, query)
			if err != nil {
				rest.ResponseErr(w, err)
				return
			}
		}
		rest.Response(w, nil, err, http.StatusOK, "")
	case "get":

		// if the user can make the action on * can make it for the id
		var (
			cid         []string       // array of ids
			c           []string       // array of collections
			permissions []structs.RBAC // response from database
		)

		// get all collections entity are able to make anything
		c = append(c, "*")
		if rbac.Collection != "*" {
			c = append(c, rbac.Collection)
		}

		// get all ids user are able to make anything
		cid = append(cid, "*")
		if rbac.CollectionID != "*" {
			cid = append(cid, rbac.CollectionID)
		}

		query := map[string]interface{}{
			"did": rbac.DomainID,
			"iid": rbac.IdentityID,
			"c": bson.M{
				"$in": c,
			},
			"cid": bson.M{
				"$in": cid,
			},
		}

		err = a.mongo.GetAll(ps.ByName("id"), constants.ColRBAC, query, []string{}, &permissions)
		if err != nil {
			rest.ResponseErr(w, err)
			return
		}

		// now build the response with all the permissions the user has
		for _, perm := range permissions {
			if !in(perm.Permission, rbac.Permissions) {
				rbac.Permissions = append(rbac.Permissions, perm.Permission)
			}
		}
		// return rbac request filled
		rest.Response(w, rbac, err, http.StatusOK, "")
	}

}

func in(item string, items []string) bool {
	for _, i := range items {
		if i == item {
			return true
		}
	}
	return false
}
