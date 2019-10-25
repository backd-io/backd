package store

import "github.com/fernandezvara/backd/internal/structs"

type command struct {
	Op    string          `json:"op,omitempty"`
	Key   string          `json:"key,omitempty"`
	Value structs.Session `json:"value,omitempty"`
}
