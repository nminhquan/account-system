package client

import (
	"context"
	"log"
	"time"

	"gitlab.zalopay.vn/quannm4/mas/credentials"

	pb "gitlab.zalopay.vn/quannm4/mas/proto"
	"google.golang.org/grpc"
	grpc_creds "google.golang.org/grpc/credentials"
)

type NodeClient struct {
	tcAddr     string
	raftLeader string
}

func CreateNodeClient(tcAddr string, raftClusterAddr string) *NodeClient {
	return &NodeClient{tcAddr, raftClusterAddr}
}

func (client *NodeClient) RegisterNodeToCluster(clusterId string, nodeId string, rmAddr string, raftAddr string, join bool) bool {
	log.Println("RegisterHostWithTC: Registering ", clusterId, " ", nodeId, " ", rmAddr, " with ", client.tcAddr)
	token, err := credentials.NewTokenFromFile(credentials.JWT_TOKEN_FILE)
	if err != nil {
		log.Fatalf("Cannot get token: %v", err)
	}

	creds, err := grpc_creds.NewClientTLSFromFile(credentials.SSL_SERVER_CERT, "")
	if err != nil {
		log.Fatalf("could not load tls cert: %s", err)
	}
	conn, err := grpc.Dial(client.tcAddr,
		grpc.WithTransportCredentials(creds),
		grpc.WithPerRPCCredentials(token))

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewNodeServiceClient(conn)

	// clientDeadline := time.Now().Add(time.Duration(10000) * time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)

	defer cancel()

	result, err := c.AddNode(ctx, &pb.NodeRequest{ClusterId: clusterId, NodeId: nodeId, RmNodeAddr: rmAddr, RaftNodeAddr: raftAddr, JoinRaft: join})

	if err != nil {
		log.Printf("RegisterHostWithTC Could not transfer: %v", err)
		return false
	}
	log.Println("RegisterHostWithTC Returned msg= ", result.Message)
	return true
}

func (client *NodeClient) JoinExistingRaft(clusterId string, nodeId string, addr string) bool {
	leaderAdd := client.raftLeader
	log.Println("JoinExistingRaft: Registering ", clusterId, " ", nodeId, " ", addr, " with ", leaderAdd)
	conn, err := grpc.Dial(leaderAdd, grpc.WithInsecure())

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewTransactionServiceClient(conn)

	// clientDeadline := time.Now().Add(time.Duration(10000) * time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)

	defer cancel()

	result, err := c.ProposeAddNode(ctx, &pb.NodeRequest{ClusterId: clusterId, NodeId: nodeId, RaftNodeAddr: addr})

	if err != nil {
		log.Printf("RegisterHostWithTC Could not transfer: %v", err)
		return false
	}
	log.Println("RegisterHostWithTC Returned msg= ", result.Message)
	return true
}
