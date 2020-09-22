package mockstore

import (
	"fmt"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/storage"

	"github.com/getumen/replicator/pkg/store"
)

type leveldbStore struct {
	internal *leveldb.DB
}

// New creates leveldb Store
func New() (store.Store, error) {

	s := storage.NewMemStorage()

	db, err := leveldb.Open(s, nil)
	if err != nil {
		return nil, err
	}
	return &leveldbStore{
		internal: db,
	}, nil
}

func (s *leveldbStore) Snapshot() (store.Snapshot, error) {

	snap, err := s.internal.GetSnapshot()

	if err != nil {
		return nil, err
	}

	return &snapshot{
		internal: snap,
	}, nil
}

func (s *leveldbStore) DiscardAll() error {
	snap, err := s.internal.GetSnapshot()

	if err != nil {
		return err
	}

	defer snap.Release()

	iter := snap.NewIterator(nil, nil)

	for iter.Next() {
		err = s.internal.Delete(iter.Key(), nil)
		if err != nil {
			return err
		}
	}

	err = iter.Error()
	if err != nil {
		return err
	}

	return nil
}

func (s *leveldbStore) CreateBatch() store.Batch {
	return &batch{
		internal: new(leveldb.Batch),
	}
}

func (s *leveldbStore) Write(ba store.Batch) error {
	if b, ok := ba.(*batch); ok {
		return s.internal.Write(b.internal, nil)
	}
	return fmt.Errorf("fail to cast batch")
}

type snapshot struct {
	internal *leveldb.Snapshot
}

func (s *snapshot) NewIterator() store.Iterator {
	return &iter{
		internal: s.internal.NewIterator(
			nil,
			&opt.ReadOptions{
				DontFillCache: true,
			}),
	}
}

func (s *snapshot) Get(key []byte) ([]byte, error) {
	value, err := s.internal.Get(key, nil)
	if err == leveldb.ErrNotFound {
		return nil, store.ErrNotFound
	} else if err != nil {
		return nil, err
	}
	return value, nil
}
func (s *snapshot) Has(key []byte) (bool, error) {
	value, err := s.internal.Has(key, nil)
	if err == leveldb.ErrNotFound {
		return false, store.ErrNotFound
	} else if err != nil {
		return false, err
	}
	return value, nil
}
func (s *snapshot) Release() {
	s.internal.Release()
}

type batch struct {
	internal *leveldb.Batch
}

func (b *batch) Len() int {
	return b.internal.Len()
}

func (b *batch) Put(key, value []byte) {
	b.internal.Put(key, value)
}

func (b *batch) Delete(key []byte) {
	b.internal.Delete(key)
}

func (b *batch) Reset() {
	b.internal.Reset()
}

type iter struct {
	internal iterator.Iterator
}

func (i *iter) Key() []byte {
	return i.internal.Key()
}

func (i *iter) Value() []byte {
	return i.internal.Value()
}

func (i *iter) Error() error {
	return i.internal.Error()
}

func (i *iter) First() bool {
	return i.internal.First()
}

func (i *iter) Last() bool {
	return i.internal.Last()
}

func (i *iter) Seek(key []byte) bool {
	return i.internal.Seek(key)
}

func (i *iter) Next() bool {
	return i.internal.Next()
}

func (i *iter) Prev() bool {
	return i.internal.Prev()
}

func (i *iter) Release() {
	i.internal.Release()
}
