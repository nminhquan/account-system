package resource_manager

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/theodesp/blockingQueues"

	"github.com/nminhquan/mas/consensus"
	"github.com/nminhquan/mas/dao"
	"github.com/nminhquan/mas/db"
	"github.com/nminhquan/mas/model"
	"github.com/nminhquan/mas/utils"

	"go.etcd.io/etcd/etcdserver/api/snap"
)

const DEFAULT_BALANCE float64 = 0.0

type ConsensusService interface {
}

type AccountService interface {
	Start()
	GetAccount(string) (*model.AccountInfo, error)
	CreateAccount(string, float64) (bool, error)
	ProcessSendPayment(string, string, float64) (bool, error)
	ProcessReceivePayment(string, string, float64) (bool, error)
	ProcessRollbackPayment(accountNum string, amount float64) (bool, error)
	Propose(interface{})                      //propose to RaftNode
	ReadCommits(<-chan *string, <-chan error) //read commits from RaftNode
}

type AccountServiceImpl struct {
	commitC     <-chan *string
	proposeC    chan<- string
	mu          sync.RWMutex
	snapshotter *snap.Snapshotter
	resultC     *blockingQueues.BlockingQueue
	errorC      <-chan error
	accDAO      dao.AccountDAO
	// pmtDAO      dao.PaymentDAO
	raftNode *consensus.RaftNode
}

func (accServ *AccountServiceImpl) Start() {
	log.Printf("[AccountService] Waiting")
	gob.Register(model.AccountInfo{})
	gob.Register(model.PaymentInfo{})
	accServ.ReadCommits(accServ.commitC, accServ.errorC)
	go accServ.ReadCommits(accServ.commitC, accServ.errorC)
}

func NewAccountService(accDb *db.RocksDB, commitC <-chan *string, proposeC chan<- string, snapshotter *snap.Snapshotter, errorC <-chan error, raftNode *consensus.RaftNode) *AccountServiceImpl {
	resultC, _ := blockingQueues.NewArrayBlockingQueue(uint64(100000))
	var accDAO dao.AccountDAO = dao.NewAccountDAO(accDb)
	return &AccountServiceImpl{commitC: commitC, proposeC: proposeC, snapshotter: snapshotter, resultC: resultC, errorC: errorC, accDAO: accDAO, raftNode: raftNode}
}
func (accServ *AccountServiceImpl) GetAccount(accountNumber string) (*model.AccountInfo, error) {
	acc, err := accServ.accDAO.GetAccount(accountNumber)
	return acc, err
}

func (accServ *AccountServiceImpl) CreateAccount(accountNumber string, balance float64) (bool, error) {
	start := time.Now()
	accInfo := model.AccountInfo{
		Id:      utils.NewSHAHash(accountNumber),
		Number:  accountNumber,
		Balance: balance,
	}

	ins := model.Instruction{
		Type: model.INS_TYPE_CREATE_ACCOUNT,
		Data: accInfo,
	}
	elapsed := time.Since(start)
	log.Printf("[CreateAccount %v] Time proposed 0: %v", accountNumber, float64(elapsed.Nanoseconds()/int64(time.Millisecond)))
	start1 := time.Now()
	accServ.Propose(ins)

	log.Printf("[CreateAccount %v] Waiting for result message ", accountNumber)
	var result = true
	// select {
	// case rs := <-accServ.resultC:
	// 	result = result && rs
	// }
	elapsed1 := time.Since(start1)
	log.Printf("[CreateAccount %v] Time proposed 1: %v", accountNumber, float64(elapsed1.Nanoseconds()/int64(time.Millisecond)))

	r, _ := accServ.resultC.Get()
	result = result && r.(bool)

	elapsed2 := time.Since(start1)
	log.Printf("[CreateAccount %v] Time proposed 2: %v", accountNumber, float64(elapsed2.Nanoseconds()/int64(time.Millisecond)))
	return result, nil
}

func (accServ *AccountServiceImpl) ProcessSendPayment(fromAcc string, toAcc string, amount float64) (bool, error) {
	// check From balance
	fromaccInfo, _ := accServ.getAccount(fromAcc)
	if fromaccInfo.Balance < amount {
		log.Println("[ProcessSendPayment] not enough balance")
		return false, errors.New(fmt.Sprintf("[ProcessSendPayment] not enough balance, current balance: %v", fromaccInfo.Balance))
	}
	currentT := time.Now().Format(time.RFC850)
	pmInfo := model.PaymentInfo{
		Id:     utils.NewSHAHash(fromAcc, toAcc, currentT),
		From:   fromAcc,
		To:     toAcc,
		Amount: amount,
	}
	ins := model.Instruction{
		Type: model.INS_TYPE_SEND_PAYMENT,
		Data: pmInfo,
	}
	accServ.Propose(ins)

	log.Printf("[ProcessSendPayment] Waiting for result message")
	var result = true
	// select {
	// case rs := <-accServ.resultC:
	// 	result = result && rs
	// }

	r, _ := accServ.resultC.Get()
	result = result && r.(bool)

	return result, nil
}

