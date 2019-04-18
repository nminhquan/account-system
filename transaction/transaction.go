package transaction

import (
	"log"
	"time"

	"github.com/theodesp/blockingQueues"
	"gitlab.zalopay.vn/quannm4/mas/client"
	"gitlab.zalopay.vn/quannm4/mas/model"
	"gitlab.zalopay.vn/quannm4/mas/utils"
)

type Transaction interface {
	Prepare() bool
	Begin() bool
	Commit()
	Rollback()
}

type GlobalTransaction struct {
	subTxn []Transaction
	txnId  string
}

// Create transaction with peers address
// Also contain the data for transferring
func NewGlobalTransaction(subTxn []Transaction, txnId string) *GlobalTransaction {
	// Create Tx record hold the switch
	return &GlobalTransaction{subTxn, txnId}
}

func (tx *GlobalTransaction) Prepare() bool {
	txnDao.CreateTransactionEntry(tx.txnId, utils.GetCurrentTimeInMillis(), model.TXN_STATE_PENDING, "")
	var result bool = true
	for _, sub := range tx.subTxn {
		result = result && sub.Prepare()
	}
	return result
}
func (tx *GlobalTransaction) Begin() bool {
	start := time.Now()
	var queue, _ = blockingQueues.NewArrayBlockingQueue(uint64(2))
	var result bool = true

	for i, localTxn := range tx.subTxn {
		go func(index int, transaction Transaction) {
			tResult := transaction.Begin()
			queue.Put(tResult)
		}(i, localTxn)
	}
	for _ = range tx.subTxn {
		thisResult, _ := queue.Get()
		result = result && thisResult.(bool)
	}
	elapsed := time.Since(start)
	log.Printf("[Global %v] Time elapsed to Begin: %v", tx.txnId, float64(elapsed.Nanoseconds()/int64(time.Millisecond)))

	return result
}

func (tx *GlobalTransaction) Commit() {
	for _, localTxn := range tx.subTxn {
		go localTxn.Commit()
	}
}

func (tx *GlobalTransaction) Rollback() {
	for _, localTxn := range tx.subTxn {
		go localTxn.Rollback()

	}
}

/*ATOMIC TRANSACTION*/
type LocalTransaction struct {
	rmClient    client.Client
	globalLock  client.LockClient
	data        model.Instruction
	state       bool
	localTxnId  string
	globalTxnId string
}

// Create transaction with peers address
// Also contain the data for transferring
func NewLocalTransaction(rmClient client.Client, globalLock client.LockClient, data model.Instruction, localTxnId string, globalTxnId string) *LocalTransaction {
	// Create Tx record hold the switch
	return &LocalTransaction{rmClient, globalLock, data, true, localTxnId, globalTxnId}
}

func (localTxn *LocalTransaction) Prepare() bool {
	txnDao.CreateSubTransactionEntry(localTxn.localTxnId, model.TXN_STATE_PENDING, localTxn.data, localTxn.globalTxnId)
	return localTxn.globalLock.CreateLockRequest()
}

func (localTxn *LocalTransaction) Begin() bool {
	start := time.Now()
	inst := localTxn.data
	success := localTxn.rmClient.CreatePhase1Request(inst)
	localTxn.state = success
	elapsed := time.Since(start)
	log.Printf("[Local %v] Time elapsed to Begin: %v", localTxn.globalTxnId, float64(elapsed.Nanoseconds()/int64(time.Millisecond)))

	return success
}

//TODO: Should we use CreatePhase2CommitRequest instead????
// close txn
// update Txn status
// delete undo snapshot
// release globalock
func (localTxn *LocalTransaction) Commit() {
	txnDao.CreateSubTransactionEntry(localTxn.localTxnId, model.TXN_STATE_COMMITTED, localTxn.data, localTxn.globalTxnId)
	localTxn.globalLock.CreateReleaseRequest()
}

// call local rollback on failed txn only
// update Txn status
// delete undo snapshot
// release globalock
func (localTxn *LocalTransaction) Rollback() {
	// create roll back log before rolling back
	// only rollback transaction which was committed
	if localTxn.state {
		localTxn.rmClient.CreatePhase2RollbackRequest(localTxn.data)
	}
	txnDao.CreateSubTransactionEntry(localTxn.localTxnId, model.TXN_STATE_ABORTED, localTxn.data, localTxn.globalTxnId)
	localTxn.globalLock.CreateReleaseRequest()
}
