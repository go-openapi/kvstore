package persist

import (
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

func goleveldbRewriteValueError(value Value, err error) (Value, error) {
	return value, goleveldbRewriteError(err)
}

type goleveldbStore struct {
	DB *leveldb.DB
}

func (g *goleveldbStore) Put(key string, value Value) error {
	return goleveldbRewriteError(g.DB.Put([]byte(key), []byte(value), goleveldbSyncWrite()))
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
		result = append(result, KeyValue{Key: string(iter.Key()), Value: iter.Value()})
	}
	iter.Release()

	err := iter.Error()
	if err != nil {
		return nil, err
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
