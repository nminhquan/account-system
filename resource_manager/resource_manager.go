package resource_manager

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"

	"go.etcd.io/etcd/raft/raftpb"

	"gitlab.zalopay.vn/quannm4/mas/model"
	pb "gitlab.zalopay.vn/quannm4/mas/proto"
	"gitlab.zalopay.vn/quannm4/mas/utils"

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

func (server *ResourceManager) ProcessPhase1(ctx context.Context, in *pb.TXRequest) (*pb.TXReply, error) {
	ins := utils.DeserializeMessage(in.Data)
	log.Println("[ProcessPhase1] Decoded instruction = ", ins)
	var message string
	switch ins.Type {
	case model.INS_TYPE_QUERY_ACCOUNT:
		data := ins.Data.(model.AccountInfo)
		message = fmt.Sprintf("%v", *server.accService.GetAccount(data.Number))
	case model.INS_TYPE_CREATE_ACCOUNT:
		data := ins.Data.(model.AccountInfo)
		message = server.accService.CreateAccount(data.Number, data.Balance)
	case model.INS_TYPE_SEND_PAYMENT:
		data := ins.Data.(model.PaymentInfo)
		message = server.accService.ProcessSendPayment(data.From, data.To, data.Amount)
	case model.INS_TYPE_RECEIVE_PAYMENT:
		data := ins.Data.(model.PaymentInfo)
		message = server.accService.ProcessReceivePayment(data.From, data.To, data.Amount)
	}

	log.Println("Phase 1 return message: ", message)
	return &pb.TXReply{Message: message}, nil
}

func (server *ResourceManager) ProcessPhase2Commit(ctx context.Context, in *pb.TXRequest) (*pb.TXReply, error) {
	panic("func (localTxn *LocalTransaction) not yet impl")
}

func (server *ResourceManager) ProcessPhase2Rollback(ctx context.Context, in *pb.TXRequest) (*pb.TXReply, error) {
	ins := utils.DeserializeMessage(in.Data)
	log.Println("[ProcessPhase2Rollback] Decoded instruction = ", ins)

	var message string
	switch ins.Type {
	case model.INS_TYPE_SEND_PAYMENT:
		data := ins.Data.(model.PaymentInfo)
		message = server.accService.ProcessRollbackPayment(data.From, data.Amount)
	case model.INS_TYPE_RECEIVE_PAYMENT:
		data := ins.Data.(model.PaymentInfo)
		message = server.accService.ProcessRollbackPayment(data.To, (-1)*data.Amount)
	}

	return &pb.TXReply{Message: message}, nil
}

func (server *ResourceManager) ReverseSnapshot(ins model.Instruction) string {
	data := ins.Data.(model.AccountInfo)
	message := server.accService.CreateAccount(data.Number, data.Balance)
	return message
}

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
