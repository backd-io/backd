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
