package transaction

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"

	"gitlab.zalopay.vn/quannm4/mas/config"

	"gitlab.zalopay.vn/quannm4/mas/client"
	"gitlab.zalopay.vn/quannm4/mas/credentials"
	"gitlab.zalopay.vn/quannm4/mas/dao"
	"gitlab.zalopay.vn/quannm4/mas/db"
	"gitlab.zalopay.vn/quannm4/mas/model"
	pb "gitlab.zalopay.vn/quannm4/mas/proto"
	"gitlab.zalopay.vn/quannm4/mas/utils"

	"google.golang.org/grpc"
	grpc_creds "google.golang.org/grpc/credentials"
)

var cacheService = db.NewCacheService(config.RedisHost, "")
var txnDao = dao.NewTxnCoordinatorDAO(cacheService)

type TxnCoordinator struct {
	port string
	mtx  sync.Mutex
}

func CreateTCServer(port string) *TxnCoordinator {

	return &TxnCoordinator{port, sync.Mutex{}}
}

func (tc *TxnCoordinator) GetAccount(ctx context.Context, in *pb.AccountRequest) (*pb.AccountReply, error) {
	// peerBucket := txnDao.GetPeerBucket(in.AccountNumber)
	// rmClient := client.CreateRMClient(strings.Split(peerBucket, ",")[0])
	// accInfo := &model.AccountInfo{Number: in.AccountNumber}
	// instruction := model.Instruction{Type: model.INS_TYPE_QUERY_ACCOUNT, Data: accInfo}
	// message := rmClient.CreateGetAccountRequest(instruction)
	return &pb.AccountReply{Message: "message"}, nil
}

func (tc *TxnCoordinator) CreateAccount(ctx context.Context, in *pb.AccountRequest) (*pb.AccountReply, error) {
	if bucket := txnDao.GetPeerBucket(in.AccountNumber); bucket != "" {
		message := fmt.Sprintln("FAIL: Account already exists, id = ", in.AccountNumber, " bucket = ", bucket)
		log.Println(message)
		return &pb.AccountReply{Message: message}, nil
	}
	accInfo := &model.AccountInfo{Number: in.AccountNumber, Balance: in.Balance}
	nextId := txnDao.IncrMaxId()
	pl := txnDao.GetPeersList()
	log.Println("[CreateAccount] Peers: ", pl)
	thisPeer := nextId % int64(len(pl))
	log.Println("[CreateAccount] nextId ", nextId)
	thisPeers := strconv.FormatInt(thisPeer+1, 10)
	log.Println("[CreateAccount] thisPeers", pl[thisPeers])
	peers := strings.Split(pl[thisPeers], ",")
	log.Println("[[CreateAccount]] Selected peer: ", peers)
	rmClient := client.CreateRMClient(peers[0])

	globalLock := client.CreateLockClient(config.LockServHost, accInfo.Number)
	instruction := model.Instruction{Type: model.INS_TYPE_CREATE_ACCOUNT, Data: accInfo}
	globalTxnId := utils.GenXid()
	localTxnId := utils.GenXid()
	var localTxn = NewLocalTransaction(rmClient, globalLock, instruction, localTxnId, globalTxnId)
	subTxns := []Transaction{
		localTxn,
	}

	var txn Transaction = NewGlobalTransaction(subTxns)
	var message string
	if txn.Begin() {
		txn.Commit()
		err := txnDao.InsertPeerBucket(accInfo.Number, pl[thisPeers])
		if err != nil {
			message = model.RPC_MESSAGE_FAIL
		}
		message = model.RPC_MESSAGE_OK
	} else {
		message = model.RPC_MESSAGE_FAIL
	}

	return &pb.AccountReply{Message: message}, nil
}

