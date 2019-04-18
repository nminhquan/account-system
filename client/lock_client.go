package client

import (
	"context"
	"log"
	"time"

	pb "gitlab.zalopay.vn/quannm4/mas/proto"
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
	conn, err := grpcConnMap.Get(serverAddress)
	if err != nil {
		log.Fatalf("Cannot get Connection for LockClient: %v", err)
	}
	c := pb.NewLockServiceClient(conn)

	return &LockClientImpl{c, resourceId}
}

func (client *LockClientImpl) CreateLockRequest() bool {
	// clientDeadline := time.Now().Add(time.Duration(10000) * time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	message, err := client.AcquireLock(ctx, &pb.LockRequest{LockId: client.resourceId})

	if err != nil {
		log.Printf("[LockClient] ccould not transfer: %v", err)
		return false
	}

	if message.Message != "OK" {
		return false
	}
	return true
}

func (client *LockClientImpl) CreateReleaseRequest() bool {
	// clientDeadline := time.Now().Add(time.Duration(10000) * time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	message, err := client.ReleaseLock(ctx, &pb.LockRequest{LockId: client.resourceId})
	if err != nil {
		log.Printf("[LockClient] ccould not transfer: %v", err)
		return false
	}

	if message.Message != "OK" {
		return false
	}
	return true
}
