package logstore

//go:generate mockgen -source=$GOFILE -destination=mock$GOPACKAGE/mock_$GOFILE -package=mock$GOPACKAGE

import "github.com/hashicorp/raft"

// LogStore is raft log store
type LogStore interface {
	raft.LogStore
}
