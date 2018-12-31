package store

import (
	"fmt"
	"time"
)

func (s *Store) startTicker() {

	ticker := time.Tick(s.tick)

	go (func() {
		for {
			select {
			case <-ticker:
				s.removeExpired()
			}
		}
	})()

}

// removeExpired is called periodically to remove those sessions already expired
func (s *Store) removeExpired() {

	s.mu.Lock()
	defer s.mu.Unlock()

	var expired int64
	now := time.Now()

	for key, value := range s.m {
		if value.IsExpired() {
			delete(s.m, key)
			expired++
		}
	}

	fmt.Println("Expired keys:", expired, "Remain:", len(s.m), "Took:", time.Since(now).String())
}
