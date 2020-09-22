package avro

import (
	"fmt"
	"io"
	"log"

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

type snapshotter struct{}

func (s *snapshotter) CreateSnapshot(store store.Store) (raft.FSMSnapshot, error) {

	snapshot, err := store.Snapshot()

	if err != nil {
		return nil, err
	}

	return &fsmSnapshot{
		snapshot: snapshot,
	}, nil
}

func (s *snapshotter) Restore(store store.Store, reader io.ReadCloser) error {
	defer reader.Close()

	err := store.DiscardAll()
	if err != nil {
		return err
	}

	r, err := goavro.NewOCFReader(reader)
	if err != nil {
		return err
	}

	batch := store.CreateBatch()

	for r.Scan() {
		record, err := r.Read()
		if err != nil {
			return err
		}
		if m, ok := record.(map[string]interface{}); ok {
			log.Println(m)
			key := m["key"].([]byte)
			if value, ok := m["value"]; ok {
				if valueMap, ok := value.([]byte); ok {
					batch.Put(key, valueMap)
				} else {
					return fmt.Errorf("unsupported type: %v", value)
				}
			} else {
				return fmt.Errorf("value not found")
			}
		} else {
			return fmt.Errorf("fail to cast record")
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
		sink.Cancel()
	}

	return err
}

func (f *fsmSnapshot) Release() {
	f.snapshot.Release()
}
