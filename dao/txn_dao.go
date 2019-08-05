package dao

import (
	"encoding/json"
	"log"
	"time"

	"github.com/nminhquan/mas/utils"

	"github.com/nminhquan/mas/db"

	"github.com/nminhquan/mas/model"
)

type TxnCoordinatorDAO interface {
	GetMaxId() string
	IncrMaxId() int64
	GetPeersList() map[string]string
	CreateTransactionEntry(string, int64, string, string)
	CreateSubTransactionEntry(string, string, model.Instruction, string)
	CheckLock(lockId string) (int64, error)
	CreateLock(lockId string, state string, timeout time.Duration) (bool, error)
	RefreshLock(lockId string, timeout time.Duration) (bool, error)
	DeleteLock(lockId string) (int64, error)
	GetPeerBucket(id string) string
	InsertPeerBucket(id string, peer string) error
	InsertPeers(id string, peers string) (bool, error)
}

type TxnCoordinatorDAOImpl struct {
	cache   *db.CacheService
	rocksdb *db.RocksDB
}

func NewTxnCoordinatorDAO(cacheService *db.CacheService, lockDB *db.RocksDB) TxnCoordinatorDAO {
	return &TxnCoordinatorDAOImpl{cacheService, lockDB}
}

func (tcDao *TxnCoordinatorDAOImpl) InsertPeers(id string, peers string) (bool, error) {
	rs, err := tcDao.cache.HSet("peers", id, peers).Result()

	if err != nil {
		log.Println("cannot InsertPeers")
		return rs, err
	}

	return rs, nil
}

func (tcDao *TxnCoordinatorDAOImpl) GetMaxId() string {
	rs, err := tcDao.cache.Get("current_id").Result()

	if err != nil {
		log.Println("cannot get current_id")
	}

	return rs
}

func (tcDao *TxnCoordinatorDAOImpl) IncrMaxId() int64 {
	rs, err := tcDao.cache.Incr("current_id").Result()

	if err != nil {
		log.Println("cannot set current_id")
	}

	return rs
}

func (tcDao *TxnCoordinatorDAOImpl) CheckLock(lockId string) (int64, error) {
	result, err := tcDao.cache.Exists(lockId).Result()
	return result, err
}

func (tcDao *TxnCoordinatorDAOImpl) CreateLock(lockId string, state string, timeout time.Duration) (bool, error) {
	updated, err := tcDao.cache.SetNX(lockId, state, timeout).Result()
	if err != nil || !updated {
		return false, err
	}
	return true, nil
}

// func (tcDao *TxnCoordinatorDAOImpl) CheckLock(lockId string) (int64, error) {
// 	status, err := tcDao.rocksdb.Get(lockId)
// 	if status == "" {
// 		return 0, nil
// 	} else if err != nil {
// 		return 0, err
// 	}
// 	return 100, nil
// }

// func (tcDao *TxnCoordinatorDAOImpl) CreateLock(lockId string, state string, timeout time.Duration) (bool, error) {
// 	data, err := tcDao.rocksdb.Get(lockId)
// 	if data == "" {
// 		err = tcDao.rocksdb.Set(lockId, state)
// 		return true, nil
// 	}
// 	return false, err
// }

func (tcDao *TxnCoordinatorDAOImpl) RefreshLock(lockId string, timeout time.Duration) (bool, error) {
	updated, err := tcDao.cache.Expire(lockId, timeout).Result()
	if err != nil {
		return false, err
	}
	return updated, nil
}

func (tcDao *TxnCoordinatorDAOImpl) DeleteLock(lockId string) (int64, error) {
	result, err := tcDao.cache.Del(lockId).Result()
	if err != nil {
		return result, err
	}
	return result, nil
}

// func (tcDao *TxnCoordinatorDAOImpl) DeleteLock(lockId string) (int64, error) {
// 	err := tcDao.rocksdb.Del(lockId)
// 	if err != nil {
// 		return 0, err
// 	}
// 	return 100, nil
// }

func (tcDao *TxnCoordinatorDAOImpl) GetPeersList() map[string]string {
	rs, _ := tcDao.cache.HGetAll("peers").Result()
	return rs
}

// func (tcDao *TxnCoordinatorDAOImpl) GetPeerBucket(id string) string {
// 	rs, _ := tcDao.cache.HGet("buckets", id).Result()
// 	return rs
// }

// func (tcDao *TxnCoordinatorDAOImpl) InsertPeerBucket(id string, peer string) error {
// 	_, err := tcDao.cache.HSet("buckets", id, peer).Result()
// 	if err != nil {
// 		log.Println("cannot set peer Bucket for id: ", id, " value: ", peer)
// 	}
// 	return err
// }

func (tcDao *TxnCoordinatorDAOImpl) GetPeerBucket(id string) string {
	rs, _ := tcDao.rocksdb.Get("buckets:" + id)
	return rs
}

func (tcDao *TxnCoordinatorDAOImpl) InsertPeerBucket(id string, peer string) error {
	err := tcDao.rocksdb.Set("buckets:"+id, peer)
	if err != nil {
		log.Println("cannot set peer Bucket for id: ", id, " value: ", peer)
	}
	return err
}

func (tcDao *TxnCoordinatorDAOImpl) CreateTransactionEntry(txnId string, ts int64, state string, lockIds string) {
	txnHKey := "txn:" + txnId
	fields := map[string]interface{}{
		"ts":    ts,
		"state": state,
		"locks": lockIds,
	}
	var b, _ = json.Marshal(fields)
	err := tcDao.rocksdb.Set(txnHKey, string(b))

	if err != nil {
		log.Println("cannot CreateTransactionEntry")
	}
}

func (tcDao *TxnCoordinatorDAOImpl) CreateSubTransactionEntry(txnId string, state string, inst model.Instruction, globalTxnId string) {
	txnHKey := "sub_txn:" + globalTxnId + ":" + txnId
	fields := map[string]interface{}{
		"state": state,
		"inst":  utils.SerializeMessage(inst),
	}
	var b, _ = json.Marshal(fields)
	err := tcDao.rocksdb.Set(txnHKey, string(b))

	if err != nil {
		log.Println("cannot CreateSubTransactionEntry ", err)
	}
}