func (accServ *AccountServiceImpl) ProcessReceivePayment(fromAcc string, toAcc string, amount float64) (bool, error) {
	currentT := time.Now().Format(time.RFC850)
	pmInfo := model.PaymentInfo{
		Id:     utils.NewSHAHash(fromAcc, toAcc, currentT),
		From:   fromAcc,
		To:     toAcc,
		Amount: amount,
	}
	ins := model.Instruction{
		Type: model.INS_TYPE_RECEIVE_PAYMENT,
		Data: pmInfo,
	}
	accServ.Propose(ins)

	log.Printf("[ProcessReceivePayment] Waiting for result message")
	var result = true
	// select {
	// case rs := <-accServ.resultC:
	// 	result = result && rs

	// }
	r, _ := accServ.resultC.Get()
	result = result && r.(bool)
	log.Println("[ProcessReceivePayment], message: ", result)
	return result, nil
}

func (accServ *AccountServiceImpl) ProcessRollbackPayment(accountNum string, amount float64) (bool, error) {
	pmInfo := model.PaymentInfo{
		From:   accountNum,
		To:     accountNum,
		Amount: amount,
	}
	ins := model.Instruction{
		Type: model.INS_TYPE_ROLLBACK,
		Data: pmInfo,
	}
	accServ.Propose(ins)

	log.Printf("[ProcessRollbackPayment] Waiting for result message")
	var result = true
	// select {
	// case rs := <-accServ.resultC:
	// 	result = result && rs

	// }
	r, _ := accServ.resultC.Get()
	result = result && r.(bool)
	log.Println("[ProcessRollbackPayment] ProcessReceivePayment, message: ", result)
	return result, nil
}

func (accServ *AccountServiceImpl) getAccount(accountNumber string) (*model.AccountInfo, error) {
	accInfo, err := accServ.accDAO.GetAccount(accountNumber)
	return accInfo, err
}

// Account service as one part inside the Resource manager will manage data of its raft group, propose to its local raft node channel
// Change then will be redirected to the leader of the node and replicated to other nodes
func (accServ *AccountServiceImpl) Propose(data interface{}) {
	log.Println("[AccountService] Propose data to raft ", data)
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(data); err != nil {
		log.Fatal(err)
	}
	accServ.proposeC <- buf.String()
}

func (accServ *AccountServiceImpl) ReadCommits(commitC <-chan *string, errorC <-chan error) {
	for data := range commitC {
		if data == nil {
			log.Printf("AccountService::getcommitC triggered load snapshot")
			// done replaying log; new data incoming
			// OR signaled to load snapshot
			snapshot, err := accServ.snapshotter.Load()
			if err == snap.ErrNoSnapshot {
				return
			}
			if err != nil {
				log.Panic(err)
			}
			log.Printf("loading snapshot at term %d and index %d", snapshot.Metadata.Term, snapshot.Metadata.Index)
			if err := accServ.recoverFromSnapshot(snapshot.Data); err != nil {
				log.Panic(err)
			}
			continue
		}

		dec := gob.NewDecoder(bytes.NewBufferString(*data))
		var ins model.Instruction
		dec.Decode(&ins)

		go func() {
			// select {
			// case accServ.resultC <- true:
			// }
			accServ.resultC.Put(true)
		}()
		status := accServ.ApplyInstructionToStateMachine(ins)
		log.Println("[Readcommits] ApplyInstructionToStateMachine: ", status)

	}
	if err, ok := <-errorC; ok {
		log.Fatal(err)
	}
}

func (accServ *AccountServiceImpl) ApplyInstructionToStateMachine(ins model.Instruction) bool {
	log.Printf("[ApplyInstructionToStateMachine]: query = %s", ins)

	var status bool
	switch ins.Type {
	case model.INS_TYPE_CREATE_ACCOUNT:
		accInfo := ins.Data.(model.AccountInfo)
		status, _ = accServ.accDAO.CreateAccount(accInfo)

	case model.INS_TYPE_SEND_PAYMENT:
		pmInfo := ins.Data.(model.PaymentInfo)
		// client := client.CreateLockClient(model.LOCK_SERVICE_HOST, pmInfo.From)
		// raftState := accServ.raftNode.Node().Status().RaftState
		// log.Println("[ApplyInstructionToStateMachine] raftState = ", raftState.String())
		status, _ = accServ.accDAO.UpdateAccountBalance(pmInfo.From, (-1)*pmInfo.Amount)

	case model.INS_TYPE_RECEIVE_PAYMENT:
		pmInfo := ins.Data.(model.PaymentInfo)
		// client := client.CreateLockClient(model.LOCK_SERVICE_HOST, pmInfo.To)
		// raftState := accServ.raftNode.Node().Status().RaftState
		// log.Println("[ApplyInstructionToStateMachine] raftState = ", raftState.String())

		status, _ = accServ.accDAO.UpdateAccountBalance(pmInfo.To, pmInfo.Amount)

	case model.INS_TYPE_ROLLBACK:
		pmInfo := ins.Data.(model.PaymentInfo)
		status, _ = accServ.accDAO.UpdateAccountBalance(pmInfo.From, pmInfo.Amount)

	default:
		return false
	}

	return status
}

// TODO func (s *kvstore) getSnapshot() ([]byte, error) {
// 	s.mu.RLock()
// 	defer s.mu.RUnlock()
// 	return json.Marshal(s.kvStore)
// }

func (s *AccountServiceImpl) recoverFromSnapshot(snapshot []byte) error {
	panic("recoverFromSnapshot is not yet impl")
}

func (s *AccountServiceImpl) SetAccDAO(impl dao.AccountDAO) {
	s.accDAO = impl
}
