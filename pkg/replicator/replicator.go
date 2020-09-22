package replicator

import (
	"fmt"
	"io"

	"github.com/getumen/replicator/pkg/handler"
	"github.com/getumen/replicator/pkg/log"
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

// Replicator is a replication functionality by raft algorithm
type Replicator struct {
	raft *raft.Raft

	// store must be thread safe
	store store.Store

	snapshots snapshotstore.SnapshotStore

	transport transport.Transport

	logStore logstore.LogStore

	stableStore stablestore.StableStore

	serializer serializer.Serializer

	snapshotter snapshotter.Snapshotter

	handler handler.Handler

	logger log.Logger

	RaftBind string
}

// Start prepares a replicator
func (r *Replicator) Start(
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
		r.snapshots,
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
func (r *Replicator) Join(nodeID, addr string) error {
	r.logger.Info(
		"received join request for remote node",
		log.Fields{
			"nodeID": nodeID,
			"addr":   addr,
		},
	)

	configFuture := r.raft.GetConfiguration()
	if err := configFuture.Error(); err != nil {
		r.logger.Warning(
			"failed to get raft configuration",
			log.Fields{
				"err": err,
			},
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
				r.logger.Info(
					"node already member of cluster, ignoring join request",
					log.Fields{
						"nodeID": nodeID,
						"addr":   addr,
					},
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
	r.logger.Info(
		"node joined successfully",
		log.Fields{
			"nodeID": nodeID,
			"addr":   addr,
		},
	)
	return nil
}

// Apply log is invoked once a log entry is committed.
// It returns a value which will be made available in the
// ApplyFuture returned by Raft.Apply method if that
// method was called on the same Raft node as the FSM.
func (r *Replicator) Apply(raftLog *raft.Log) interface{} {
	command, err := r.serializer.Deserialize(raftLog.Data)
	if err != nil {
		r.logger.Fatal(
			"fail to deserialize",
			log.Fields{
				"err": err,
			},
		)
	}

	err = r.handler.Apply(r.store, command)

	if err != nil {
		r.logger.Fatal(
			"fail to apply command",
			log.Fields{
				"err": err,
			},
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
func (r *Replicator) Snapshot() (raft.FSMSnapshot, error) {
	return r.snapshotter.CreateSnapshot(r.store)
}

// Restore is used to restore an FSM from a snapshot. It is not called
// concurrently with any other command. The FSM must discard all previous
// state.
func (r *Replicator) Restore(reader io.ReadCloser) error {
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
func (r *Replicator) ApplyBatch(logs []*raft.Log) []interface{} {
	commands := make([]*models.Command, len(logs))

	for i := range logs {
		command, err := r.serializer.Deserialize(logs[i].Data)
		if err != nil {
			r.logger.Fatal(
				"fail to deserialize",
				log.Fields{
					"err": err,
				},
			)
		}
		commands[i] = command
	}

	err := r.handler.ApplyBatch(r.store, commands)

	if err != nil {
		if err != nil {
			r.logger.Fatal(
				"fail to apply commands",
				log.Fields{
					"err": err,
				},
			)
		}
	}
	return nil
}
