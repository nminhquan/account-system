package db

import (
	"errors"
	"log"
	"strconv"

	"github.com/tecbot/gorocksdb"
)

var (
	ErrRockDBKeyNotExist = errors.New("Key does not exist")
)

type KeyValDB interface {
	Set(key interface{}, value interface{}) error
	Get(key interface{}) (interface{}, error)
}

type RocksDB struct {
	path  string
	db    *gorocksdb.DB
	read  *gorocksdb.ReadOptions
	write *gorocksdb.WriteOptions
}

func NewRocksDB(path string) (*RocksDB, error) {
	log.Println("Create RocksDB under ", path)
	bbto := gorocksdb.NewDefaultBlockBasedTableOptions()
	bbto.SetBlockCache(gorocksdb.NewLRUCache(3 << 30))

	opts := gorocksdb.NewDefaultOptions()
	opts.SetBlockBasedTableFactory(bbto)
	opts.SetCreateIfMissing(true)
	opts.IncreaseParallelism(4)

	database, err := gorocksdb.OpenDb(opts, path)
	if err != nil {
		return nil, err
	}

	read := gorocksdb.NewDefaultReadOptions()
	write := gorocksdb.NewDefaultWriteOptions()

	rocksDB := RocksDB{db: database, read: read, write: write}

	return &rocksDB, nil
}

func (rocksDB *RocksDB) Close() {
	rocksDB.db.Close()
	rocksDB.read.Destroy()
	rocksDB.write.Destroy()
}

func (rocksDB *RocksDB) GetAccountBalance(accountId string) (string, error) {
	data, err := rocksDB.db.Get(rocksDB.read, []byte(accountId))
	if data.Data() == nil || err != nil {
		err = ErrRockDBKeyNotExist
	}

	return string(data.Data()), err
}

func (rocksDB *RocksDB) SetAccountBalance(accountId string, amount float64) error {
	err := rocksDB.db.Put(rocksDB.write, []byte(accountId), []byte(strconv.FormatFloat(amount, 'f', 6, 64)))
	return err
}

func (rocksDB *RocksDB) Set(key string, value string) error {
	err := rocksDB.db.Put(rocksDB.write, []byte(key), []byte(value))
	return err
}

func (rocksDB *RocksDB) Get(key string) (string, error) {
	data, err := rocksDB.db.Get(rocksDB.read, []byte(key))
	return string(data.Data()), err
}

func (rocksDB *RocksDB) Del(key string) error {
	err := rocksDB.db.Delete(rocksDB.write, []byte(key))
	return err
}

func (rocksDB *RocksDB) DB() *gorocksdb.DB {
	return rocksDB.db
}
