package hashstorage

import (
	"github.com/rs/zerolog"
	"github.com/syndtr/goleveldb/leveldb"
	leveldbStorage "github.com/syndtr/goleveldb/leveldb/storage"
	"golang.org/x/xerrors"

	"github.com/spikeekips/avl"
)

var (
	leveldbPrefixRootKey []byte = []byte{0x00, 0x00}
	leveldbPrefixNodeKey []byte = []byte{0x00, 0x01}
	leveldbPrefixRawKey  []byte = []byte{0x00, 0x02}
)

type LevelDBStorage struct {
	*avl.Logger
	db *leveldb.DB
}

func NewLevelDBStorage(path string) (*LevelDBStorage, error) {
	var st leveldbStorage.Storage
	if path == "memory://" {
		st = leveldbStorage.NewMemStorage()
	} else {
		if s, err := leveldbStorage.OpenFile(path, false); err != nil {
			return nil, err
		} else {
			st = s
		}
	}
	db, err := leveldb.Open(st, nil)
	if err != nil {
		return nil, err
	}

	dbs := &LevelDBStorage{
		Logger: avl.NewLogger(func(c zerolog.Context) zerolog.Context {
			return c.Str("module", "avl_tree_leveldb_storage")
		}),
		db: db,
	}

	dbs.Log().Debug().
		Str("path", path).
		Msg("new leveldbStorage")

	return dbs, nil
}

func (ls *LevelDBStorage) GetRoot(name []byte) ([]byte, error) {
	b, err := ls.db.Get(ls.rootKey(name), nil)
	if xerrors.Is(err, leveldb.ErrNotFound) {
		return nil, nil
	}

	return b, err
}

func (ls *LevelDBStorage) SetRoot(name, rootKey []byte) error {
	ls.Log().Debug().
		Bytes("name", name).
		Msg("SetRoot")
	return ls.db.Put(ls.rootKey(name), rootKey, nil)
}

func (ls *LevelDBStorage) rootKey(name []byte) []byte {
	return append(leveldbPrefixRootKey[:], name...)
}

func (ls *LevelDBStorage) nodeKey(key []byte) []byte {
	return append(leveldbPrefixNodeKey[:], key...)
}

func (ls *LevelDBStorage) rawKey(key []byte) []byte {
	return append(leveldbPrefixRawKey[:], key...)
}

func (ls *LevelDBStorage) GetNode(key []byte) ([]byte, error) {
	ls.Log().Debug().
		Bytes("key", key).
		Msg("GetNode")

	b, err := ls.db.Get(ls.nodeKey(key), nil)
	if xerrors.Is(err, leveldb.ErrNotFound) {
		return nil, nil
	}

	return b, err
}

func (ls *LevelDBStorage) SetNode(key, value []byte) error {
	ls.Log().Debug().
		Bytes("key", key).
		Msg("SetNode")
	return ls.db.Put(ls.nodeKey(key), value, nil)
}

func (ls *LevelDBStorage) Batch() BatchStorage {
	return &LevelDBBatchStorage{
		ls:    ls,
		batch: new(leveldb.Batch),
	}
}

func (ls *LevelDBStorage) GetRaw(key []byte) ([]byte, error) {
	ls.Log().Debug().
		Bytes("key", key).
		Msg("GetRaw")
	b, err := ls.db.Get(ls.rawKey(key), nil)
	if xerrors.Is(err, leveldb.ErrNotFound) {
		return nil, nil
	}

	return b, err
}

func (ls *LevelDBStorage) SetRaw(key, value []byte) error {
	ls.Log().Debug().
		Bytes("key", key).
		Msg("SetRaw")
	return ls.db.Put(ls.rawKey(key), value, nil)
}

type LevelDBBatchStorage struct {
	ls    *LevelDBStorage
	batch *leveldb.Batch
}

func (lsb *LevelDBBatchStorage) SetRoot(name []byte, rootKey []byte) error {
	lsb.batch.Put(lsb.ls.rootKey(name), rootKey)
	return nil
}

func (lsb *LevelDBBatchStorage) SetNode(key, value []byte) error {
	lsb.batch.Put(lsb.ls.nodeKey(key), value)
	return nil
}

func (lsb *LevelDBBatchStorage) Commit() error {
	return lsb.ls.db.Write(lsb.batch, nil)
}
