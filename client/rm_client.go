package client

import (
	"bytes"
	"context"
	"encoding/gob"
	"log"
	"time"

	"github.com/nminhquan/mas/model"
	pb "github.com/nminhquan/mas/proto"
	"github.com/nminhquan/mas/utils"

	"google.golang.org/grpc"
)

var grpcConnMap *ConnMap = NewConnMap()

type Client interface {
	CreateGetAccountRequest(ins model.Instruction) (string, error)
	CreatePhase1Request(instruction model.Instruction) bool
	CreatePhase2CommitRequest(instruction model.Instruction) bool
	CreatePhase2RollbackRequest(data model.Instruction) bool
}
type RMClient struct {
	serverAddress []string
}

func CreateRMClient(serverAddress []string) *RMClient {
	return &RMClient{serverAddress}
}

func (client *RMClient) getConnection(address string) *grpc.ClientConn {
	conn, err := grpcConnMap.Get(address)
	if err != nil || conn == nil {
		log.Fatalf("Cannot get connection: %v", err)
	}

	return conn
}

func (client *RMClient) CreateGetAccountRequest(ins model.Instruction) (string, error) {
	c := pb.NewTransactionServiceClient(client.getConnection(client.serverAddress[0]))

	// clientDeadline := time.Now().Add(time.Duration(10000) * time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)

	defer cancel()
	byteArr := utils.SerializeMessage(ins)
	var result *pb.TXReply

	result, err := c.ProcessPhase1(ctx, &pb.TXRequest{Data: byteArr})

	if err != nil {
		log.Printf("[CreateGetAccountRequest] Could not transfer: %v", err)
		return "Cannot get account", err
	}

	return result.Message, nil
}

func (client *RMClient) CreatePhase1Request(ins model.Instruction) bool {
	c := pb.NewTransactionServiceClient(client.getConnection(client.serverAddress[0]))

	// clientDeadline := time.Now().Add(time.Duration(10000) * time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()
	byteArr := utils.SerializeMessage(ins)
	var result *pb.TXReply
	start2 := time.Now()
	result, err := c.ProcessPhase1(ctx, &pb.TXRequest{Data: byteArr})

	if err != nil {
		log.Printf("[CreatePhase1Request] ccould not transfer: %v", err)
		return false
	}

	// wait for PaymentReply
	if result.Message != "OK" {
		return false
	}
	elapsed2 := time.Since(start2)
	log.Printf("[CreatePhase1Request %v] Time elapsed RM 1: %v", ins.XID, float64(elapsed2.Nanoseconds()/int64(time.Millisecond)))

	return true
}

func (client *RMClient) CreatePhase2CommitRequest(ins model.Instruction) bool {
	c := pb.NewTransactionServiceClient(client.getConnection(client.serverAddress[0]))

	// clientDeadline := time.Now().Add(time.Duration(10000) * time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)

	defer cancel()
	var network bytes.Buffer // Stand-in for a network connection
	enc := gob.NewEncoder(&network)
	enc.Encode(ins)

	var result *pb.TXReply

	result, err := c.ProcessPhase2Commit(ctx, &pb.TXRequest{Data: network.Bytes()})

	if err != nil {
		log.Printf("[CreatePhase2CommitRequest] ccould not transfer: %v", err)
		return false
	}

	// wait for PaymentReply
	if result.Message != "OK" {
		return false
	}
	return true
}

func (client *RMClient) CreatePhase2RollbackRequest(data model.Instruction) bool {
	c := pb.NewTransactionServiceClient(client.getConnection(client.serverAddress[0]))

	// clientDeadline := time.Now().Add(time.Duration(10000) * time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)

	defer cancel()
	var network bytes.Buffer // Stand-in for a network connection
	enc := gob.NewEncoder(&network)
	enc.Encode(data)

	var result *pb.TXReply

	result, err := c.ProcessPhase2Rollback(ctx, &pb.TXRequest{Data: network.Bytes()})

	if err != nil {
		log.Printf("[CreatePhase2RollbackRequest] ccould not transfer: %v", err)
		return false
	}

	// wait for PaymentReply
	if result.Message != "OK" {
		return false
	}

	return true
}
