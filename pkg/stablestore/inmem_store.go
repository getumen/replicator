package stablestore

import "github.com/hashicorp/raft"

// NewInmemStore returns in-memory stable store
func NewInmemStore() StableStore {
	return raft.NewInmemStore()
}