func (tc *TxnCoordinator) CreatePayment(ctx context.Context, in *pb.PaymentRequest) (*pb.PaymentReply, error) {
	// check MetaDB
	fromPeerBucket := txnDao.GetPeerBucket(in.FromAccountNumber)
	toPeerBucket := txnDao.GetPeerBucket(in.ToAccountNumber)

	paymentValid := tc.validatePayment(in.FromAccountNumber, in.ToAccountNumber, fromPeerBucket, toPeerBucket)

	if !paymentValid {
		message := "FAIL: Invalid payment information, please check again"
		log.Println(message)

		return &pb.PaymentReply{Message: message}, nil
	}

	rmClientFrom := client.CreateRMClient(strings.Split(fromPeerBucket, ",")[0])
	rmClientTo := client.CreateRMClient(strings.Split(toPeerBucket, ",")[0])
	log.Printf("rmClientFrom: %v, rmClientTo: %v", rmClientFrom, rmClientTo)
	pmInfo := model.PaymentInfo{From: in.FromAccountNumber, To: in.ToAccountNumber, Amount: in.Amount}
	instructionFrom := model.Instruction{Type: model.INS_TYPE_SEND_PAYMENT, Data: pmInfo}
	instructionTo := model.Instruction{Type: model.INS_TYPE_RECEIVE_PAYMENT, Data: pmInfo}
	globalLockFrom := client.CreateLockClient(config.LockServHost, pmInfo.From)
	globalLockTo := client.CreateLockClient(config.LockServHost, pmInfo.To)

	globalTxnId := utils.GenXid()
	var localTxnFrom = NewLocalTransaction(rmClientFrom, globalLockFrom, instructionFrom, utils.GenXid(), globalTxnId)
	var localTxnTo = NewLocalTransaction(rmClientTo, globalLockTo, instructionTo, utils.GenXid(), globalTxnId)
	subTxns := []Transaction{
		localTxnFrom,
		localTxnTo,
	}
	var txn Transaction = NewGlobalTransaction(subTxns)
	currentTs := utils.GetCurrentTimeInMillis()
	txnDao.CreateTransactionEntry(globalTxnId, currentTs, model.TXN_STATE_PENDING, fmt.Sprintf("%v,%v", pmInfo.From, pmInfo.To))

	var message string
	if txn.Begin() {
		txn.Commit()
		txnDao.CreateTransactionEntry(globalTxnId, currentTs, model.TXN_STATE_COMMITTED, fmt.Sprintf("%v,%v", pmInfo.From, pmInfo.To))
		message = model.TXN_STATE_COMMITTED
	} else {
		log.Println("rollback")
		txnDao.CreateTransactionEntry(globalTxnId, currentTs, model.TXN_STATE_ABORTED, fmt.Sprintf("%v,%v", pmInfo.From, pmInfo.To))
		message = model.TXN_STATE_ABORTED
		txn.Rollback()
	}

	return &pb.PaymentReply{Message: message}, nil
}

func (tc *TxnCoordinator) validatePayment(from string, to string, fromPeerBucket string, toPeerBucket string) bool {
	if from == to {
		log.Println("Cannot send money to yourself.")
		return false
	} else if fromPeerBucket == "" {
		log.Println("From account: ", from, " doesn't exist")
		return false
	} else if toPeerBucket == "" {
		log.Println("To account: ", from, " doesn't exist")
		return false
	}

	return true
}

func (tc *TxnCoordinator) Start() {
	lis, err := net.Listen("tcp", tc.port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// register JWTServerInterceptor for authentication
	creds, err := grpc_creds.NewServerTLSFromFile(credentials.SSL_SERVER_CERT, credentials.SSL_SERVER_PRIVATE_KEY)
	if err != nil {
		log.Fatalf("Cannot get credentials: %v", err)
	}

	s := grpc.NewServer(grpc.UnaryInterceptor(credentials.JWTServerInterceptor), grpc.Creds(creds))
	pb.RegisterAccountServiceServer(s, tc)
	log.Printf("[TxnCoordinator] RPC server started at port: %s", tc.port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
