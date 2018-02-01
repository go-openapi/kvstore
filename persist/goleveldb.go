package persist

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
)

// NewGoLevelDBStore creates a new store backed by goleveldb
func NewGoLevelDBStore(cfg *viper.Viper) (Store, error) {
	db, err := leveldb.OpenFile(cfg.GetString("store.path"), nil)
	if err != nil {
		return nil, err
	}
	return &goleveldbStore{
		DB: db,
	}, nil
}

var (
	goleveldbSyncWrite   *opt.WriteOptions
	goleveldbNoCacheRead *opt.ReadOptions
)

func init() {
	goleveldbSyncWrite = &opt.WriteOptions{Sync: true}
	goleveldbNoCacheRead = &opt.ReadOptions{DontFillCache: true}
}

func goleveldbRewriteError(err error) error {
	switch err {
	case leveldb.ErrNotFound:
		return ErrNotFound
	case leveldb.ErrReadOnly:
		return ErrReadOnly
	case leveldb.ErrClosed:
		return ErrClosed
	case leveldb.ErrSnapshotReleased:
		return ErrSnapshotReleased
	case leveldb.ErrIterReleased:
		return ErrIterReleased
	default:
		return err
	}
}

func goleveldbRewriteValueError(value []byte, err error) (Value, error) {
	if err != nil {
		return Value{}, goleveldbRewriteError(err)
	}
	var result Value
	_, e := result.UnmarshalMsg(value)
	if e != nil {
		return Value{}, fmt.Errorf("msgp unmarshal failed: %v", e)
	}
	return result, nil
}

type goleveldbStore struct {
	DB *leveldb.DB
}

func (g *goleveldbStore) Put(key string, value *Value) error {
	// need this to be 0 when this is a new entry
	newVersion := value.Version

	prev, err := goleveldbRewriteValueError(g.DB.Get(UnsafeStringToBytes(key), goleveldbNoCacheRead))
	if err != nil {
		if err != ErrNotFound {
			return goleveldbRewriteError(err)
		}
		if err == ErrNotFound && newVersion != 0 {
			return ErrGone
		}
		value.Version = VersionOf(value.Value)
	}

	if prev.Version != newVersion {
		return ErrVersionMismatch
	}

	value.LastUpdated = time.Now().UTC().UnixNano()
	data, err := value.MarshalMsg(nil)
	if err != nil {
		return err
	}

	return goleveldbRewriteError(g.DB.Put(UnsafeStringToBytes(key), data, goleveldbSyncWrite))
}

func (g *goleveldbStore) Get(key string) (Value, error) {
	return goleveldbRewriteValueError(g.DB.Get(UnsafeStringToBytes(key), nil))
}

func (g *goleveldbStore) FindByPrefix(prefix string) ([]KeyValue, error) {
	var rg *util.Range
	if prefix != "" {
		rg = util.BytesPrefix(UnsafeStringToBytes(prefix))
	}

	iter := g.DB.NewIterator(rg, nil)
	var result []KeyValue
	for iter.Next() {
		value, err := goleveldbRewriteValueError(iter.Value(), nil)
		if err != nil {
			iter.Release()
			return nil, err
		}
		result = append(result, KeyValue{Key: UnsafeBytesToString(iter.Key()), Value: value})
	}
	iter.Release()

	err := iter.Error()
	if err != nil {
		return nil, goleveldbRewriteError(err)
	}
	return result, nil
}

func (g *goleveldbStore) Delete(key string) error {
	return goleveldbRewriteError(g.DB.Delete(UnsafeStringToBytes(key), goleveldbSyncWrite))
}

func (g *goleveldbStore) Close() error {
	return g.DB.Close()
}
