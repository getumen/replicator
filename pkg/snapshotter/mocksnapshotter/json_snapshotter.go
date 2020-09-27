package mocksnapshotter

import (
	"encoding/base64"
	"encoding/json"
	io "io"

	"github.com/getumen/replicator/pkg/snapshotter"
	"github.com/getumen/replicator/pkg/store"
	raft "github.com/hashicorp/raft"
)

// NewJSONSnapshotter is an implementation for tests
func NewJSONSnapshotter() snapshotter.Snapshotter {
	return &snapshotterImpl{}
}

type snapshotterImpl struct {
}

func (s *snapshotterImpl) CreateSnapshot(store store.Store) (raft.FSMSnapshot, error) {

	snapshot, err := store.Snapshot()

	if err != nil {
		return nil, err
	}

	return &fsmSnapshot{
		snapshot: snapshot,
	}, nil
}

func (s *snapshotterImpl) Restore(store store.Store, reader io.ReadCloser) error {
	defer reader.Close()

	err := store.DiscardAll()
	if err != nil {
		return err
	}

	batch := store.CreateBatch()

	var m map[string]string
	err = json.NewDecoder(reader).Decode(&m)
	if err != nil {
		return err
	}

	for key, value := range m {
		k, err := base64.StdEncoding.DecodeString(key)
		if err != nil {
			return err
		}
		v, err := base64.StdEncoding.DecodeString(value)
		if err != nil {
			return err
		}
		batch.Put(k, v)
	}

	// write the remaining block
	if batch.Len() > 0 {
		err = store.Write(batch)
		if err != nil {
			return err
		}
	}

	return nil
}

type fsmSnapshot struct {
	snapshot store.Snapshot
}

func (f *fsmSnapshot) Persist(sink raft.SnapshotSink) error {

	err := func() error {
		var m map[string]string

		iter := f.snapshot.NewIterator()

		for iter.Next() {
			key := base64.StdEncoding.EncodeToString(iter.Key())
			value := base64.StdEncoding.EncodeToString(iter.Value())

			m[key] = value
		}

		err := json.NewEncoder(sink).Encode(m)
		if err != nil {
			return err
		}

		return sink.Close()
	}()

	if err != nil {
		sink.Cancel()
	}

	return err
}

func (f *fsmSnapshot) Release() {
	f.snapshot.Release()
}
