package mocklogstore

import (
	"github.com/getumen/replicator/pkg/logstore"
	"github.com/hashicorp/raft"
)

// NewInmemStore returns in-memory log store
func NewInmemStore() logstore.LogStore {
	return raft.NewInmemStore()
}
