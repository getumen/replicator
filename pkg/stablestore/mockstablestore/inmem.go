package mockstablestore

import (
	"github.com/getumen/replicator/pkg/stablestore"
	"github.com/hashicorp/raft"
)

// NewInmemStore returns in-memory stable store
func NewInmemStore() stablestore.StableStore {
	return raft.NewInmemStore()
}
