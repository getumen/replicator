package logstore

import "github.com/hashicorp/raft"

// NewInmemStore returns in-memory log store
func NewInmemStore() LogStore {
	return raft.NewInmemStore()
}
