package resource_manager

import (
	"bytes"
	"encoding/gob"
	"log"
	"mas/client"
	"mas/dao"
	"mas/db"
	"mas/model"
	"mas/utils"
	"sync"
	"time"

	"go.etcd.io/etcd/etcdserver/api/snap"
)

const DEFAULT_BALANCE float64 = 0.0

type ConsensusService interface {
}

type AccountService interface {
	Start()
	CreateAccount(string, float64) string
	ProcessSendPayment(string, string, float64) string
	ProcessReceivePayment(string, string, float64) string
	ProcessRollbackPayment(accountNum string, amount float64) string
	Propose(interface{})                      //propose to RaftNode
	ReadCommits(<-chan *string, <-chan error) //read commits from RaftNode
}

type AccountServiceImpl struct {
	accDb       *db.AccountDB
	commitC     <-chan *string
	proposeC    chan<- string
	mu          sync.RWMutex
	snapshotter *snap.Snapshotter
	resultC     chan string
	errorC      <-chan error
	accDAO      dao.AccountDAO
	pmtDAO      dao.PaymentDAO
}

func (accServ *AccountServiceImpl) Start() {
	log.Printf("AccountService::Waiting for commits")
	gob.Register(model.AccountInfo{})
	gob.Register(model.PaymentInfo{})
	accServ.ReadCommits(accServ.commitC, accServ.errorC)
	go accServ.ReadCommits(accServ.commitC, accServ.errorC)
}

func NewAccountService(accDb *db.AccountDB, commitC <-chan *string, proposeC chan<- string, snapshotter *snap.Snapshotter, errorC <-chan error) *AccountServiceImpl {
	resultC := make(chan string)
	var accDAO dao.AccountDAO = dao.NewAccountDAO(accDb)
	log.Printf("accDAO = %p", accDAO)
	var pmtDAO dao.PaymentDAO = dao.NewPaymentDAO(accDb)
	return &AccountServiceImpl{accDb: accDb, commitC: commitC, proposeC: proposeC, snapshotter: snapshotter, resultC: resultC, errorC: errorC, accDAO: accDAO, pmtDAO: pmtDAO}
}

func (accServ *AccountServiceImpl) CreateAccount(accountNumber string, balance float64) string {
	// currentT := time.Now().Format(time.RFC850)
	accInfo := model.AccountInfo{
		Id:      utils.NewSHAHash(accountNumber),
		Number:  accountNumber,
		Balance: balance,
	}

	ins := model.Instruction{
		Type: model.INS_TYPE_CREATE_ACCOUNT,
		Data: accInfo,
	}

	accServ.Propose(ins)

	log.Printf("Waiting for result message")
	var returnMessage string
	select {
	case rs := <-accServ.resultC:
		returnMessage = rs
	}
	log.Println("CreateAccount return message: ", returnMessage)
	return returnMessage
}

