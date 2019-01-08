package backd

// Objects is the struct that contains all the actions doable with the Objects API
type Objects struct {
	backd         *Backd
	applicationID string
}

// Objects returns an instance of the Objects struct
func (b *Backd) Objects(applicationID string) *Objects {
	return &Objects{
		backd:         b,
		applicationID: applicationID,
	}
}

func (o *Objects) headers() map[string]string {
	return map[string]string{
		HeaderSessionID:     o.backd.sessionID,
		HeaderApplicationID: o.applicationID,
	}
}

// GetMany returns all objects that matches the conditions especified
func (o *Objects) GetMany(collection string, queryOptions QueryOptions, object interface{}) error {
	return o.backd.get(objectsMS, []string{collection}, queryOptions, object, o.headers())
}

// GetByID returns an object by its ID
func (o *Objects) GetByID(collection, id string, object interface{}) error {
	return o.backd.getByID(objectsMS, []string{collection, id}, object, o.headers())
}

// Insert inserts a new object on the desired collection if the user have the required permissions
func (o *Objects) Insert(collection string, object interface{}) (id string, err error) {
	return o.backd.insert(objectsMS, []string{collection}, object, o.headers())
}

// Update updates the required object if the user has permissions for
//   from is the original object updated by the user
//   to   is the object retreived by the API
func (o *Objects) Update(collection, id string, from, to interface{}) error {
	return o.backd.update(objectsMS, []string{collection, id}, from, to, o.headers())
}

// Delete removes a object by ID
func (o *Objects) Delete(collection, id string) error {
	return o.backd.delete(objectsMS, []string{collection, id}, o.headers())
}
