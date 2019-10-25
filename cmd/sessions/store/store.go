package store

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"sync"
	"time"

	"github.com/fernandezvara/backd/internal/instrumentation"
	"github.com/fernandezvara/backd/internal/structs"
	"github.com/hashicorp/raft"
	"go.uber.org/zap"
)

// Store is the struct that holds the information on the nodes
// In order to replicate and distribute in a secure way the information
//   the keys and values are changed via distributed consensus (raft).
type Store struct {
	// where to bind the raft server
	raftBind string

	// Mutex to ensure safe concurrency
	mu sync.Mutex

	// m is the map that holds the sessions and ensures TTL are enforced
	m map[string]structs.Session

	// tick sets the duration between cleanups, this will run in every server
	//   being master or not, since deletions by TTL will be coherent if all
	//   servers are in time sync. This avoids overcharge of master servers.
	//   If there is an update of the session_id and that one was deleted,
	//   (s *Store) Set will insert it again.
	tick time.Duration

	// raft server itself (hashicorp/raft)
	raft *raft.Raft

	// instrumentation and logging
	inst *instrumentation.Instrumentation
}

const (
	raftTimeout = time.Duration(30 * time.Second)
)

// New returns the Store
func New(inst *instrumentation.Instrumentation, tick time.Duration) *Store {
	return &Store{
		inst: inst,
		tick: tick,
	}
}

// Open opens the store. If `enableSingle` this server will become leader
//   since it will hold all the work, but it will be a single point of failure.
// Use it only for tests.
func (s *Store) Open(enableSingle bool, id string, port uint16) error {

	var (
		// raft
		raftConfig    *raft.Config
		raftTransport *raft.NetworkTransport
		// raftSnapshots   *raft.InmemSnapshotStore
		// raftLogStore    raft.LogStore
		// raftStableStore raft.StableStore
		raftRaft *raft.Raft

		addr *net.TCPAddr

		err error
	)

	// set the bind port (uint16 = 0 to 65535)
	s.raftBind = fmt.Sprintf(":%d", port)

	raftConfig = raft.DefaultConfig()
	raftConfig.LocalID = raft.ServerID(id)

	addr, err = net.ResolveTCPAddr("tcp", s.raftBind)
	if err != nil {
		return err
	}
	// transport settings:
	//  connection pool = 3
	//  timeout         = 10 seconds
	raftTransport, err = raft.NewTCPTransport(s.raftBind, addr, 3, raftTimeout, os.Stderr)
	if err != nil {
		return err
	}

	raftSnapshots := raft.NewInmemSnapshotStore()
	raftLogStore := raft.NewInmemStore()
	raftStableStore := raft.NewInmemStore()

	raftRaft, err = raft.NewRaft(raftConfig, (*fsm)(s), raftLogStore, raftStableStore, raftSnapshots, raftTransport)
	if err != nil {
		return err
	}
	s.raft = raftRaft
	s.m = make(map[string]structs.Session)

	if enableSingle {

		configuration := raft.Configuration{
			Servers: []raft.Server{
				{
					ID:      raftConfig.LocalID,
					Address: raftTransport.LocalAddr(),
				},
			},
		}

		raftRaft.BootstrapCluster(configuration)

	}

	// start cleanup by TTL
	s.startTicker()

	return nil

}

// Get returns the value for the given key.
func (s *Store) Get(key string) (structs.Session, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.m[key], nil
}

// Set sets the value for the given key.
func (s *Store) Set(key string, value structs.Session) error {

	if s.raft.State() != raft.Leader {
		return ErrNotLeader
	}

	c := &command{
		Op:    "set",
		Key:   key,
		Value: value,
	}
	b, err := json.Marshal(c)
	if err != nil {
		return err
	}

	f := s.raft.Apply(b, raftTimeout)
	return f.Error()
}

// Delete deletes the given key.
func (s *Store) Delete(key string) error {
	if s.raft.State() != raft.Leader {
		return ErrNotLeader
	}

	c := &command{
		Op:  "delete",
		Key: key,
	}
	b, err := json.Marshal(c)
	if err != nil {
		return err
	}

	f := s.raft.Apply(b, raftTimeout)
	return f.Error()
}

// Join joins a node, identified by nodeID and located at addr, to this store.
// The node must be ready to respond to Raft communications at that address.
func (s *Store) Join(nodeID, addr string) error {

	s.inst.Info("raft", zap.String("op", "join request"), zap.String("node_id", nodeID), zap.String("node_addr", addr))
	// s.logger.Printf("received join request for remote node %s at %s", nodeID, addr)

	configFuture := s.raft.GetConfiguration()
	if err := configFuture.Error(); err != nil {
		s.inst.Error("raft", zap.String("err", err.Error()))
		return err
	}

	for _, srv := range configFuture.Configuration().Servers {
		// If a node already exists with either the joining node's ID or address,
		// that node may need to be removed from the config first.
		if srv.ID == raft.ServerID(nodeID) || srv.Address == raft.ServerAddress(addr) {
			// However if *both* the ID and the address are the same, then nothing -- not even
			// a join operation -- is needed.
			if srv.Address == raft.ServerAddress(addr) && srv.ID == raft.ServerID(nodeID) {
				s.inst.Info("raft", zap.String("op", "ignore join request"), zap.String("node_id", nodeID), zap.String("node_addr", addr))
				return nil
			}

			future := s.raft.RemoveServer(srv.ID, 0, 0)
			if err := future.Error(); err != nil {
				s.inst.Error("raft", zap.String("op", "remove node"), zap.String("node_id", nodeID), zap.String("node_addr", addr), zap.String("err", err.Error()))
				return ErrRemovingNode
			}
		}
	}

	f := s.raft.AddVoter(raft.ServerID(nodeID), raft.ServerAddress(addr), 0, 0)
	if f.Error() != nil {
		return f.Error()
	}
	s.inst.Info("raft", zap.String("op", "join success"), zap.String("node_id", nodeID), zap.String("node_addr", addr))
	return nil
}
