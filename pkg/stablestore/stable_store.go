package stablestore

import "github.com/hashicorp/raft"

// StableStore is raft log stable store
type StableStore interface {
	raft.StableStore
}
