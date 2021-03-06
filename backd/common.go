package backd

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"path"

	"github.com/google/go-querystring/query"
)

// auth paths
const (
	pathSession string = "session"
)

// admin paths
const (
	pathBootstrap = "bootstrap"
)

// local client constants
const (
	clientName = "Backd Go Client"
)

// Permission is the required level of permission required to operate
type Permission string

// Headers for the endpoints
const (
	HeaderSessionID     = "X-Session-ID"
	HeaderApplicationID = "X-Application-ID"
)

// Exported permissions
const (
	PermissionRead   Permission = "read"
	PermissionCreate Permission = "create"
	PermissionUpdate Permission = "update"
	PermissionDelete Permission = "delete"
	PermissionAdmin  Permission = "admin"
)

// Session state
const (
	StateAnonymous int = iota
	StateExpired
	StateLoggedIn
)

// RBAC Actions
const (
	RBACActionGet    = "get"
	RBACActionSet    = "set"
	RBACActionAdd    = "add"
	RBACActionRemove = "remove"
)

// HTTP Errors returns if found
var (
	ErrBadRequest       = errors.New(http.StatusText(http.StatusBadRequest))
	ErrUnauthorized     = errors.New(http.StatusText(http.StatusUnauthorized))
	ErrNotFound         = errors.New(http.StatusText(http.StatusNotFound))
	ErrMethodNotAllowed = errors.New(http.StatusText(http.StatusMethodNotAllowed))
	ErrConflict         = errors.New(http.StatusText(http.StatusConflict))
)

type microservice int

const (
	adminMS microservice = iota
	authMS
	objectsMS
	functionsMS
)

func (b *Backd) buildPath(m microservice, parts []string) string {

	var (
		urlString string
	)

	switch m {
	case adminMS:
		urlString = b.adminURL
	case authMS:
		urlString = b.authURL
	case objectsMS:
		urlString = b.objectsURL
	case functionsMS:
		urlString = b.functionsURL
	}

	u, err := url.Parse(urlString)
	if err != nil {
		panic(err)
	}

	for _, part := range parts {
		u.Path = path.Join(u.Path, part)
	}

	return u.String()

}

func (b *Backd) buildPathWithQueryOptions(m microservice, parts []string, options QueryOptions) string {

	var (
		urlString   string
		queryString string
		q           string
		values      url.Values
		err         error
	)

	urlString = b.buildPath(m, parts)

	values, err = query.Values(options)
	if err == nil {
		queryString = values.Encode()
	}

	if options.Q != nil {
		// get the string
		b, _ := json.Marshal(options.Q)
		if len(b) > 0 {
			if len(queryString) > 0 {
				q = fmt.Sprintf("&q=%s", string(b))
			} else {
				q = fmt.Sprintf("q=%s", string(b))
			}
		}
	}

	// concatenate
	queryString = queryString + q

	if urlString != "" && queryString != "" {
		return fmt.Sprintf("%s?%s", urlString, queryString)
	}
	return urlString

}
