package stablestore

//go:generate mockgen -source=$GOFILE -destination=mock$GOPACKAGE/mock_$GOFILE -package=mock$GOPACKAGE

import "github.com/hashicorp/raft"

// StableStore is raft log stable store
type StableStore interface {
	raft.StableStore
}
