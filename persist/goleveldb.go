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

func goleveldbSyncWrite() *opt.WriteOptions {
	return &opt.WriteOptions{Sync: true}
}

func goleveldbRewriteError(err error) error {
	if err == leveldb.ErrNotFound {
		return ErrNotFound
	}
	if err == leveldb.ErrReadOnly {
		return ErrReadOnly
	}
	if err == leveldb.ErrClosed {
		return ErrClosed
	}
	if err == leveldb.ErrSnapshotReleased {
		return ErrSnapshotReleased
	}
	if err == ErrIterReleased {
		return ErrIterReleased
	}
	return err
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

func (g *goleveldbStore) Put(key string, value Value) error {
	// need this to be 0 when this is a new entry
	newVersion := value.Version

	opts := opt.ReadOptions{DontFillCache: true}
	prev, err := goleveldbRewriteValueError(g.DB.Get([]byte(key), &opts))
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

	return goleveldbRewriteError(g.DB.Put([]byte(key), data, goleveldbSyncWrite()))
}

func (g *goleveldbStore) Get(key string) (Value, error) {
	return goleveldbRewriteValueError(g.DB.Get([]byte(key), nil))
}

func (g *goleveldbStore) FindByPrefix(prefix string) ([]KeyValue, error) {
	var rg *util.Range
	if prefix != "" {
		rg = util.BytesPrefix([]byte(prefix))
	}
	iter := g.DB.NewIterator(rg, nil)
	var result []KeyValue
	for iter.Next() {
		value, err := goleveldbRewriteValueError(iter.Value(), nil)
		if err != nil {
			iter.Release()
			return nil, err
		}
		result = append(result, KeyValue{Key: string(iter.Key()), Value: value})
	}

	err := iter.Error()
	if err != nil {
		return nil, goleveldbRewriteError(err)
	}
	return result, nil
}

func (g *goleveldbStore) Delete(key string) error {
	return goleveldbRewriteError(g.DB.Delete([]byte(key), goleveldbSyncWrite()))
}

func (g *goleveldbStore) Close() error {
	if g.DB == nil {
		return nil
	}
	return g.DB.Close()
}
