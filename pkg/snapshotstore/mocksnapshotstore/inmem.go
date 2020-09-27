package mocksnapshotstore

import (
	"github.com/getumen/replicator/pkg/snapshotstore"
	"github.com/hashicorp/raft"
)

// NewInMemSnapshotStore returns in-memory SnapshotStore
func NewInMemSnapshotStore() snapshotstore.SnapshotStore {
	return raft.NewInmemSnapshotStore()
}
