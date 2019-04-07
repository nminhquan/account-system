package node

// import (
// 	"strconv"
// 	"go.etcd.io/etcd/raft/raftpb"
// 	"gitlab.zalopay.vn/quannm4/mas/model"
// 	"context"
// 	"log"
// 	"net"

// 	"gitlab.zalopay.vn/quannm4/mas/credentials"
// 	pb "gitlab.zalopay.vn/quannm4/mas/proto"
// 	"google.golang.org/grpc"
// 	grpc_creds "google.golang.org/grpc/credentials"
// )

// type RMNodeManager struct {
// 	port string
// 	confChangeC chan<- raftpb.ConfChange
// }

// func CreateRMNodeManagerServer(port string, confChangeC chan<- raftpb.ConfChange) *RMNodeManager {
// 	return &RMNodeManager{port, confChangeC}
// }

// func (nm *RMNodeManager) Start() {
// 	lis, err := net.Listen("tcp", nm.port)
// 	if err != nil {
// 		log.Fatalf("[RMNodeManager] failed to listen: %v port : %v", err, nm.port)
// 	}
// 	// register JWTServerInterceptor for authentication
// 	creds, err := grpc_creds.NewServerTLSFromFile(credentials.SSL_SERVER_CERT, credentials.SSL_SERVER_PRIVATE_KEY)
// 	if err != nil {
// 		log.Fatalf("[RMNodeManager] Cannot get credentials: %v", err)
// 	}

// 	s := grpc.NewServer(grpc.UnaryInterceptor(credentials.JWTServerInterceptor), grpc.Creds(creds))
// 	pb.RegisterNodeServiceServer(s, nm)
// 	log.Printf("[RMNodeManager] RPC server started at port: %s", nm.port)
// 	if err := s.Serve(lis); err != nil {
// 		log.Fatalf("[RMNodeManager] failed to serve: %v", err)
// 	}
// }

// func (nm *RMNodeManager) AddNode(ctx context.Context, in *pb.NodeRequest) (*pb.NodeReturn, error) {
// 	log.Println("[RMNodeManager] AddNode for ", in.NodeId, " ", in.NodeAddr)
// 	nodeId, _ := strconv.Atoi(in.NodeId)
// 	cc := raftpb.ConfChange{
// 		Type:    raftpb.ConfChangeAddNode,
// 		NodeID:  uint64(nodeId),
// 		Context: []byte(in.NodeAddr),
// 	}
// 	nm.confChangeC <- cc
// 	return &pb.NodeReturn{Message: model.RPC_MESSAGE_OK}, nil
// }
