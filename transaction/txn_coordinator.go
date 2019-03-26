package transaction

import (
	"context"
	"fmt"
	"log"
	"mas/client"
	"mas/dao"
	"mas/db"
	"mas/model"
	pb "mas/proto"
	"mas/utils"
	"net"
	"strconv"
	"strings"
	"sync"

	"google.golang.org/grpc"
)

var cacheService = db.NewCacheService("localhost", "")
var txnDao = dao.NewTxnCoordinatorDAO(cacheService)

type TxnCoordinator struct {
	port string
	mtx  sync.Mutex
}

func CreateTCServer(port string) *TxnCoordinator {

	return &TxnCoordinator{port, sync.Mutex{}}
}

func (tc *TxnCoordinator) CreateAccount(ctx context.Context, in *pb.AccountRequest) (*pb.AccountReply, error) {
	accInfo := &model.AccountInfo{Number: in.AccountNumber, Balance: in.Balance}
	nextId := txnDao.IncrMaxId()
	pl := txnDao.GetPeersList()
	log.Println("Peers: ", pl)
	thisPeer := nextId % int64(len(pl))

	log.Println("nextId ", nextId)

	thisPeers := strconv.FormatInt(thisPeer+1, 10)
	log.Println("thisPeer:", thisPeers)
	log.Println("val ", pl[thisPeers])
	peers := strings.Split(pl[thisPeers], ",")

	log.Println("Selected peer: ", peers)
	rmClient := client.CreateRMClient(peers[0])

	globalLock := client.CreateLockClient(model.LOCK_SERVICE_HOST, accInfo.Number)
	instruction := model.Instruction{Type: model.INS_TYPE_CREATE_ACCOUNT, Data: accInfo}
	globalTxnId := utils.GenXid()

	var localTxn = NewLocalTransaction(rmClient, globalLock, instruction, globalTxnId)
	subTxns := []Transaction{
		localTxn,
	}

	var txn Transaction = NewGlobalTransaction(subTxns)

	txn.Begin()
	txnDao.InsertPeerBucket(accInfo.Number, pl[thisPeers])
	return &pb.AccountReply{Message: "Account Added"}, nil
}

func (tc *TxnCoordinator) CreatePayment(ctx context.Context, in *pb.PaymentRequest) (*pb.PaymentReply, error) {
	// check MetaDB
	fromPeerBucket := txnDao.GetPeerBucket(in.FromAccountNumber)
	toPeerBucket := txnDao.GetPeerBucket(in.ToAccountNumber)

	rmClientFrom := client.CreateRMClient(strings.Split(fromPeerBucket, ",")[0])
	rmClientTo := client.CreateRMClient(strings.Split(toPeerBucket, ",")[0])

	log.Printf("rmClientFrom: %v, rmClientTo: %v", rmClientFrom, rmClientTo)

	pmInfo := model.PaymentInfo{From: in.FromAccountNumber, To: in.ToAccountNumber, Amount: in.Amount}

	instructionFrom := model.Instruction{Type: model.INS_TYPE_SEND_PAYMENT, Data: pmInfo}
	instructionTo := model.Instruction{Type: model.INS_TYPE_RECEIVE_PAYMENT, Data: pmInfo}
	globalLockFrom := client.CreateLockClient(model.LOCK_SERVICE_HOST, pmInfo.From)
	globalLockTo := client.CreateLockClient(model.LOCK_SERVICE_HOST, pmInfo.To)

	globalTxnId := utils.GenXid()
	var localTxnFrom = NewLocalTransaction(rmClientFrom, globalLockFrom, instructionFrom, globalTxnId)
	var localTxnTo = NewLocalTransaction(rmClientTo, globalLockTo, instructionTo, globalTxnId)
	subTxns := []Transaction{
		localTxnFrom,
		localTxnTo,
	}
	var txn Transaction = NewGlobalTransaction(subTxns)

	currentTs := utils.GetCurrentTimeInMillis()

	txnDao.CreateTransactionEntry(globalTxnId, currentTs, model.TXN_STATE_PENDING, fmt.Sprintf("%v,%v", pmInfo.From, pmInfo.To))
	//------BEGIN SYNC TRANSACTION/ SERIALIZED
	// tc.mtx.Lock()
	// create transaction object
	// transaction object will use its client to connect to RM servers

	if txn.Begin() {
		txn.Commit()
		txnDao.CreateTransactionEntry(globalTxnId, currentTs, model.TXN_STATE_COMMITTED, fmt.Sprintf("%v,%v", pmInfo.From, pmInfo.To))
	} else {
		log.Println("rollback")
		txnDao.CreateTransactionEntry(globalTxnId, currentTs, model.TXN_STATE_ABORTED, fmt.Sprintf("%v,%v", pmInfo.From, pmInfo.To))
		txn.Rollback()
	}
	// tc.mtx.Unlock()

	return &pb.PaymentReply{Message: "OK"}, nil
}

func (tc *TxnCoordinator) Start() {
	lis, err := net.Listen("tcp", tc.port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterAccountServiceServer(s, tc)
	log.Printf("[TxnCoordinator] RPC server started at port: %s", tc.port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
