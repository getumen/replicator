package snapshotstore

//go:generate mockgen -source=$GOFILE -destination=mock$GOPACKAGE/mock_$GOFILE -package=mock$GOPACKAGE

import "github.com/hashicorp/raft"

// SnapshotStore is a snapshot of raft log
type SnapshotStore interface {
	raft.SnapshotStore
}
