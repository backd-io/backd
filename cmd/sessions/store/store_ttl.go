package store

import (
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
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

	s.inst.Metric("sessions_in_use").(*prometheus.GaugeVec).WithLabelValues(s.inst.Hostname()).Set(float64(len(s.m)))

	fmt.Println("Expired keys:", expired, "Remain:", len(s.m), "Took:", time.Since(now).String())
}
