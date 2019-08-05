package client

import (
	"context"
	"log"
	"time"

	"github.com/nminhquan/mas/credentials"
	pb "github.com/nminhquan/mas/proto"
	"google.golang.org/grpc"
	grpc_creds "google.golang.org/grpc/credentials"
)

type AccountClient struct {
	serverAddress string
	apiToken      string
	apiCert       string
	conn          *grpc.ClientConn
}

func CreateAccountClient(serverAddress string, apiToken string, apiCert string) *AccountClient {
	clientConn, err := createConn(apiCert, apiToken, serverAddress)
	if err != nil {
		log.Fatalf("CreateAccountClient: Failed to create client Connection: %v", err)
	}
	return &AccountClient{serverAddress, apiToken, apiCert, clientConn}
}

func (client *AccountClient) GetAccountRequest(accountNumber string) (string, error) {
	c := pb.NewAccountServiceClient(client.conn)

	// clientDeadline := time.Now().Add(time.Duration(10000) * time.Millisecond)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	repl, err := c.GetAccount(ctx, &pb.AccountRequest{AccountNumber: accountNumber})
	if err != nil {
		log.Fatalf("[AccountClient] Could not transfer: %v", err)
		return err.Error(), err
	}

	return repl.Message, nil
}

func (client *AccountClient) CreateAccountRequest(accountNumber string, balance float64) (string, error) {
	c := pb.NewAccountServiceClient(client.conn)

	// clientDeadline := time.Now().Add(time.Duration(10000) * time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	repl, err := c.CreateAccount(ctx, &pb.AccountRequest{AccountNumber: accountNumber, Balance: balance})
	if err != nil {
		log.Fatalf("Could not transfer: %v", err)
		return "", err
	}
	return repl.Message, nil
}

func (client *AccountClient) CreatePaymentRequest(fromAcc string, toAcc string, amount float64) (string, error) {
	c := pb.NewAccountServiceClient(client.conn)

	// clientDeadline := time.Now().Add(time.Duration(10000) * time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	repl, err := c.CreatePayment(ctx, &pb.PaymentRequest{FromAccountNumber: fromAcc, ToAccountNumber: toAcc, Amount: amount})
	if err != nil {
		log.Fatalf("Could not transfer: %v", err)
		return "", err
	}

	return repl.Message, nil
}

func createConn(apiCert, apiToken, serverAddress string) (*grpc.ClientConn, error) {
	token, err := credentials.NewTokenFromFile(apiToken)
	if err != nil {
		log.Printf("Cannot get token: %v", err)
		return nil, err
	}

	creds, err := grpc_creds.NewClientTLSFromFile(apiCert, "")
	if err != nil {
		log.Printf("could not load tls cert: %s", err)
		return nil, err
	}
	conn, err := grpc.Dial(serverAddress,
		grpc.WithTransportCredentials(creds),
		grpc.WithPerRPCCredentials(token))
	if err != nil {
		log.Printf("did not connect: %v", err)
		return nil, err
	}

	return conn, nil
}

func (client *AccountClient) Close() error {
	return client.conn.Close()
}

func (client *AccountClient) CreatePingPongRequest(message interface{}) (string, error) {
	c := pb.NewAccountServiceClient(client.conn)

	// clientDeadline := time.Now().Add(time.Duration(10000) * time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	rep, err := c.TestMethod(ctx, &pb.TestMessage{Message: message.(string)})
	if err != nil {
		return "", err
	}

	return rep.Message, nil
}
