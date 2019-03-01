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

type CoordinateServer struct {
	accService AccountService
	port       string
}

func CreateRPCServer(accDB *db.AccountDB, cluster *RaftClusterInfo, port string) *CoordinateServer {
	var accountService AccountService = CreateAccountService(accDB, cluster.CommitC, cluster.ProposeC, cluster.SnapshotterReady)
	go accountService.Start()
	log.Println("ACCOUNTSERVICE created and started")
	return &CoordinateServer{accountService, port}
}

func (s *CoordinateServer) CreateAccount(ctx context.Context, in *pb.AccountRequest) (*pb.AccountReply, error) {
	log.Printf("Create account request for account: %v, balance: %v", in.AccountNumber, in.Balance)
	id := s.createAccount(in.AccountNumber, in.Balance)
	log.Printf("")
	return &pb.AccountReply{AccountId: id}, nil
}

func (s *CoordinateServer) createAccount(accountNumber string, balance float64) int64 {
	return s.accService.CreateAccount(accountNumber, balance)
}

func (server *CoordinateServer) Start() {
	lis, err := net.Listen("tcp", server.port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterAccountServiceServer(s, server)
	log.Printf("RPC server started at port: %s", server.port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
