package logstore

import "github.com/hashicorp/raft"

// LogStore is raft log store
type LogStore interface {
	raft.LogStore
}
