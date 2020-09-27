package avro

import (
	"fmt"
	"io"

	"github.com/getumen/replicator/pkg/snapshotter"
	"github.com/getumen/replicator/pkg/store"
	"github.com/hashicorp/raft"
	"github.com/linkedin/goavro/v2"
)

const blockSize = 1024
const schema = `
{
	"name": "snapshot",
	"type": "record",
	"fields": [
		{"name": "key", "type": "bytes"},
		{"name": "value", "type": "bytes"}
	]
}
`

type snapshotterImpl struct{}

// NewAvroSnapshotter returns avro snapshotter
func NewAvroSnapshotter() snapshotter.Snapshotter {
	return &snapshotterImpl{}
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

func (s *snapshotterImpl) Restore(store store.Store, reader io.ReadCloser) (err error) {
	defer func() {
		closeErr := reader.Close()
		if err == nil && closeErr != nil {
			err = closeErr
		}
	}()
	err = store.DiscardAll()
	if err != nil {
		return
	}

	r, err := goavro.NewOCFReader(reader)
	if err != nil {
		return
	}

	batch := store.CreateBatch()

	for r.Scan() {
		record, readErr := r.Read()
		if readErr != nil {
			err = readErr
			return
		}
		if m, ok := record.(map[string]interface{}); ok {
			key := m["key"].([]byte)
			if value, ok := m["value"]; ok {
				if valueMap, ok := value.([]byte); ok {
					batch.Put(key, valueMap)
				} else {
					err = fmt.Errorf("unsupported type: %v", value)
					return
				}
			} else {
				err = fmt.Errorf("value not found")
				return
			}
		} else {
			err = fmt.Errorf("fail to cast record")
			return
		}

		if batch.Len() == blockSize {
			err = store.Write(batch)
			if err != nil {
				return err
			}
			batch.Reset()
		}
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
		codec, err := goavro.NewCodec(schema)
		if err != nil {
			return err
		}

		writer, err := goavro.NewOCFWriter(
			goavro.OCFConfig{
				W:               sink,
				Codec:           codec,
				CompressionName: "snappy",
			},
		)

		if err != nil {
			return err
		}

		iter := f.snapshot.NewIterator()

		block := make([]interface{}, 0, blockSize)

		for iter.Next() {
			key := make([]byte, len(iter.Key()))
			copy(key, iter.Key())

			if iter.Value() != nil {
				value := make([]byte, len(iter.Value()))
				copy(value, iter.Value())

				block = append(block, map[string]interface{}{
					"key":   key,
					"value": value,
				})
			} else {
				block = append(block, map[string]interface{}{
					"key":   key,
					"value": []byte{},
				})
			}

			if len(block) == blockSize {
				err = writer.Append(block)
				if err != nil {
					return err
				}
				block = block[:0]
			}
		}

		// write the remaining block
		err = writer.Append(block)
		if err != nil {
			return err
		}

		return sink.Close()
	}()

	if err != nil {
		err = sink.Cancel()
		if err != nil {
			return err
		}
	}

	return err
}

func (f *fsmSnapshot) Release() {
	f.snapshot.Release()
}
