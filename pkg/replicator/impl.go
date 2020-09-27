package replicator

import (
	"fmt"
	"io"
	"time"

	"github.com/getumen/replicator/pkg/handler"
	"github.com/getumen/replicator/pkg/logger"
	"github.com/getumen/replicator/pkg/logstore"
	"github.com/getumen/replicator/pkg/models"
	"github.com/getumen/replicator/pkg/serializer"
	"github.com/getumen/replicator/pkg/snapshotstore"
	"github.com/getumen/replicator/pkg/snapshotter"
	"github.com/getumen/replicator/pkg/stablestore"
	"github.com/getumen/replicator/pkg/store"
	"github.com/getumen/replicator/pkg/transport"
	"github.com/hashicorp/raft"
)

const (
	retainSnapshotCount = 2
	raftTimeout         = 10 * time.Second
)

// ErrNotLeader represents this node is not leader
var ErrNotLeader = fmt.Errorf("not leader")

// Replicator is a replication functionality by raft algorithm
type replicator struct {
	raft *raft.Raft
	// store must be thread safe
	store         store.Store
	snapshotStore snapshotstore.SnapshotStore
	transport     transport.Transport
	logStore      logstore.LogStore
	stableStore   stablestore.StableStore
	serializer    serializer.Serializer
	snapshotter   snapshotter.Snapshotter
	handler       handler.Handler
	logger        logger.Logger
}

// NewReplicator returns replicator implementation
func NewReplicator(
	store store.Store,
	snapshotStore snapshotstore.SnapshotStore,
	transport transport.Transport,
	logStore logstore.LogStore,
	stableStore stablestore.StableStore,
	serializer serializer.Serializer,
	snapshotter snapshotter.Snapshotter,
	handler handler.Handler,
	logger logger.Logger,
) (Replicator, error) {
	return &replicator{
		store:         store,
		snapshotStore: snapshotStore,
		transport:     transport,
		logStore:      logStore,
		stableStore:   stableStore,
		serializer:    serializer,
		snapshotter:   snapshotter,
		handler:       handler,
		logger:        logger,
	}, nil
}

// Start prepares a replicator
func (r *replicator) Start(
	localID string,
	enableSingle bool,

) error {

	// Setup Raft configuration.
	config := raft.DefaultConfig()
	config.LocalID = raft.ServerID(localID)

	ra, err := raft.NewRaft(
		config,
		r,
		r.logStore,
		r.stableStore,
		r.snapshotStore,
		r.transport,
	)
	if err != nil {
		return fmt.Errorf("new raft: %s", err)
	}
	r.raft = ra

	if enableSingle {
		configuration := raft.Configuration{
			Servers: []raft.Server{
				{
					ID:      config.LocalID,
					Address: r.transport.LocalAddr(),
				},
			},
		}
		ra.BootstrapCluster(configuration)
	}

	return nil

}

// Join joins a node, identified by nodeID and located at addr, to this store.
// The node must be ready to respond to Raft communications at that address.
func (r *replicator) Join(nodeID, addr string) error {
	r.logger.WithFields(
		logger.Fields{
			"nodeID": nodeID,
			"addr":   addr,
		},
	).Info(
		"received join request for remote node",
	)

	configFuture := r.raft.GetConfiguration()
	if err := configFuture.Error(); err != nil {
		r.logger.WithFields(
			logger.Fields{
				"err": err,
			},
		).Warning(
			"failed to get raft configuration",
		)
		return err
	}

	for _, srv := range configFuture.Configuration().Servers {
		// If a node already exists with either the joining node's ID or address,
		// that node may need to be removed from the config first.
		if srv.ID == raft.ServerID(nodeID) || srv.Address == raft.ServerAddress(addr) {
			// However if *both* the ID and the address are the same, then nothing -- not even
			// a join operation -- is needed.
			if srv.Address == raft.ServerAddress(addr) && srv.ID == raft.ServerID(nodeID) {
				r.logger.WithFields(
					logger.Fields{
						"nodeID": nodeID,
						"addr":   addr,
					},
				).Info(
					"node already member of cluster, ignoring join request",
				)
				return nil
			}

			future := r.raft.RemoveServer(srv.ID, 0, 0)
			if err := future.Error(); err != nil {
				return fmt.Errorf("error removing existing node %s at %s: %s",
					nodeID, addr, err)
			}
		}
	}

	f := r.raft.AddVoter(raft.ServerID(nodeID), raft.ServerAddress(addr), 0, 0)
	if f.Error() != nil {
		return f.Error()
	}
	r.logger.WithFields(
		logger.Fields{
			"nodeID": nodeID,
			"addr":   addr,
		},
	).Info(
		"node joined successfully",
	)
	return nil
}