func (accServ *AccountServiceImpl) ProcessSendPayment(fromAcc string, toAcc string, amount float64) string {
	// check From balance
	log.Println("ProcessPayment")
	fromaccInfo := accServ.getAccount(fromAcc)
	log.Println("ProcessPayment, fromAcc: ", fromaccInfo)
	if fromaccInfo.Balance < amount {
		log.Println("not enough balance")
		return model.RPC_MESSAGE_FAIL
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

	log.Printf("Waiting for result message")
	var returnMessage string
	select {
	case rs := <-accServ.resultC:
		returnMessage = rs
	}

	log.Println("ProcessSendPayment, message: ", returnMessage)
	return returnMessage
}

func (accServ *AccountServiceImpl) ProcessReceivePayment(fromAcc string, toAcc string, amount float64) string {
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

	log.Printf("Waiting for result message")
	var returnMessage string
	select {
	case rs := <-accServ.resultC:
		returnMessage = rs

	}
	log.Println("ProcessReceivePayment, message: ", returnMessage)
	return returnMessage
}

func (accServ *AccountServiceImpl) ProcessRollbackPayment(accountNum string, amount float64) string {
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

	log.Printf("Waiting for result message")
	var returnMessage string
	select {
	case rs := <-accServ.resultC:
		returnMessage = rs

	}
	log.Println("ProcessReceivePayment, message: ", returnMessage)
	return returnMessage
}

func (accServ *AccountServiceImpl) getAccount(accountNumber string) *model.AccountInfo {
	log.Printf("accServ.accDAO = %p", accServ.accDAO)
	accInfo := accServ.accDAO.GetAccount(accountNumber)
	return accInfo
}

//Account service as one part inside the Resource manager will manage data of its raft group, propose to its local raft node channel
// Change then will be redirected to the leader of the node and replicated to other nodes
func (accServ *AccountServiceImpl) Propose(data interface{}) {
	log.Println("Propose data to raft ", data)
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(data); err != nil {
		log.Fatal(err)
	}
	accServ.proposeC <- buf.String()
}

func (accServ *AccountServiceImpl) ReadCommits(commitC <-chan *string, errorC <-chan error) {
	log.Printf("AccountService::ReadCommits")
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

		log.Println("Readcommits decoded: ins = ", ins)
		status := accServ.ApplyInstructionToStateMachine(ins)
		log.Println("Readcommits ApplyInstructionToStateMachine: ", status)
		go func() {
			switch {
			case status == true:
				accServ.resultC <- model.RPC_MESSAGE_OK
			case status == false:
				log.Printf("Put FAIL message")
				accServ.resultC <- model.RPC_MESSAGE_FAIL
			}

		}()
	}
	if err, ok := <-errorC; ok {
		log.Fatal(err)
	}
}

func (accServ *AccountServiceImpl) ApplyInstructionToStateMachine(ins model.Instruction) bool {
	log.Printf("AccountDB::ApplyInstructionToStateMachine: query = %s", ins)

	var status bool
	switch ins.Type {
	case model.INS_TYPE_CREATE_ACCOUNT:
		accInfo := ins.Data.(model.AccountInfo)
		status = accServ.accDAO.CreateAccount(accInfo)

	case model.INS_TYPE_SEND_PAYMENT:
		pmInfo := ins.Data.(model.PaymentInfo)
		client := client.CreateLockClient(model.LOCK_SERVICE_HOST, pmInfo.From)
		lockStatus := client.CreateLockRequest()
		if lockStatus {
			status = accServ.accDAO.UpdateAccountBalance(pmInfo.From, (-1)*pmInfo.Amount)
		} else {
			log.Println("cannot get lock, lock timeout")
			status = false
		}
		// lockStatus = client.CreateReleaseRequest()

	case model.INS_TYPE_RECEIVE_PAYMENT:
		pmInfo := ins.Data.(model.PaymentInfo)
		client := client.CreateLockClient(model.LOCK_SERVICE_HOST, pmInfo.To)
		lockStatus := client.CreateLockRequest()
		if lockStatus {
			status = accServ.accDAO.UpdateAccountBalance(pmInfo.To, pmInfo.Amount)
		} else {
			log.Println("cannot get lock, lock timeout")
			status = false
		}
		// lockStatus = client.CreateReleaseRequest()

	case model.INS_TYPE_ROLLBACK:
		pmInfo := ins.Data.(model.PaymentInfo)
		status = accServ.accDAO.UpdateAccountBalance(pmInfo.From, pmInfo.Amount)

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
	// var store map[string]string
	// if err := json.Unmarshal(snapshot, &store); err != nil {
	// 	return err
	// }
	// s.mu.Lock()
	// defer s.mu.Unlock()
	// s.kvStore = store
	// return nil
	panic("recoverFromSnapshot is not yet impl")
}

func (s *AccountServiceImpl) SetAccDAO(impl dao.AccountDAO) {
	s.accDAO = impl
}
