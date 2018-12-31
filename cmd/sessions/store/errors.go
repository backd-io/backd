package store

import "errors"

var (
	// ErrNotLeader is returned when an operation that must be done by
	//   the cluster leader is being tryed on a follower
	ErrNotLeader = errors.New("RAFT: Not leader")
	// ErrRemovingNode is returned when service fails to remove a node from
	//   the cluster
	ErrRemovingNode = errors.New("RAFT: Error removing node")
)
