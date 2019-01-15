package backd

import (
	"bytes"
	"net"
	"net/http"
	"time"

	"github.com/dghubble/sling"
)

// Backd is the struct that holds the client for the service
type Backd struct {
	sling      *sling.Sling
	authURL    string
	objectsURL string
	adminURL   string
	sessionID  string
	expiresAt  int64
}

// NewClient returns an usable client to connect to an instance of Backd
func NewClient(authURL, objectsURL, adminURL string) *Backd {

	var (
		backd Backd
	)

	backd.authURL = authURL
	backd.objectsURL = objectsURL
	backd.adminURL = adminURL
	backd.ConnectionTimeouts(5, 5, 10)

	return &backd

}

// ConnectionTimeouts allow to change the client timeouts for:
//  - Dialer
//  - TLS Handshake
//  - HTTP timeout
func (b *Backd) ConnectionTimeouts(dialer, tlsHandshake, timeout time.Duration) {

	b.sling = sling.New().Client(&http.Client{
		Timeout: timeout * time.Second,
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout: dialer * time.Second,
			}).Dial,
			TLSHandshakeTimeout: tlsHandshake * time.Second,
		},
	}).Set("User-Agent", clientName)

}

// BootstrapCluster creates the first user on the cluster with full Administration
// permissions on the the backd application
func (b *Backd) BootstrapCluster(code, name, username, email, password string) error {

	var (
		body     BootstrapRequest
		failure  APIError
		response *http.Response
		err      error
	)

	body = BootstrapRequest{
		Code:     code,
		Name:     name,
		Username: username,
		Email:    email,
		Password: password,
	}

	response, err = b.sling.Post(b.buildPath(adminMS, []string{pathBootstrap})).BodyJSON(&body).Receive(nil, &failure)

	return failure.wrapErr(err, response, http.StatusCreated)

}

// headers returns the common headers needed to operate (session ID)
func (b *Backd) headers() map[string]string {
	return map[string]string{
		HeaderSessionID: b.sessionID,
	}
}

// get is the generic getMany items from somewhere:
// - expects parts as part of the full URL
// - expects queryStrings map
func (b *Backd) get(m microservice, parts []string, queryOptions QueryOptions, data interface{}, headers map[string]string) error {

	var (
		failure  APIError
		response *http.Response
		sling    *sling.Sling
		err      error
	)

	sling = b.sling

	for key, value := range headers {
		sling.Set(key, value)
	}

	response, err = sling.
		Get(b.buildPathWithQueryOptions(m, parts, queryOptions)).
		Receive(data, &failure)

	// rebuild and return err
	return failure.wrapErr(err, response, http.StatusOK)

}

// getByID returns something by its id
func (b *Backd) getByID(m microservice, parts []string, object interface{}, headers map[string]string) error {

	var (
		failure  APIError
		response *http.Response
		sling    *sling.Sling
		err      error
	)

	sling = b.sling

	for key, value := range headers {
		sling.Set(key, value)
	}

	response, err = sling.
		Get(b.buildPath(m, parts)).
		Receive(object, &failure)

	// rebuild and return err
	return failure.wrapErr(err, response, http.StatusOK)

}

// insert allows to insert the required object on the API
func (b *Backd) insert(m microservice, parts []string, object interface{}, headers map[string]string) (success map[string]interface{}, err error) {

	var (
		failure  APIError
		response *http.Response
		sling    *sling.Sling
	)

	sling = b.sling

	for key, value := range headers {
		sling.Set(key, value)
	}

	if _, ok := object.([]byte); ok {
		response, err = sling.
			Post(b.buildPath(m, parts)).
			Body(bytes.NewReader(object.([]byte))).
			Receive(&success, &failure)
	} else {
		response, err = sling.
			Post(b.buildPath(m, parts)).
			BodyJSON(object).
			Receive(&success, &failure)
	}

	// rebuild err
	err = failure.wrapErr(err, response, http.StatusCreated)
	return
}

// insertRBAC allows to insert the required object on the API
func (b *Backd) insertRBAC(m microservice, parts []string, object interface{}, headers map[string]string) (err error) {

	var (
		success  map[string]interface{}
		failure  APIError
		response *http.Response
		sling    *sling.Sling
	)

	sling = b.sling

	for key, value := range headers {
		sling.Set(key, value)
	}

	response, err = sling.
		Post(b.buildPath(m, parts)).
		BodyJSON(object).
		Receive(&success, &failure)

	// rebuild err
	err = failure.wrapErr(err, response, http.StatusOK)
	return
}

// update updates the required object if the user has permissions for
//   from is the original object updated by the user
//   to   is the object retreived by the API
func (b *Backd) update(m microservice, parts []string, from, to interface{}, headers map[string]string) error {

	var (
		failure  APIError
		response *http.Response
		sling    *sling.Sling
		err      error
	)

	sling = b.sling

	for key, value := range headers {
		sling.Set(key, value)
	}

	if _, ok := from.([]byte); ok {
		response, err = sling.
			Put(b.buildPath(m, parts)).
			Body(bytes.NewReader(from.([]byte))).
			Receive(to, &failure)
	} else {
		response, err = sling.
			Put(b.buildPath(m, parts)).
			BodyJSON(from).
			Receive(to, &failure)
	}

	// rebuild and return err
	return failure.wrapErr(err, response, http.StatusOK)

}

// delete removes a object by ID
func (b *Backd) delete(m microservice, parts []string, headers map[string]string) error {

	var (
		failure  APIError
		response *http.Response
		sling    *sling.Sling
		err      error
	)

	sling = b.sling

	for key, value := range headers {
		sling.Set(key, value)
	}

	response, err = sling.
		Delete(b.buildPath(m, parts)).
		Receive(nil, &failure)

	return failure.wrapErr(err, response, http.StatusNoContent)

}
