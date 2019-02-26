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
	"log"
	"net"

	"fmt"
	"google.golang.org/grpc"
	pb "grpc-school/txn_example/message"
)

const (
	port = ":50051"
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	balance float32
}

func (s *server) DepositMoney(ctx context.Context, in *pb.DepositRequest) (*pb.DepositReply, error) {
	log.Printf("Deposit request for user: %v, amount: %v", in.UserId, in.Amount)
	s.addMoney(&in.Amount)
	message := fmt.Sprintf("Hello %v, added: %v, current balance: %v", in.UserId, in.Amount, s.balance)
	return &pb.DepositReply{Message: message}, nil
}

func (s *server) addMoney(amount *float32) {
	s.balance += *amount
}

func run() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterMoneyServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
