package client

import (
	"bytes"
	"context"
	"encoding/gob"
	"log"
	"time"

	"gitlab.zalopay.vn/quannm4/mas/model"
	pb "gitlab.zalopay.vn/quannm4/mas/proto"
	"gitlab.zalopay.vn/quannm4/mas/utils"

	"google.golang.org/grpc"
)

type Client interface {
	CreateGetAccountRequest(ins model.Instruction) string
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

func (client *RMClient) CreateGetAccountRequest(ins model.Instruction) string {
	conn, err := grpc.Dial(client.serverAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewTransactionServiceClient(conn)

	// clientDeadline := time.Now().Add(time.Duration(10000) * time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*6000000)

	defer cancel()
	byteArr := utils.SerializeMessage(ins)
	var result *pb.TXReply

	result, err = c.ProcessPhase1(ctx, &pb.TXRequest{Data: byteArr})

	// //For simple, At RM site only have commit phase, auto-commit is ON
	// result, err = c.Commit(ctx, &pb.TXRequest{Data: network.Bytes()})

	if err != nil {
		log.Printf("ccould not transfer: %v", err)
		return err.Error()
	}
	log.Println("[RMClient] CreateGetAccountRequest Returned msg= ", result.Message)
	return result.Message
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
