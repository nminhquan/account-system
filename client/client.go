package main

import (
	"context"
	"flag"
	"log"
	pb "mas/proto"
	"time"

	"google.golang.org/grpc"
)

const (
	serverAddress             = "localhost:50051"
	defaultAccountNum         = ""
	defaultAction             = "createAcc"
	defaultAccountBal float64 = 0.0
)

func main() {
	address := flag.String("address", serverAddress, "address of coordinator")
	fromAccountNumber := flag.String("fromAcc", defaultAccountNum, "aaa")
	toAccountNumber := flag.String("toAcc", defaultAccountNum, "aaa")
	amount := flag.Float64("amount", defaultAccountBal, "bbb")

	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.Dial(*address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewAccountServiceClient(conn)

	// clientDeadline := time.Now().Add(time.Duration(10000) * time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)

	defer cancel()

	// this is where gRPC is triggered
	// switch *action {
	// case "createAcc":
	// 	r, err := c.CreateAccount(ctx, &pb.AccountRequest{AccountNumber: *toAccountNumber, Balance: *amount})
	// case "createPmt":
	// 	r, err := c.CreatePayment(ctx, &pb.PaymentRequest{FromAccountNumber: *fromAccountNumber, ToAccountNumber: *toAccountNumber, Balance: *amount})
	// }

	if *fromAccountNumber == "" {
		_, err = c.CreateAccount(ctx, &pb.AccountRequest{AccountNumber: *toAccountNumber, Balance: *amount})
	} else {
		_, err = c.CreatePayment(ctx, &pb.PaymentRequest{FromAccountNumber: *fromAccountNumber, ToAccountNumber: *toAccountNumber, Amount: *amount})

	}

	if err != nil {
		log.Fatalf("ccould not transfer: %v", err)
	}
}
