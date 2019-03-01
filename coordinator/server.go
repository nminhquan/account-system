package coordinator

import (
	"context"
	"log"
	. "mas/consensus"
	"mas/db"
	pb "mas/proto"
	"net"

	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

type CoordinateServer struct {
	accService AccountService
}

func CreateRPCServer(accDB *db.AccountDB, cluster *RaftClusterInfo) *CoordinateServer {
	var accountService AccountService = CreateAccountService(accDB, cluster.CommitC, cluster.ProposeC, cluster.SnapshotterReady)
	go accountService.Start()
	log.Println("ACCOUNTSERVICE created and started")
	return &CoordinateServer{accountService}
}

func (s *CoordinateServer) CreateAccount(ctx context.Context, in *pb.AccountRequest) (*pb.AccountReply, error) {
	log.Printf("Create account request for account: %v, balance: %v", in.AccountNumber, in.Balance)
	id := s.createAccount(in.AccountNumber, in.Balance)
	return &pb.AccountReply{AccountId: id}, nil
}

func (s *CoordinateServer) createAccount(accountNumber string, balance float64) int64 {
	return s.accService.CreateAccount(accountNumber, balance)
}

func (server *CoordinateServer) Start() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterAccountServiceServer(s, server)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	log.Printf("RPC server started")
}
