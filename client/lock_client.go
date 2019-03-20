package client

import (
	"context"
	"log"
	pb "mas/proto"
	"time"

	"google.golang.org/grpc"
)

type LockClient interface {
	CreateLockRequest() bool
	CreateReleaseRequest() bool
}

type LockClientImpl struct {
	pb.LockServiceClient
	resourceId string
}

func CreateLockClient(serverAddress string, resourceId string) LockClient {
	conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	c := pb.NewLockServiceClient(conn)

	return &LockClientImpl{c, resourceId}
}

func (client *LockClientImpl) CreateLockRequest() bool {
	log.Println("Creating lock request for acc: ", client.resourceId)
	// clientDeadline := time.Now().Add(time.Duration(10000) * time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	message, err := client.AcquireLock(ctx, &pb.LockRequest{LockId: client.resourceId})

	if err != nil {
		log.Printf("ccould not transfer: %v", err)
		return false
	}
	log.Println("[LockClient] lock returned message: ", message)

	if message.Message != "OK" {
		return false
	}
	return true
}

func (client *LockClientImpl) CreateReleaseRequest() bool {
	log.Println("Creating release request for acc: ", client.resourceId)
	// clientDeadline := time.Now().Add(time.Duration(10000) * time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	message, err := client.ReleaseLock(ctx, &pb.LockRequest{LockId: client.resourceId})
	if err != nil {
		log.Printf("ccould not transfer: %v", err)
		return false
	}
	log.Println("[LockClient] release returned message: ", message)
	if message.Message != "OK" {
		return false
	}
	return true
}
