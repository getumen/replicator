package snapshotter

//go:generate mockgen -source=$GOFILE -destination=mock$GOPACKAGE/mock_$GOFILE -package=mock$GOPACKAGE

import (
	"io"

	"github.com/getumen/replicator/pkg/store"
	"github.com/hashicorp/raft"
)

// Snapshotter creates snapshot of FSM
type Snapshotter interface {
	CreateSnapshot(store.Store) (raft.FSMSnapshot, error)
	Restore(store.Store, io.ReadCloser) error
}
