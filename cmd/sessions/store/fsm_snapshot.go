package store

import (
	"encoding/json"

	"github.com/fernandezvara/backd/internal/structs"
	"github.com/hashicorp/raft"
)

type fsmSnapshot struct {
	store map[string]structs.Session
}

func (f *fsmSnapshot) Persist(sink raft.SnapshotSink) error {
	err := func() error {
		// Encode data.
		b, err := json.Marshal(f.store)
		if err != nil {
			return err
		}

		// Write data to sink.
		if _, err := sink.Write(b); err != nil {
			return err
		}

		// Close the sink.
		return sink.Close()
	}()

	if err != nil {
		sink.Cancel()
	}

	return err
}

func (f *fsmSnapshot) Release() {}
