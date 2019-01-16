package backd

// Apps is a simple struct that holds the operations doable for applications
type Apps struct {
	backd *Backd
}

// Apps returns an instance of the Apps struct
func (b *Backd) Apps() *Apps {
	return &Apps{
		backd: b,
	}
}

// AdminApplication is the struct that contains all the actions doable with the Objects API
type AdminApplication struct {
	backd *Backd
	appID string
	RBAC  *AdminAppRBAC
	// Users  *AdminUsers
	// Groups *AdminGroups
}

// App returns an instance of the Admin struct
func (b *Backd) App(appID string) *AdminApplication {
	return &AdminApplication{
		backd: b,
		appID: appID,
	}
}

// AdminAppRBAC holds groups operations
type AdminAppRBAC struct {
	backd *Backd
	appID string
}

// RBAC returns an instance of the AdminAppRBAC struct
func (b *Backd) RBAC(appID string) *AdminAppRBAC {
	return &AdminAppRBAC{
		backd: b,
		appID: appID,
	}
}

// apps

// GetMany returns all applications that matches the conditions especified
func (a *Apps) GetMany(queryOptions QueryOptions, object interface{}) error {
	return a.backd.get(adminMS, []string{"applications"}, queryOptions, object, a.backd.headers())
}

// GetByID returns an application by its ID
func (a *Apps) GetByID(id string, object interface{}) error {
	return a.backd.getByID(adminMS, []string{"applications", id}, object, a.backd.headers())
}

// Insert inserts a new application on the desired collection if the user have the required permissions
func (a *Apps) Insert(object interface{}) (map[string]interface{}, error) {
	return a.backd.insert(adminMS, []string{"applications"}, object, a.backd.headers())
}

// Update updates the required application if the user has permissions for
//   from is the original domain updated by the user
//   to   is the updated application retreived by the API
func (a *Apps) Update(id string, from, to interface{}) error {
	return a.backd.update(adminMS, []string{"applications", id}, from, to, a.backd.headers())
}

// Delete removes an application by ID
func (a *Apps) Delete(id string) error {
	return a.backd.delete(adminMS, []string{"applications", id}, a.backd.headers())
}

// Set sets a new role permission set
func (a *AdminAppRBAC) Set(rbac RBAC) error {
	rbac.Action = RBACActionSet
	return a.backd.insertRBAC(adminMS, []string{"applications", a.appID, "rbac"}, rbac, a.backd.headers())
}

// Get get current role permission set
func (a *AdminAppRBAC) Get(rbac *RBAC) error {
	rbac.Action = RBACActionGet
	return a.backd.insertRBAC(adminMS, []string{"applications", a.appID, "rbac"}, rbac, a.backd.headers())
}

// Add adds role/s to the role permission set
func (a *AdminAppRBAC) Add(rbac RBAC) error {
	rbac.Action = RBACActionAdd
	return a.backd.insertRBAC(adminMS, []string{"applications", a.appID, "rbac"}, rbac, a.backd.headers())
}

// Remove removes role/s to the role permission set
func (a *AdminAppRBAC) Remove(rbac RBAC) error {
	rbac.Action = RBACActionRemove
	return a.backd.insertRBAC(adminMS, []string{"applications", a.appID, "rbac"}, rbac, a.backd.headers())
}
