package backd

import (
	"net/http"
	"strings"
)

// APIError is the struct that is returned when an error is returned from the APIs
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Reason  string `json:"reason,omitempty"`
	err     error
}

func (a APIError) Error() string {

	// check if http client returned an error
	if a.err != nil {
		return a.err.Error()
	}

	var texts []string

	if a.Message != "" {
		texts = append(texts, a.Message)
	}

	if a.Reason != "" {
		texts = append(texts, a.Reason)
	}

	return strings.Join(texts, " - ")
}

func (a APIError) wrapErr(err error, response *http.Response, expect int) error {

	// returns a known error to allow control it
	if response.StatusCode != expect {
		switch response.StatusCode {
		case 400:
			return ErrBadRequest
		case 401:
			return ErrUnauthorized
		case 404:
			return ErrNotFound
		case 405:
			return ErrMethodNotAllowed
		case 409:
			return ErrConflict
		}
	}

	if err == nil && a.Code == 0 {
		return nil
	}

	a.err = err
	return a

}

// Login is the struct that is expected by the API as request for an user authentication
type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Domain   string `json:"domain"`
}

// LoginResponse is the response if success. Upon a successful login it returns a
//   Session ID and ExpiresAt expiration date (as seconds from epoch)
type LoginResponse struct {
	ID        string `json:"id"`
	ExpiresAt int64  `json:"expires_at"`
}

// BootstrapRequest is the request to initialize a `backd` cluster
type BootstrapRequest struct {
	Code     string `json:"code"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RequestByID is a request that only especifies an ID (for simple things like group membership)
type RequestByID struct {
	ID string `json:"_id"`
}

// Metadata is the struct that defines how metadata is stored on the API
type Metadata struct {
	CreatedBy string `json:"cby" mapstructure:"cby"`
	UpdatedBy string `json:"uby" mapstructure:"uby"`
	CreatedAt int64  `json:"cat" mapstructure:"cat"`
	UpdatedAt int64  `json:"uat" mapstructure:"uat"`
}

// User is the struct that API expects to get on user operations
type User struct {
	ID                string                 `json:"_id" mapstructure:"_id"`                                 // (required) ID generated by the API
	Username          string                 `json:"username" mapstructure:"username"`                       // (required) Username is the entity that will be used for logon. If email will be used as username then both must match
	Name              string                 `json:"name" mapstructure:"name"`                               // (required) Name of the user (it can get filled with the data from the remote authorization authority)
	Email             string                 `json:"email" mapstructure:"email"`                             // (required) Email of the user (the one used to notify by mail)
	Description       string                 `json:"desc,omitempty" mapstructure:"desc,omitempty"`           // (optional) User description
	Password          string                 `json:"password,omitempty" mapstructure:"-"`                    // (optional) Password is only used to get the initial password on user creation
	GeneratedPassword string                 `json:"generated_password,omitempty" mapstructure:"-"`          // GeneratedPassword will be filled only if the user didn't set a password on user creation, so it generates one randomly
	Active            bool                   `json:"active,omitempty" mapstructure:"active,omitempty"`       // (required) Active defines when the user can interact with the APIs (some authorizations can leave it as active if the authentication system will allow or restrict the user)
	Validated         bool                   `json:"validated,omitempty" mapstructure:"validated,omitempty"` // (required) Validated shows if the user needs to make any action to active its email (and probably its account too)
	Data              map[string]interface{} `json:"data,omitempty" mapstructure:"data,omitempty"`           // (optional) Data is the arbitrary information that can be stored for the user
	Metadata          `json:"meta" mapstructure:"meta"`
}

// Group is the struct that api expects
type Group struct {
	ID          string `json:"_id"`            // (required) ID generated by the API
	Name        string `json:"name,omitempty"` // (required) Name of the group
	Description string `json:"desc,omitempty"` // (optional) Description
	Metadata    `json:"meta"`
}

// DomainType defines the behavior to build a session on a backd defined domain
type DomainType string

const (
	// DomainTypeBackd when set the domain will use natively only the backd users/groups
	DomainTypeBackd DomainType = "b"
	// DomainTypeActiveDirectory when set the domain will inherit the groups from the users
	//   on logon. So user membership will be updated from the ones received when the user
	//   creates a session.
	DomainTypeActiveDirectory DomainType = "ad"
)

// Domain struct
type Domain struct {
	ID          string                 `json:"_id"`
	Description string                 `json:"desc"`
	Type        DomainType             `json:"type"`
	Config      map[string]interface{} `json:"config,omitempty"`
	Metadata    `json:"meta"`
}

// QueryOptions is the builder of query parameters used for getMany queries
type QueryOptions struct {
	Q       map[string]interface{} `json:"q,omitempty" url:"-"` // for url it must be decode to string
	Sort    []string               `json:"sort,omitempty" url:"sort,omitempty"`
	Page    int                    `json:"page,omitempty" url:"page,omitempty"`
	PerPage int                    `json:"per_page,omitempty" url:"per_page,omitempty"`
}

// RBAC is the struct used to manage roles and permissions by the API
type RBAC struct {
	Action       string   `json:"action,omitempty"` // allowed actions: add / remove / set
	DomainID     string   `json:"domain_id"`        // domain
	IdentityID   string   `json:"identity_id"`      // user_id / group_id
	Collection   string   `json:"collection"`       // collection if application, if domain there is no concept of 'collection' you can manage or not if entity_id match
	CollectionID string   `json:"collection_id"`    // id if application, entity_id if domain
	Permissions  []string `json:"permissions"`      // array of permissions matching entity and item
}

// Relation is the representation of linked data.
type Relation struct {
	ID            string `json:"_id"`
	Source        string `json:"src"`
	SourceID      string `json:"sid"`
	Destination   string `json:"dst"`
	DestinationID string `json:"did"`
	Relation      string `json:"rel"`
	Metadata      `json:"meta"`
}
