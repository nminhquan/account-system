package consensus

import (
	"flag"
	"go.etcd.io/etcd/raft/raftpb"
	"log"
	"mas/coordinator"
	"mas/db"
	"strings"
)

func RaftInit() {
	cluster := flag.String("cluster", "http://127.0.0.1:9021", "comma separated cluster peers")
	id := flag.Int("id", 1, "node ID")
	kvport := flag.Int("port", 9121, "key-value server port")
	join := flag.Bool("join", false, "join an existing cluster")
	flag.Parse()

	proposeC := make(chan string)
	defer close(proposeC)
	confChangeC := make(chan raftpb.ConfChange)
	defer close(confChangeC)

	// var kvs *kvstore
	// getSnapshot := func() ([]byte, error) { return kvs.getSnapshot() }

	commitC, errorC, snapshotterReady := NewRaftNode(*id, strings.Split(*cluster, ","), *join, proposeC, confChangeC)

	accountDB := db.CreateAccountDB("localhost", "root", "123456", commitC, proposeC, snapshotterReady)

	StartServer(accountDB, *kvport, confChangeC, errorC)
}

func StartServer(db *db.AccountDB, port int, confChangeC chan raftpb.ConfChange, errorC <-chan error) {
	log.Println("Start server")
	rpcServer := coordinator.CreateRPCServer(db)

	go func() {
		log.Printf("Starting RPC Server: %v \n", 1)
		rpcServer.Start()
	}()
}
