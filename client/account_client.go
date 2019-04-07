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

type AccountClient struct {
	serverAddress string
	apiToken      string
	apiCert       string
}

func CreateAccountClient(serverAddress string, apiToken string, apiCert string) *AccountClient {
	return &AccountClient{serverAddress, apiToken, apiCert}
}

func (client *AccountClient) GetAccountRequest(accountNumber string) (string, error) {
	token, err := credentials.NewTokenFromFile(client.apiToken)
	if err != nil {
		log.Printf("Cannot get token: %v", err)
		return "", err
	}

	creds, err := grpc_creds.NewClientTLSFromFile(client.apiCert, "")
	if err != nil {
		log.Printf("could not load tls cert: %s", err)
		return "", err
	}
	conn, err := grpc.Dial(client.serverAddress,
		grpc.WithTransportCredentials(creds),
		grpc.WithPerRPCCredentials(token))
	if err != nil {
		log.Printf("did not connect: %v", err)
		return "", err
	}

	defer conn.Close()
	c := pb.NewAccountServiceClient(conn)

	// clientDeadline := time.Now().Add(time.Duration(10000) * time.Millisecond)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	repl, err := c.GetAccount(ctx, &pb.AccountRequest{AccountNumber: accountNumber})
	if err != nil {
		log.Printf("Could not transfer: %v", err)
		return "", err
	}

	return repl.Message, nil
}

func (client *AccountClient) CreateAccountRequest(accountNumber string, balance float64) (string, error) {
	token, err := credentials.NewTokenFromFile(client.apiToken)
	if err != nil {
		log.Fatalf("Cannot get token: %v", err)
		return "", err
	}

	creds, err := grpc_creds.NewClientTLSFromFile(client.apiCert, "")
	if err != nil {
		log.Fatalf("could not load tls cert: %s", err)
		return "", err
	}
	conn, err := grpc.Dial(client.serverAddress,
		grpc.WithTransportCredentials(creds),
		grpc.WithPerRPCCredentials(token))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
		return "", err
	}

	defer conn.Close()
	c := pb.NewAccountServiceClient(conn)

	// clientDeadline := time.Now().Add(time.Duration(10000) * time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	repl, err := c.CreateAccount(ctx, &pb.AccountRequest{AccountNumber: accountNumber, Balance: balance})
	if err != nil {
		log.Fatalf("Could not transfer: %v", err)
	}

	return repl.Message, nil
}

func (client *AccountClient) CreatePaymentRequest(fromAcc string, toAcc string, amount float64) (string, error) {
	token, err := credentials.NewTokenFromFile(client.apiToken)
	if err != nil {
		log.Fatalf("Cannot get token: %v", err)
		return "", err
	}

	creds, err := grpc_creds.NewClientTLSFromFile(client.apiCert, "")
	if err != nil {
		log.Fatalf("could not load tls cert: %s", err)
		return "", err
	}
	conn, err := grpc.Dial(client.serverAddress,
		grpc.WithTransportCredentials(creds),
		grpc.WithPerRPCCredentials(token))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
		return "", err
	}
	defer conn.Close()
	c := pb.NewAccountServiceClient(conn)

	// clientDeadline := time.Now().Add(time.Duration(10000) * time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	repl, err := c.CreatePayment(ctx, &pb.PaymentRequest{FromAccountNumber: fromAcc, ToAccountNumber: toAcc, Amount: amount})
	if err != nil {
		log.Fatalf("Could not transfer: %v", err)
	}

	return repl.Message, nil
}
