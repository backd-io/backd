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
