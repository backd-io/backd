package main

import (
	"fmt"
	"net/http"

	"github.com/backd-io/backd/internal/utils"

	"github.com/backd-io/backd/backd"
	"github.com/backd-io/backd/internal/constants"
	"github.com/backd-io/backd/internal/db"
	"github.com/backd-io/backd/internal/rest"
	"github.com/backd-io/backd/internal/structs"
	"github.com/julienschmidt/httprouter"
	"github.com/sethvargo/go-password/password"
	"go.uber.org/zap"
)

// if the service hasn't been configured for the first time `backd` won't make anything than
//   return 401, but we need to bootstrap the service by using the cli.
// for the bootstrap process this service will want to receive a struct to create the first
//   user that will be admin.
func (a *apiStruct) isBootstrapped() error {

	var (
		err error
	)

	err = a.mongo.IsInitialized(constants.DBBackdApp)

	if err == db.ErrDatabaseNotInitialized {
		a.bootstrapCode, err = password.Generate(32, 16, 0, true, true)
		a.inst.Info("server not bootstrapped", zap.String("code", a.bootstrapCode))
	}
	return err

}

func (a *apiStruct) isBootstrapRequestOK(bootstrapRequest *backd.BootstrapRequest) bool {

	if bootstrapRequest.Code != a.bootstrapCode {
		return false
	}

	if bootstrapRequest.Email == "" || bootstrapRequest.Name == "" || bootstrapRequest.Password == "" || bootstrapRequest.Username == "" {
		return false
	}

	return true

}

func (a *apiStruct) postBootstrap(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	// if the server is already bootstrapped don't allow to reconfigure
	if a.bootstrapCode == "" {
		rest.NotAllowed(w, r)
		return
	}

	var (
		bootstrapRequest backd.BootstrapRequest
		user             structs.User
		membership       structs.Membership
		group            structs.Group
		rbac             structs.RBAC
		rbac1            structs.RBAC
		err              error
	)

	// ensure we fill the struct with the required data
	err = rest.GetFromBody(r, &bootstrapRequest)
	if err != nil || a.isBootstrapRequestOK(&bootstrapRequest) != true {
		rest.BadRequest(w, r, "error getting data from body")
		return
	}

	// build db, user and set permissions
	err = a.mongo.CreateBackdDatabases()
	if err != nil {
		rest.BadRequest(w, r, "error creating backd databases")
		return
	}

	// Create Administrator User
	user.ID = db.NewXID().String()
	user.Name = bootstrapRequest.Name
	user.Username = bootstrapRequest.Username
	user.Email = bootstrapRequest.Email
	user.Active = true
	user.Validated = true
	user.SetPassword(bootstrapRequest.Password)
	user.SetCreate(constants.DBBackdDom, constants.SystemUserID)
	err = a.mongo.Insert(constants.DBBackdDom, constants.ColUsers, &user)
	if err != nil {
		fmt.Println("error creating admin user")
		utils.Prettify(user)
		rest.BadRequest(w, r, "error creating admin user")
		return
	}

	// Create Domain Administrators group
	group.ID = db.NewXID().String()
	group.Name = constants.GroupDomainAdministrators
	group.Description = "backd Domain Administrators"
	group.SetCreate(constants.DBBackdDom, constants.SystemUserID)
	err = a.mongo.Insert(constants.DBBackdDom, constants.ColGroups, &group)
	if err != nil {
		rest.BadRequest(w, r, "error creating admin group")
		return
	}

	// Add user to group
	membership.GroupID = group.ID
	membership.UserID = user.ID
	err = a.mongo.Insert(constants.DBBackdDom, constants.ColMembership, &membership)
	if err != nil {
		rest.BadRequest(w, r, "error adding user to admin group")
		return
	}

	// add admin permissions over `_backd` domain to Domain Administrators
	rbac.ID = db.NewXID().String()
	rbac.DomainID = constants.DBBackdDom
	rbac.IdentityID = group.ID
	rbac.Collection = constants.ColDomains
	rbac.CollectionID = "*"
	rbac.Permission = string(backd.PermissionAdmin)
	err = a.mongo.Insert(constants.DBBackdDom, constants.ColRBAC, &rbac)
	if err != nil {
		rest.BadRequest(w, r, "error adding user to admin group")
		return
	}
	// add admin permissions over `backd` application to Domain Administrators
	rbac1.ID = db.NewXID().String()
	rbac1.DomainID = constants.DBBackdDom
	rbac1.IdentityID = group.ID
	rbac1.Collection = "*"
	rbac1.CollectionID = "*"
	rbac1.Permission = string(backd.PermissionAdmin)
	err = a.mongo.Insert(constants.DBBackdApp, constants.ColRBAC, &rbac1)

	// do not allow this operation again
	a.bootstrapCode = ""

	a.inst.Info("server bootstrapped correctly")
	rest.Response(w, nil, err, http.StatusCreated, "")

}
