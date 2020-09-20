package snapshotstore

import "github.com/hashicorp/raft"

// SnapshotStore is a snapshot of raft log
type SnapshotStore interface {
	raft.SnapshotStore
}
