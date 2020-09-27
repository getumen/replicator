package replicator

import (
	"github.com/getumen/replicator/pkg/models"
	"github.com/getumen/replicator/pkg/store"
)

// Replicator is storage engine replicated by raft algorithm
type Replicator interface {
	Join(nodeID, addr string) error
	Propose(command *models.Command) error
	GetSnapshot() (store.Snapshot, error)
	IsLeader() bool
	Leader() string
	LeaderCh() <-chan bool
	Start(string, bool) error
}
