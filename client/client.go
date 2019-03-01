package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"google.golang.org/grpc"
	pb "mas/proto"
)

const (
	address                   = "localhost:50051"
	defaultAccountNum         = "0000"
	defaultAccountBal float64 = 0.0
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewAccountServiceClient(conn)

	// Contact the server and print out its response.
	accountNumber := defaultAccountNum
	accountBal := defaultAccountBal
	if len(os.Args) > 1 {
		accountNumber = os.Args[1]
		accountBal, err = strconv.ParseFloat(os.Args[2], 64)
		if err != nil {
			accountBal = float64(defaultAccountBal)
		}
	}
	// clientDeadline := time.Now().Add(time.Duration(10000) * time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)

	defer cancel()

	// this is where gRPC is triggered
	r, err := c.CreateAccount(ctx, &pb.AccountRequest{AccountNumber: accountNumber, Balance: accountBal})
	if err != nil {
		log.Fatalf("ccould not transfer: %v", err)
	}
	log.Printf("Account created info: %d", r.AccountId)
}
