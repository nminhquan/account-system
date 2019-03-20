package client

import (
	"bytes"
	"context"
	"encoding/gob"
	"log"
	"mas/model"
	pb "mas/proto"
	"mas/utils"
	"time"

	"google.golang.org/grpc"
)

type Client interface {
	CreatePhase1Request(instruction model.Instruction) bool
	CreatePhase2CommitRequest(instruction model.Instruction) bool
	CreatePhase2RollbackRequest(data model.Instruction) bool
}
type RMClient struct {
	serverAddress string
}

func CreateRMClient(serverAddress string) *RMClient {
	return &RMClient{serverAddress}
}

func (client *RMClient) CreatePhase1Request(ins model.Instruction) bool {
	conn, err := grpc.Dial(client.serverAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewTransactionServiceClient(conn)

	// clientDeadline := time.Now().Add(time.Duration(10000) * time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)

	defer cancel()
	byteArr := utils.SerializeMessage(ins)
	var result *pb.TXReply

	log.Printf("Sending msg= %v %v to RM: %v", ins.Type, ins.Data, client.serverAddress)
	result, err = c.ProcessPhase1(ctx, &pb.TXRequest{Data: byteArr})

	// //For simple, At RM site only have commit phase, auto-commit is ON
	// result, err = c.Commit(ctx, &pb.TXRequest{Data: network.Bytes()})

	if err != nil {
		log.Printf("ccould not transfer: %v", err)
		return false
	}

	// wait for PaymentReply
	if result.Message != "OK" {
		return false
	}
	log.Println("[RMClient] CreatePhase1Request Returned msg= ", result.Message)
	return true
}

func (client *RMClient) CreatePhase2CommitRequest(ins model.Instruction) bool {
	conn, err := grpc.Dial(client.serverAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewTransactionServiceClient(conn)

	// clientDeadline := time.Now().Add(time.Duration(10000) * time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)

	defer cancel()
	var network bytes.Buffer // Stand-in for a network connection
	enc := gob.NewEncoder(&network)
	enc.Encode(ins)

	var result *pb.TXReply

	result, err = c.ProcessPhase2Commit(ctx, &pb.TXRequest{Data: network.Bytes()})

	if err != nil {
		log.Printf("ccould not transfer: %v", err)
		return false
	}

	// wait for PaymentReply
	if result.Message != "OK" {
		return false
	}
	log.Println("[RMClient] CreatePhase2CommitRequest Returned msg= ", result.Message)
	return true
}

func (client *RMClient) CreatePhase2RollbackRequest(data model.Instruction) bool {
	conn, err := grpc.Dial(client.serverAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewTransactionServiceClient(conn)

	// clientDeadline := time.Now().Add(time.Duration(10000) * time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)

	defer cancel()
	var network bytes.Buffer // Stand-in for a network connection
	enc := gob.NewEncoder(&network)
	enc.Encode(data)

	var result *pb.TXReply

	result, err = c.ProcessPhase2Rollback(ctx, &pb.TXRequest{Data: network.Bytes()})

	if err != nil {
		log.Printf("ccould not transfer: %v", err)
		return false
	}

	// wait for PaymentReply
	if result.Message != "OK" {
		return false
	}
	log.Println("[RMClient] CreatePhase2RollbackRequest Returned msg= ", result.Message)
	return true
}
