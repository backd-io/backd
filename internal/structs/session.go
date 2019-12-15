package structs

import (
	"crypto/rand"
	"encoding/base64"
	"time"
)

// Session is the struct that reflects the information of the user
//   currently logged into the domain
type Session struct {
	ID        string   `json:"_id"`
	DomainID  string   `json:"did"`
	User      User     `json:"uid"`
	ExpiresAt int64    `json:"eat"`
	CreatedAt int64    `json:"cat"`
	Groups    []string `json:"g"`
}

// IsExpired returns expiration status of the session
func (s *Session) IsExpired() bool {
	return time.Now().Unix() > s.ExpiresAt
}

// NewToken creates a random token of 43 characters
func (s *Session) NewToken() (err error) {

	var (
		arr []byte
	)

	arr = make([]byte, 32)

	_, err = rand.Read(arr)
	if err != nil {
		return
	}

	s.ID = base64.RawURLEncoding.EncodeToString(arr)

	return

}

// SessionResponse is the struct that will be returned to the client
//   when a session has been established
type SessionResponse struct {
	ID        string `json:"id"`
	ExpiresAt int64  `json:"expires_at"`
}

// SessionRequest is the struct that defines how an user creates a session
//   on the auth service.
type SessionRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	DomainID string `json:"domain"`
}
