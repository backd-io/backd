package backd

// Funcs is the struct that contains all helpers to work with the
//    functions on the admin endpoint
type Funcs struct {
	backd         *Backd
	applicationID string
}

// Functions returns an instance of Funcs exposing the helper functions
//    for the function publishing workflow
func (b *Backd) Functions(applicationID string) *Funcs {
	return &Funcs{
		backd:         b,
		applicationID: applicationID,
	}
}

func (a *Funcs) headers() map[string]string {
	return map[string]string{
		HeaderSessionID:     a.backd.sessionID,
		HeaderApplicationID: a.applicationID,
	}
}

// GetMany returns all functions that matches the conditions especified
func (a *Funcs) GetMany(queryOptions QueryOptions, object interface{}) error {
	return a.backd.get(adminMS, []string{"applications", a.applicationID, "functions"}, queryOptions, object, a.headers())
}

// GetByID returns an function by its ID
func (a *Funcs) GetByID(id string, object interface{}) error {
	return a.backd.getByID(adminMS, []string{"applications", a.applicationID, "functions", id}, object, a.headers())
}

// Insert inserts a new function if the user have the required permissions
func (a *Funcs) Insert(object interface{}) (map[string]interface{}, error) {
	return a.backd.insert(adminMS, []string{"applications", a.applicationID, "functions"}, object, a.headers())
}

// Update updates the required function if the user has permissions for
//   from is the original domain updated by the user
//   to   is the updated domain retreived by the API
func (a *Funcs) Update(id string, from, to interface{}) error {
	return a.backd.update(adminMS, []string{"applications", a.applicationID, "functions", id}, from, to, a.headers())
}

// Delete removes a function by ID
func (a *Funcs) Delete(id string) error {
	return a.backd.delete(adminMS, []string{"applications", a.applicationID, "functions", id}, a.headers())
}

// Run executes a function by its ID(name), input expected is a map[string]interface{}
//   outputs a map[string]interface{} from the function itself
func (a *Funcs) Run(id string, input map[string]interface{}) (map[string]interface{}, error) {
	return a.backd.insert(functionsMS, []string{"functions", id}, input, a.headers())
}
