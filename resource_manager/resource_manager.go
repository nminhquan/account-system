package resource_manager

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	"go.etcd.io/etcd/raft/raftpb"

	"github.com/nminhquan/mas/model"
	pb "github.com/nminhquan/mas/proto"
	"github.com/nminhquan/mas/utils"

	"google.golang.org/grpc"
)

// ResourceManager (RM) contains RPCServer for receiving requests from Coordinator
// and Account service for processing with Raft
type ResourceManager struct {
	accService  AccountService
	port        string
	confChangeC chan<- raftpb.ConfChange
}

func CreateResourceManager(accountService AccountService, port string, confChangeC chan<- raftpb.ConfChange) *ResourceManager {
	go accountService.Start()
	log.Println("ACCOUNTSERVICE created and started")
	return &ResourceManager{accountService, port, confChangeC}
}

func (server *ResourceManager) ProcessPhase1(ctx context.Context, in *pb.TXRequest) (*pb.TXReply, error) {
	start := time.Now()
	ins := utils.DeserializeMessage(in.Data)
	log.Println("[RMProcessPhase1] Decoded instruction = ", ins)
	var message string
	var err error
	var status bool
	switch ins.Type {
	case model.INS_TYPE_QUERY_ACCOUNT:
		data := ins.Data.(model.AccountInfo)
		accountInfo, error := server.accService.GetAccount(data.Number)
		err = error
		message = fmt.Sprintf("%v", *accountInfo)
	case model.INS_TYPE_CREATE_ACCOUNT:
		data := ins.Data.(model.AccountInfo)
		status, err = server.accService.CreateAccount(data.Number, data.Balance)
		if status {
			message = model.RPC_MESSAGE_OK
		} else {
			message = model.RPC_MESSAGE_FAIL
		}
	case model.INS_TYPE_SEND_PAYMENT:
		data := ins.Data.(model.PaymentInfo)
		status, err = server.accService.ProcessSendPayment(data.From, data.To, data.Amount)
		if status {
			message = model.RPC_MESSAGE_OK
		} else {
			message = model.RPC_MESSAGE_FAIL
		}
	case model.INS_TYPE_RECEIVE_PAYMENT:
		data := ins.Data.(model.PaymentInfo)
		status, err = server.accService.ProcessReceivePayment(data.From, data.To, data.Amount)
		if status {
			message = model.RPC_MESSAGE_OK
		} else {
			message = model.RPC_MESSAGE_FAIL
		}
	}
	elapsed := time.Since(start)
	log.Printf("[ProcessPhase1 %v] Time proposed: %v", ins.XID, float64(elapsed.Nanoseconds()/int64(time.Millisecond)))
	return &pb.TXReply{Message: message}, err
}

func (server *ResourceManager) ProcessPhase2Commit(ctx context.Context, in *pb.TXRequest) (*pb.TXReply, error) {
	panic("func (localTxn *LocalTransaction) not yet impl")
}

func (server *ResourceManager) ProcessPhase2Rollback(ctx context.Context, in *pb.TXRequest) (*pb.TXReply, error) {
	ins := utils.DeserializeMessage(in.Data)
	log.Println("[RMProcessPhase2Rollback] Decoded instruction = ", ins)

	var message string
	var err error
	var status bool
	switch ins.Type {
	case model.INS_TYPE_SEND_PAYMENT:
		data := ins.Data.(model.PaymentInfo)
		status, err = server.accService.ProcessRollbackPayment(data.From, data.Amount)
		if status {
			message = model.RPC_MESSAGE_OK
		} else {
			message = model.RPC_MESSAGE_FAIL
		}
	case model.INS_TYPE_RECEIVE_PAYMENT:
		data := ins.Data.(model.PaymentInfo)
		status, err = server.accService.ProcessRollbackPayment(data.To, (-1)*data.Amount)
		if status {
			message = model.RPC_MESSAGE_OK
		} else {
			message = model.RPC_MESSAGE_FAIL
		}
	}

	return &pb.TXReply{Message: message}, err
}

// func (server *ResourceManager) ReverseSnapshot(ins model.Instruction) string {
// 	data := ins.Data.(model.AccountInfo)
// 	message, _ := server.accService.CreateAccount(data.Number, data.Balance)
// 	return message
// }

// Send a config change message to Raft when a new node wants to join the raft group
func (rm *ResourceManager) ProposeAddNode(ctx context.Context, in *pb.NodeRequest) (*pb.NodeReturn, error) {
	log.Println("[RMNodeManager] Register AddNode for ", in.NodeId, " ", in.RaftNodeAddr)
	nodeId, _ := strconv.Atoi(in.NodeId)
	cc := raftpb.ConfChange{
		Type:    raftpb.ConfChangeAddNode,
		NodeID:  uint64(nodeId),
		Context: []byte("http://" + in.RaftNodeAddr), // http://172.24.20.67:42379
	}
	rm.confChangeC <- cc
	return &pb.NodeReturn{Message: model.RPC_MESSAGE_OK}, nil
}

func (server *ResourceManager) StartRPCServer() {
	lis, err := net.Listen("tcp", server.port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterTransactionServiceServer(s, server)
	log.Printf("[ResourceManager]RPC server started at port: %s", server.port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
