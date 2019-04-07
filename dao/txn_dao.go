package dao

import (
	"log"
	"time"

	"gitlab.zalopay.vn/quannm4/mas/db"
	"gitlab.zalopay.vn/quannm4/mas/model"
	"gitlab.zalopay.vn/quannm4/mas/utils"
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
	cache *db.CacheService
}

func NewTxnCoordinatorDAO(cacheService *db.CacheService) TxnCoordinatorDAO {
	return &TxnCoordinatorDAOImpl{cacheService}
}

func (tcDao *TxnCoordinatorDAOImpl) InsertPeers(id string, peers string) (bool, error) {
	rs, err := tcDao.cache.HSet("peers", id, peers).Result()

	if err != nil {
		log.Println("cannot get current_id")
		return rs, err
	}

	return rs, nil
}

func (tcDao *TxnCoordinatorDAOImpl) GetMaxId() string {
	rs, err := tcDao.cache.Get("current_id").Result()

	if err != nil {
		log.Fatalln("cannot get current_id")
	}

	return rs
}

func (tcDao *TxnCoordinatorDAOImpl) IncrMaxId() int64 {
	rs, err := tcDao.cache.Incr("current_id").Result()

	if err != nil {
		log.Fatalln("cannot set current_id")
	}

	return rs
}

func (tcDao *TxnCoordinatorDAOImpl) CheckLock(lockId string) (int64, error) {
	result, err := tcDao.cache.Exists(lockId).Result()
	return result, err
}

func (tcDao *TxnCoordinatorDAOImpl) CreateLock(lockId string, state string, timeout time.Duration) (bool, error) {
	updated, err := tcDao.cache.SetNX(lockId, state, timeout).Result()
	if err != nil {
		return false, err
	} else if !updated {
		return false, err
	}
	return true, nil
}

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

func (tcDao *TxnCoordinatorDAOImpl) GetPeersList() map[string]string {
	rs, err := tcDao.cache.HGetAll("peers").Result()

	if err != nil {
		log.Println("cannot get peer list, err: ", err)
		return nil
	}

	return rs
}

func (tcDao *TxnCoordinatorDAOImpl) GetPeerBucket(id string) string {
	rs, err := tcDao.cache.HGet("buckets", id).Result()
	if err != nil {
		log.Println("cannot get peer bucket: ", err)
	}

	return rs
}

func (tcDao *TxnCoordinatorDAOImpl) InsertPeerBucket(id string, peer string) error {
	_, err := tcDao.cache.HSet("buckets", id, peer).Result()
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
	_, err := tcDao.cache.HMSet(txnHKey, fields).Result()
	tcDao.cache.Expire(txnHKey, 600000*time.Millisecond)

	if err != nil {
		log.Fatalln("cannot CreateTransactionEntry")
	}
}

func (tcDao *TxnCoordinatorDAOImpl) CreateSubTransactionEntry(txnId string, state string, inst model.Instruction, globalTxnId string) {
	txnHKey := "sub_txn:" + globalTxnId + ":" + txnId
	fields := map[string]interface{}{
		"state": state,
		"inst":  utils.SerializeMessage(inst),
	}
	_, err := tcDao.cache.HMSet(txnHKey, fields).Result()
	tcDao.cache.Expire(txnHKey, 600000*time.Millisecond)

	if err != nil {
		log.Fatalln("cannot CreateSubTransactionEntry ", err)
	}
}
