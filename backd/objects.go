package backd

import (
	"net/http"
)

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

// GetByID returns an object by it's id
func (o *Objects) GetByID(collection, id string, object interface{}) error {

	var (
		failure  APIError
		response *http.Response
		err      error
	)

	response, err = o.backd.sling.
		Set(HeaderSessionID, o.backd.sessionID).
		Set(HeaderApplicationID, o.applicationID).
		Get(o.backd.buildPath(objects, collection, id)).
		Receive(object, &failure)

	// rebuild and return err
	return failure.wrapErr(err, response, http.StatusOK)

}

// Insert inserts a new object on the desired collection if the user have the required permissions
func (o *Objects) Insert(collection string, object interface{}) (id string, err error) {

	var (
		success  map[string]interface{}
		failure  APIError
		response *http.Response
	)

	response, err = o.backd.sling.
		Set(HeaderSessionID, o.backd.sessionID).
		Set(HeaderApplicationID, o.applicationID).
		Post(o.backd.buildPath(objects, collection)).
		BodyJSON(object).
		Receive(&success, &failure)

	// rebuild err
	err = failure.wrapErr(err, response, http.StatusOK)

	if err == nil {
		id, _ = success["_id"].(string)
	}

	return
}

// Update updates the required object if the user has permissions for
//   from is the original object updated by the user
//   to   is the object retreived by the API
func (o *Objects) Update(collection, id string, from, to interface{}) error {

	var (
		failure  APIError
		response *http.Response
		err      error
	)

	response, err = o.backd.sling.
		Set(HeaderSessionID, o.backd.sessionID).
		Set(HeaderApplicationID, o.applicationID).
		Put(o.backd.buildPath(objects, collection, id)).
		BodyJSON(from).
		Receive(to, &failure)

	// rebuild and return err
	return failure.wrapErr(err, response, http.StatusOK)

}

// Delete removes a object by ID
func (o *Objects) Delete(collection, id string) error {

	var (
		failure  APIError
		response *http.Response
		err      error
	)

	response, err = o.backd.sling.
		Set(HeaderSessionID, o.backd.sessionID).
		Set(HeaderApplicationID, o.applicationID).
		Delete(o.backd.buildPath(objects, collection, id)).
		Receive(nil, &failure)

	return failure.wrapErr(err, response, http.StatusNoContent)

}
