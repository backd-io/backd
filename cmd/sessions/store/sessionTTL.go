package store

// // SessionTTL is the struct that manages the session and its Time to live
// type SessionTTL struct {
// 	*sync.RWMutex
// 	session *structs.Session
// 	expires time.Time
// }

// // isExpired is true when the expires time is in the past
// func (s *SessionTTL) isExpired() bool {

// 	s.RLock()
// 	defer s.RUnlock()

// 	return s.expires.Before(time.Now())

// }

// // setExpiration sets the desired expiration on the session
// func (s *SessionTTL) setExpiration(exp time.Time) {
// 	s.Lock()
// 	s.expires = exp
// 	s.session.ExpiresAt = exp.Unix()
// 	s.Unlock()
// }

// // addDuration sets the expiration time to the session by adding the especified time
// func (s *SessionTTL) addDuration(duration time.Duration) {
// 	s.Lock()
// 	s.expires = time.Now().Add(duration)
// 	s.session.ExpiresAt = s.expires.Unix()
// 	s.Unlock()
// }
