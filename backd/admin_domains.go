package backd

// Domains is a simple struct that holds the operations doable for domains
type Domains struct {
	backd *Backd
}

// Domains returns an instance of the Domains struct
func (b *Backd) Domains() *Domains {
	return &Domains{
		backd: b,
	}
}

// Admin is the struct that contains all the actions doable with the Objects API
type Admin struct {
	backd    *Backd
	domainID string
	Users    *AdminUsers
	Groups   *AdminGroups
}

// Domain returns an instance of the Admin struct
func (b *Backd) Domain(domainID string) *Admin {
	return &Admin{
		backd:    b,
		domainID: domainID,
		Users:    b.newAdminUsers(domainID),
		Groups:   b.newAdminGroups(domainID),
	}
}

// AdminUsers holds users operations
type AdminUsers struct {
	backd    *Backd
	domainID string
}

// AdminGroups holds groups operations
type AdminGroups struct {
	backd    *Backd
	domainID string
}

// AdminRBAC holds groups operations
type AdminRBAC struct {
	backd    *Backd
	domainID string
}

func (b *Backd) newAdminUsers(domainID string) *AdminUsers {
	return &AdminUsers{
		backd:    b,
		domainID: domainID,
	}
}

func (b *Backd) newAdminGroups(domainID string) *AdminGroups {
	return &AdminGroups{
		backd:    b,
		domainID: domainID,
	}
}

func (b *Backd) newAdminRBAC(domainID string) *AdminRBAC {
	return &AdminRBAC{
		backd:    b,
		domainID: domainID,
	}
}

// domains

// GetMany returns all domains that matches the conditions especified
func (a *Domains) GetMany(queryOptions QueryOptions, object interface{}) error {
	return a.backd.get(adminMS, []string{"domains"}, queryOptions, object, a.backd.headers())
}

// GetByID returns an domain by its ID
func (a *Domains) GetByID(id string, object interface{}) error {
	return a.backd.getByID(adminMS, []string{"domains", id}, object, a.backd.headers())
}

// Insert inserts a new domain on the desired collection if the user have the required permissions
func (a *Domains) Insert(object interface{}) (map[string]interface{}, error) {
	return a.backd.insert(adminMS, []string{"domains"}, object, a.backd.headers())
}

// Update updates the required domain if the user has permissions for
//   from is the original domain updated by the user
//   to   is the updated domain retreived by the API
func (a *Domains) Update(id string, from, to interface{}) error {
	return a.backd.update(adminMS, []string{"domains", id}, from, to, a.backd.headers())
}

// Delete removes a domain by ID
func (a *Domains) Delete(id string) error {
	return a.backd.delete(adminMS, []string{"domains", id}, a.backd.headers())
}

// users

// GetMany returns all users that matches the conditions especified
func (a *AdminUsers) GetMany(queryOptions QueryOptions, object interface{}) error {
	return a.backd.get(adminMS, []string{"domains", a.domainID, "users"}, queryOptions, object, a.backd.headers())
}

// GetByID returns an user by its ID
func (a *AdminUsers) GetByID(id string, object interface{}) error {
	return a.backd.getByID(adminMS, []string{"domains", a.domainID, "users", id}, object, a.backd.headers())
}

// Insert inserts a new user on the desired collection if the user have the required permissions
func (a *AdminUsers) Insert(object interface{}) (map[string]interface{}, error) {
	return a.backd.insert(adminMS, []string{"domains", a.domainID, "users"}, object, a.backd.headers())
}

// Update updates the required user if the user has permissions for
//   from is the original user updated by the user
//   to   is the object retreived by the API
func (a *AdminUsers) Update(id string, from, to interface{}) error {
	return a.backd.update(adminMS, []string{"domains", a.domainID, "users", id}, from, to, a.backd.headers())
}

// Delete removes a user by ID
func (a *AdminUsers) Delete(id string) error {
	return a.backd.delete(adminMS, []string{"domains", a.domainID, "users", id}, a.backd.headers())
}

// groups

// GetMany returns all groups that matches the conditions especified
func (a *AdminGroups) GetMany(queryOptions QueryOptions, object interface{}) error {
	return a.backd.get(adminMS, []string{"domains", a.domainID, "groups"}, queryOptions, object, a.backd.headers())
}

// GetByID returns an group by its ID
func (a *AdminGroups) GetByID(id string, object interface{}) error {
	return a.backd.getByID(adminMS, []string{"domains", a.domainID, "groups", id}, object, a.backd.headers())
}

// Insert inserts a new group on the desired collection if the user have the required permissions
func (a *AdminGroups) Insert(object interface{}) (map[string]interface{}, error) {
	return a.backd.insert(adminMS, []string{"domains", a.domainID, "groups"}, object, a.backd.headers())
}

// Update updates the required group if the user has permissions for
//   from is the original group updated by the user
//   to   is the group retreived by the API
func (a *AdminGroups) Update(id string, from, to interface{}) error {
	return a.backd.update(adminMS, []string{"domains", a.domainID, "groups", id}, from, to, a.backd.headers())
}

// Delete removes a group by ID
func (a *AdminGroups) Delete(id string) error {
	return a.backd.delete(adminMS, []string{"domains", a.domainID, "groups", id}, a.backd.headers())
}

// AddMember adds a new member to the group
func (a *AdminGroups) AddMember(id, userID string) error {
	return a.backd.update(adminMS, []string{"domains", a.domainID, "groups", id, "members", userID}, nil, nil, a.backd.headers())
}

// RemoveMember removes a member from a group by ID
func (a *AdminGroups) RemoveMember(id, userID string) error {
	return a.backd.delete(adminMS, []string{"domains", a.domainID, "groups", id, "members", userID}, a.backd.headers())
}

// Set sets a new role permission set
func (a *AdminRBAC) Set(rbac RBAC) error {
	rbac.Action = RBACActionSet
	return a.backd.insertRBAC(adminMS, []string{"domains", a.domainID, "rbac"}, rbac, a.backd.headers())
}

// Get get cirremt role permission set
func (a *AdminRBAC) Get(rbac RBAC) error {
	rbac.Action = RBACActionGet
	return a.backd.insertRBAC(adminMS, []string{"domains", a.domainID, "rbac"}, rbac, a.backd.headers())
}

// Add adds role/s to the role permission set
func (a *AdminRBAC) Add(rbac RBAC) error {
	rbac.Action = RBACActionAdd
	return a.backd.insertRBAC(adminMS, []string{"domains", a.domainID, "rbac"}, rbac, a.backd.headers())
}

// Remove removes role/s to the role permission set
func (a *AdminRBAC) Remove(rbac RBAC) error {
	rbac.Action = RBACActionRemove
	return a.backd.insertRBAC(adminMS, []string{"domains", a.domainID, "rbac"}, rbac, a.backd.headers())
}
