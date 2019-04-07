package transaction

import (
	"log"
	"sync"

	"gitlab.zalopay.vn/quannm4/mas/client"
	"gitlab.zalopay.vn/quannm4/mas/model"
)

type Transaction interface {
	Prepare() bool
	Begin() bool
	Commit()
	Rollback()
}

type GlobalTransaction struct {
	subTxn []Transaction
}

// Create transaction with peers address
// Also contain the data for transferring
func NewGlobalTransaction(subTxn []Transaction) *GlobalTransaction {
	// Create Tx record hold the switch
	return &GlobalTransaction{subTxn}
}

func (tx *GlobalTransaction) Prepare() bool {
	// get snapshot of each record in each sub txn
	// write record to undo_log
	log.Println("[GlobalTXN] Prepare() START")
	var result bool
	for _, sub := range tx.subTxn {
		result = result && sub.Prepare()
	}
	log.Println("[GlobalTXN] Prepare() END")
	return result
}

func (tx *GlobalTransaction) Begin() bool {
	log.Println("[GlobalTXN] Begin() START")
	var resultC = make(chan bool, 10)

	var wg sync.WaitGroup
	for _, localTxn := range tx.subTxn {
		log.Printf("localTxn: %v", localTxn)
		begin := localTxn.Begin
		wg.Add(1)
		go func() {
			result := begin()
			resultC <- result
			defer wg.Done()
		}()
	}

	log.Println("[GlobalTXN] Main: Waiting for sub TXNs to finish")
	wg.Wait()
	var result bool = true
	close(resultC)
	for rs := range resultC {
		result = result && rs
	}

	log.Println("[GlobalTXN] DONE Begin() result: ", result)
	return result
}

func (tx *GlobalTransaction) Commit() {
	// close txn
	// update Txn status
	// delete undo snapshot
	// release globalock
	for _, localTxn := range tx.subTxn {
		go localTxn.Commit()
	}
}

func (tx *GlobalTransaction) Rollback() {
	// call local rollback on failed txn only
	// update Txn status
	// delete undo snapshot
	// release globalock
	for _, localTxn := range tx.subTxn {
		go localTxn.Rollback()

	}
}

/*SINGLE TRANSACTION*/
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
	return true
}
func (localTxn *LocalTransaction) Begin() bool {
	// get snapshot of row before txn
	txnDao.CreateSubTransactionEntry(localTxn.localTxnId, model.TXN_STATE_PENDING, localTxn.data, localTxn.globalTxnId)
	log.Printf("[LocalTXN:%v] Begin() local TXN for rmClient: %v, data: %v", localTxn.localTxnId, localTxn.rmClient, localTxn.data)
	localTxn.globalLock.CreateLockRequest()
	inst := localTxn.data
	success := localTxn.rmClient.CreatePhase1Request(inst)
	localTxn.state = success
	log.Printf("[LocalTXN:%v] Begin() DONE = %v", localTxn.localTxnId, success)
	return success
}

//TODO: Should we use CreatePhase2CommitRequest instead????
func (localTxn *LocalTransaction) Commit() {
	log.Println("[LocalTXN] Commit() local TXN for rmClient: ", localTxn.rmClient, "data: ", localTxn.data)
	txnDao.CreateSubTransactionEntry(localTxn.localTxnId, model.TXN_STATE_COMMITTED, localTxn.data, localTxn.globalTxnId)
	localTxn.globalLock.CreateReleaseRequest()
	log.Println("[LocalTXN] Commit() DONE: ", localTxn.rmClient, "data: ", localTxn.data)
}

func (localTxn *LocalTransaction) Rollback() {
	// create roll back log before rolling back
	// only rollback transaction which was successful
	if localTxn.state {
		log.Printf("[LocalTXN:%v] Rollback() local Rollback", localTxn.localTxnId)
		txnDao.CreateSubTransactionEntry(localTxn.localTxnId, model.TXN_STATE_ABORTED, localTxn.data, localTxn.globalTxnId)
		localTxn.rmClient.CreatePhase2RollbackRequest(localTxn.data)
		log.Printf("[LocalTXN:%v] Rollback() DONE local Rollback", localTxn.localTxnId)
	}

	localTxn.globalLock.CreateReleaseRequest()
}
