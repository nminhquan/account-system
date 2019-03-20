package main

import (
	"context"
	"log"
	"mas/model"
	"time"

	pb "mas/proto"

	"google.golang.org/grpc"
)

type MainClient struct {
	serverAddress string
}

func CreateMainClient(serverAddress string) *MainClient {
	return &MainClient{serverAddress}
}

func main() {
	client := CreateMainClient(model.TC_SERVICE_HOST)
	// client.CreateAccountRequest("from_test_1", 2000)
	// client.CreateAccountRequest("to_test_1", 4000)
	client.CreatePaymentRequest("from_test", "to_test", 10000)
}

func (client *MainClient) CreateAccountRequest(accountNumber string, balance float64) bool {
	conn, err := grpc.Dial(client.serverAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewAccountServiceClient(conn)

	// clientDeadline := time.Now().Add(time.Duration(10000) * time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	c.CreateAccount(ctx, &pb.AccountRequest{AccountNumber: accountNumber, Balance: balance})
	if err != nil {
		log.Fatalf("ccould not transfer: %v", err)
	}

	return true
}

func (client *MainClient) CreatePaymentRequest(fromAcc string, toAcc string, amount float64) bool {
	conn, err := grpc.Dial(client.serverAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewAccountServiceClient(conn)

	// clientDeadline := time.Now().Add(time.Duration(10000) * time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	c.CreatePayment(ctx, &pb.PaymentRequest{FromAccountNumber: fromAcc, ToAccountNumber: toAcc, Amount: amount})
	if err != nil {
		log.Fatalf("ccould not transfer: %v", err)
	}

	return true
}
