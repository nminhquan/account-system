package transaction

import (
	"context"
	"log"
	"net"
	"sync"
	"time"

	pb "gitlab.zalopay.vn/quannm4/mas/proto"

	"google.golang.org/grpc"
)

type LockService struct {
	port string
	mtx  sync.Mutex
}

func CreateLockServer(port string) *LockService {

	return &LockService{port, sync.Mutex{}}
}

var lockProperties = &LockProperties{
	LockAttempts:     30,
	LockAttemptDelay: time.Duration(500 * time.Millisecond),
	LockTimeout:      time.Duration(10000 * time.Millisecond),
}

func (ls *LockService) AcquireLock(ctx context.Context, in *pb.LockRequest) (*pb.LockReply, error) {
	lock := NewLockController(in.LockId, lockProperties)

	var message string

	if stt, err := lock.Lock(); !stt || err != nil {
		message = "FAILED"
	} else {
		message = "OK"
	}
	// Acq lock
	// log.Println("[LockService] Acquire lock for ", in.LockId, " message: ", message)
	return &pb.LockReply{Message: message}, nil
}

func (ls *LockService) ReleaseLock(ctx context.Context, in *pb.LockRequest) (*pb.LockReply, error) {
	lock := NewLockController(in.LockId, lockProperties)

	var message string
	// Rls lock
	if stt, err := lock.Unlock(); !stt || err != nil {
		message = "FAILED"
	} else {
		message = "OK"
	}
	// log.Println("[LockService] Release lock for ", in.LockId, " message: ", message)
	return &pb.LockReply{Message: message}, nil
}

func (ls *LockService) Start() {
	lis, err := net.Listen("tcp", ls.port)
	if err != nil {
		log.Fatalf("[LockService] failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterLockServiceServer(s, ls)
	log.Printf("[LockService] RPC server started at port: %s", ls.port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