// Apply log is invoked once a log entry is committed.
// It returns a value which will be made available in the
// ApplyFuture returned by Raft.Apply method if that
// method was called on the same Raft node as the FSM.
func (r *replicator) Apply(raftLog *raft.Log) interface{} {

	command, err := r.serializer.Deserialize(raftLog.Data)
	if err != nil {
		r.logger.WithFields(
			logger.Fields{
				"err": err,
			},
		).Fatal(
			"fail to deserialize",
		)
	}

	err = r.handler.Apply(r.store, command)

	if err != nil {
		r.logger.WithFields(
			logger.Fields{
				"err": err,
			},
		).Fatal(
			"fail to apply command",
		)
	}

	return nil
}

// Snapshot is used to support log compaction. This call should
// return an FSMSnapshot which can be used to save a point-in-time
// snapshot of the FSM. Apply and Snapshot are not called in multiple
// threads, but Apply will be called concurrently with Persist. This means
// the FSM should be implemented in a fashion that allows for concurrent
// updates while a snapshot is happening.
func (r *replicator) Snapshot() (raft.FSMSnapshot, error) {
	return r.snapshotter.CreateSnapshot(r.store)
}

// Restore is used to restore an FSM from a snapshot. It is not called
// concurrently with any other command. The FSM must discard all previous
// state.
func (r *replicator) Restore(reader io.ReadCloser) error {
	return r.snapshotter.Restore(r.store, reader)
}

// ApplyBatch is invoked once a batch of log entries has been committed and
// are ready to be applied to the FSM. ApplyBatch will take in an array of
// log entries. These log entries will be in the order they were committed,
// will not have gaps, and could be of a few log types. Clients should check
// the log type prior to attempting to decode the data attached. Presently
// the LogCommand and LogConfiguration types will be sent.
//
// The returned slice must be the same length as the input and each response
// should correlate to the log at the same index of the input. The returned
// values will be made available in the ApplyFuture returned by Raft.Apply
// method if that method was called on the same Raft node as the FSM.
func (r *replicator) ApplyBatch(logs []*raft.Log) []interface{} {
	commands := make([]*models.Command, 0)

	for i := range logs {
		if logs[i].Type == raft.LogCommand {
			command, err := r.serializer.Deserialize(logs[i].Data)
			if err != nil {
				r.logger.WithFields(
					logger.Fields{
						"err": err,
					},
				).Fatal(
					"fail to deserialize",
				)
			}
			commands = append(commands, command)
		} else if logs[i].Type == raft.LogConfiguration {
			// TODO: do something
		} else {
			r.logger.Fatal("unexpected log type")
		}
	}

	err := r.handler.ApplyBatch(r.store, commands)

	if err != nil {
		if err != nil {
			r.logger.WithFields(
				logger.Fields{
					"err": err,
				},
			).Fatal(
				"fail to apply commands",
			)
		}
	}
	retVal := make([]interface{}, len(logs))
	return retVal
}

func (r *replicator) Propose(command *models.Command) error {
	if r.raft.State() != raft.Leader {
		return fmt.Errorf("not leader")
	}

	b, err := r.serializer.Serialize(command)
	if err != nil {
		return err
	}
	f := r.raft.Apply(b, raftTimeout)
	return f.Error()
}

func (r *replicator) GetSnapshot() (store.Snapshot, error) {
	return r.store.Snapshot()
}

func (r *replicator) IsLeader() bool {
	return r.raft.State() == raft.Leader
}

func (r *replicator) Leader() string {
	return string(r.raft.Leader())
}

func (r *replicator) LeaderCh() <-chan bool {
	return r.raft.LeaderCh()
}
