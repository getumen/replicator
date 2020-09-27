package snapshotsink

//go:generate mockgen -source=$GOFILE -destination=mock$GOPACKAGE/mock_$GOFILE -package=mock$GOPACKAGE

import "github.com/hashicorp/raft"

// SnapshotSink is wrapper of raft.SnapshotSink
// this interface is for generating mock
type SnapshotSink raft.SnapshotSink
