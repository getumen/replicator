package replicator

import (
	"bytes"
	"testing"
	"time"

	"github.com/getumen/replicator/pkg/handler/mockhandler"
	"github.com/getumen/replicator/pkg/logger/mocklogger"
	"github.com/getumen/replicator/pkg/logstore/mocklogstore"
	"github.com/getumen/replicator/pkg/models"
	"github.com/getumen/replicator/pkg/serializer/mockserializer"
	"github.com/getumen/replicator/pkg/snapshotstore/mocksnapshotstore"
	"github.com/getumen/replicator/pkg/snapshotter/mocksnapshotter"
	"github.com/getumen/replicator/pkg/stablestore/mockstablestore"
	"github.com/getumen/replicator/pkg/store/mockstore"
	"github.com/getumen/replicator/pkg/transport/mocktransport"
	"github.com/hashicorp/raft"
)

func TestReplicator_InMemSignleNode(t *testing.T) {
	var localAddr raft.ServerAddress = "127.0.0.1:0"
	store, err := mockstore.New()
	if err != nil {
		t.Fatal(err)
	}
	snapshotStore := mocksnapshotstore.NewInMemSnapshotStore()
	_, transport := mocktransport.NewInMemTransport(localAddr)
	logStore := mocklogstore.NewInmemStore()
	stableStore := mockstablestore.NewInmemStore()
	serializer := mockserializer.NewSerializer()
	snapshotter := mocksnapshotter.NewJSONSnapshotter()
	handler := mockhandler.NewKVSHandler()
	logger := mocklogger.NewLogger()

	replicator, err := NewReplicator(
		store,
		snapshotStore,
		transport,
		logStore,
		stableStore,
		serializer,
		snapshotter,
		handler,
		logger,
	)

	if err != nil {
		t.Fatal(err)
	}

	err = replicator.Start(string(localAddr), true)
	if err != nil {
		t.Fatal(err)
	}
	// Simple way to ensure there is a leader.
	time.Sleep(3 * time.Second)

	err = replicator.Propose(
		&models.Command{
			Operation: "PUT",
			Key:       []byte("foo"),
			Value:     []byte("bar"),
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	// Wait for committed log entry to be applied.
	time.Sleep(500 * time.Millisecond)
	snap, err := replicator.GetSnapshot()
	if err != nil {
		t.Fatal(err)
	}
	value, err := snap.Get([]byte("foo"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(value, []byte("bar")) {
		t.Fatalf("key has wrong value: %v", value)
	}

}
