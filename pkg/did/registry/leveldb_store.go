package registry

import (
	"context"
	"errors"
	"github.com/syndtr/goleveldb/leveldb"
)

// LevelDBStore implements Store using LevelDB.
type LevelDBStore struct {
	db *leveldb.DB
}

func NewLevelDBStore(db *leveldb.DB) (*LevelDBStore, error) {
	if db == nil {
		return nil, errors.New("leveldb db is nil")
	}
	return &LevelDBStore{db: db}, nil
}

func (s *LevelDBStore) Put(_ context.Context, did string, document []byte) error {
	return s.db.Put([]byte(did), document, nil)
}

func (s *LevelDBStore) Get(_ context.Context, did string) ([]byte, error) {
	return s.db.Get([]byte(did), nil)
}

func (s *LevelDBStore) Has(_ context.Context, did string) (bool, error) {
	return s.db.Has([]byte(did), nil)
}

func (s *LevelDBStore) Close() error {
	return s.db.Close()
}
