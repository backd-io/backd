package structs

import "time"

// Session is the struct that reflects the information of the user
//   currently logged into the domain
type Session struct {
	ID        string `json:"_id"`
	DomainID  string `json:"did"`
	User      User   `json:"uid"`
	ExpiresAt int64  `json:"eat"`
	CreatedAt int64  `json:"cat"`
}

// IsExpired returns expiration status of the session
func (s *Session) IsExpired() bool {
	return time.Now().Unix() > s.ExpiresAt
}

// SessionResponse is the struct that will be returned to the client
//   when a session has been established
type SessionResponse struct {
	ID        string `json:"_id"`
	DomainID  string `json:"domain_id"`
	UserID    string `json:"user_id"`
	ExpiresAt int64  `json:"expires_at"`
}

// SessionRequest is the struct that defines how an user creates a session
//   on the auth service.
type SessionRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	DomainID string `json:"domain"`
}
