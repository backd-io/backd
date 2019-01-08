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
	// Users  *AdminUsers
	// Groups *AdminGroups
}

// App returns an instance of the Admin struct
func (b *Backd) App(appID string) *AdminApplication {
	return &AdminApplication{
		backd: b,
		appID: appID,
		// Users:    b.newAdminUsers(domainID),
		// Groups:   b.newAdminGroups(domainID),
	}
}

// apps

// GetMany returns all domains that matches the conditions especified
func (a *Apps) GetMany(queryOptions QueryOptions, object interface{}) error {
	return a.backd.get(adminMS, []string{"applications"}, queryOptions, object, a.backd.headers())
}

// GetByID returns an domain by its ID
func (a *Apps) GetByID(id string, object interface{}) error {
	return a.backd.getByID(adminMS, []string{"applications", id}, object, a.backd.headers())
}

// Insert inserts a new domain on the desired collection if the user have the required permissions
func (a *Apps) Insert(object interface{}) (id string, err error) {
	return a.backd.insert(adminMS, []string{"applications"}, object, a.backd.headers())
}

// Update updates the required domain if the user has permissions for
//   from is the original domain updated by the user
//   to   is the updated domain retreived by the API
func (a *Apps) Update(id string, from, to interface{}) error {
	return a.backd.update(adminMS, []string{"applications", id}, from, to, a.backd.headers())
}

// Delete removes a domain by ID
func (a *Apps) Delete(id string) error {
	return a.backd.delete(adminMS, []string{"applications", id}, a.backd.headers())
}
