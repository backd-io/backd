package backd

import (
	"net/http"
	"time"
)

// Login sends a log in request to the api
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

	response, err = b.sling.Post(pathSession).BodyJSON(&body).Receive(&success, &failure)

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

	response, err = b.sling.Set(HeaderSessionID, b.sessionID).Delete(pathSession).Receive(nil, &failure)

	err = failure.wrapErr(err, response, http.StatusNoContent)

	if err != nil {
		return false, err
	}

	b.sessionID = ""
	b.expiresAt = 0
	return true, nil

}

// SessionState returns current session status and remaining time if session is established
func (b *Backd) SessionState() (int, time.Duration) {

	var (
		expiresIn time.Duration
	)

	if b.sessionID == "" {
		return StateAnonymous, expiresIn
	}

	if time.Now().Unix() > b.expiresAt {
		return StateLoggedIn, time.Since(time.Unix(b.expiresAt, 0))
	}

	return StateExpired, 0
}
