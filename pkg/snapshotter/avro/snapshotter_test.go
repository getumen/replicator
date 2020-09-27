package avro

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"reflect"
	"testing"

	"github.com/getumen/replicator/pkg/snapshotsink/mocksnapshotsink"
	"github.com/getumen/replicator/pkg/store/mockstore"
	"github.com/golang/mock/gomock"
	"github.com/linkedin/goavro/v2"
)

func TestSnapshotter_CreateSnapshot(t *testing.T) {

	key := []byte("key")
	value := []byte("value")

	store, err := mockstore.New()
	if err != nil {
		t.Fatal(err)
	}

	batch := store.CreateBatch()
	batch.Put(key, value)
	err = store.Write(batch)
	if err != nil {
		t.Fatal(err)
	}

	target := &snapshotterImpl{}

	snapshot, err := target.CreateSnapshot(store)

	if err != nil {
		t.Fatal(err)
	}

	if snap, ok := snapshot.(*fsmSnapshot); ok {
		defer snap.snapshot.Release()
		v, err := snap.snapshot.Get(key)
		if err != nil {
			t.Fatal(err)
		}
		if bytes.Compare(v, value) != 0 {
			t.Fatalf("the value is invalid: %v", v)
		}
	}
}

func TestSnapshotter_Restore(t *testing.T) {
	r, w := io.Pipe()
	store, err := mockstore.New()
	if err != nil {
		t.Fatal(err)
	}
	target := &snapshotterImpl{}

	go func() {
		defer w.Close()
		codec, err := goavro.NewCodec(schema)
		if err != nil {
			log.Fatal(err)
		}
		writer, err := goavro.NewOCFWriter(
			goavro.OCFConfig{
				W:               w,
				Codec:           codec,
				CompressionName: "snappy",
			},
		)
		if err != nil {
			log.Fatal(err)
		}
		block := []interface{}{
			map[string]interface{}{
				"key":   []byte("key2"),
				"value": []byte("value"),
			},
			map[string]interface{}{
				"key":   []byte("key1"),
				"value": []byte{},
			},
		}
		err = writer.Append(block)
		if err != nil {
			log.Fatal(err)
		}
	}()

	err = target.Restore(store, r)
	if err != nil {
		t.Fatal(err)
	}

	snap, err := store.Snapshot()
	if err != nil {
		t.Fatal(err)
	}

	v, err := snap.Get([]byte("key2"))
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Compare(v, []byte("value")) != 0 {
		t.Fatalf("expected value, but got %v", v)
	}
	exists, err := snap.Has([]byte("key1"))
	if err != nil {
		t.Fatal(err)
	}
	if !exists {
		t.Fatalf("expected false, but got true")
	}
}

func TestFsmSnapshot_Persist(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store, err := mockstore.New()
	if err != nil {
		t.Fatal(err)
	}
	batch := store.CreateBatch()

	for i := 0; i < 10; i++ {
		batch.Put(
			[]byte(fmt.Sprintf("key-%d", i)),
			[]byte(fmt.Sprintf("value-%d", i)),
		)
	}
	err = store.Write(batch)
	if err != nil {
		t.Fatal(err)
	}

	snapshot, err := store.Snapshot()
	if err != nil {
		t.Fatal(err)
	}
	target := &fsmSnapshot{
		snapshot: snapshot,
	}

	snapshotSink := mocksnapshotsink.NewMockSnapshotSink(ctrl)

	buf := new(bytes.Buffer)
	snapshotSink.
		EXPECT().
		Write(gomock.Any()).
		Do(func(b []byte) (int, error) {
			return buf.Write(b)
		}).AnyTimes()

	snapshotSink.EXPECT().Close()

	target.Persist(snapshotSink)

	r, err := goavro.NewOCFReader(buf)
	if err != nil {
		t.Fatal(err)
	}

	counter := 0
	for r.Scan() {
		record, err := r.Read()
		if err != nil {
			t.Fatal(err)
		}

		if m, ok := record.(map[string][]byte); ok {
			reflect.DeepEqual(
				map[string][]byte{
					"key":   []byte(fmt.Sprintf("key-%d", counter)),
					"value": []byte(fmt.Sprintf("value-%d", counter)),
				},
				m,
			)
		}
		counter++
	}

	if counter != 10 {
		t.Fatal("value number is not match")
	}
}
