package backd

import (
	"net/http"
	"time"
)

// Login sends a log in request to the API
func (b *Backd) Login(username, password, domain string) (bool, error) {

	var (
		body     Login
		success  LoginResponse
		failure  APIError
		response *http.Response
		err      error
	)

	body = Login{
		Username: username,
		Password: password,
		Domain:   domain,
	}

	response, err = b.sling.Post(b.buildPath(authMS, []string{pathSession})).BodyJSON(&body).Receive(&success, &failure)

	err = failure.wrapErr(err, response, http.StatusOK)

	if err != nil {
		return false, err
	}

	b.sessionID = success.ID
	b.expiresAt = success.ExpiresAt
	return true, nil

}

// Logout deletes the session on the API so the client will make request (if any) as anonymous
func (b *Backd) Logout() (bool, error) {

	var (
		failure  APIError
		response *http.Response
		err      error
	)

	response, err = b.sling.Set(HeaderSessionID, b.sessionID).Delete(b.buildPath(authMS, []string{pathSession})).Receive(nil, &failure)

	err = failure.wrapErr(err, response, http.StatusNoContent)

	if err != nil {
		return false, err
	}

	b.sessionID = ""
	b.expiresAt = 0
	return true, nil

}

// Session returns current session status and remaining time if session is established
func (b *Backd) Session() (string, int, time.Time) {

	var (
		expiresAt time.Time
	)

	if b.sessionID == "" {
		return b.sessionID, StateAnonymous, expiresAt
	}

	if time.Now().Unix() < b.expiresAt {
		return b.sessionID, StateLoggedIn, time.Unix(b.expiresAt, 0)
	}

	return b.sessionID, StateExpired, expiresAt

}

// SetSession sets a sessionID and expires information from elsewhere, used as commodity for the cli
//   No check will be done on the client library so errors (if any) will arise when requesting the API
func (b *Backd) SetSession(sessionID string, expiresAt int64) {
	b.sessionID = sessionID
	b.expiresAt = expiresAt
}
