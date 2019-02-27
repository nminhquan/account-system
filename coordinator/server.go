/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

//go:generate protoc -I ../helloworld --go_out=plugins=grpc:../helloworld ../helloworld/helloworld.proto

package coordinator

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"mas/db"
	_ "mas/db/model"
	pb "mas/proto"
	"net"
)

const (
	port = ":50051"
)

// server is used to implement helloworld.GreeterServer.
type CoordinateServer struct {
	accService AccountService
}

func CreateRPCServer(accDB *db.AccountDB) *CoordinateServer {
	var accountService AccountService = CreateAccountService(accDB)
	log.Println("Service created")
	return &CoordinateServer{accountService}
}

func (s *CoordinateServer) CreateAccount(ctx context.Context, in *pb.AccountRequest) (*pb.AccountReply, error) {
	log.Printf("Create account request for account: %v, balance: %v", in.AccountNumber, in.Balance)
	id := s.createAccount(in.AccountNumber, in.Balance)
	return &pb.AccountReply{AccountId: id}, nil
}

func (s *CoordinateServer) createAccount(accountNumber string, balance float64) int64 {
	return s.accService.CreateAccount(accountNumber, balance)
}

func (server *CoordinateServer) Start() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterAccountServiceServer(s, server)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
