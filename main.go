package main

import (
	"flag"
	"fmt"
	"log"
	"mas/consensus"
	"mas/coordinator"
	"mas/db"
	"time"
)

func makeConfig() consensus.RaftClusterConfig {
	cluster := flag.String("cluster", "http://127.0.0.1:9021", "comma separated cluster peers")
	id := flag.Int("id", 1, "node ID")
	dbport := flag.Int("port", 9121, "key-value server port")
	join := flag.Bool("join", false, "join an existing cluster")
	flag.Parse()

	clusterConfig := consensus.RaftClusterConfig{cluster, id, dbport, join}
	return clusterConfig
}

func main() {
	clusterInfo := consensus.RaftInit(makeConfig())
	log.Println("Raft node created")
	time.Sleep(20000 * time.Millisecond)

	accountDB := db.CreateAccountDB("localhost", "root", "123456")
	accountServ := coordinator.CreateAccountService(accountDB, clusterInfo.CommitC, clusterInfo.ProposeC, clusterInfo.SnapshotterReady)
	go accountServ.ReadCommits(clusterInfo.CommitC, clusterInfo.ErrorC)

	startServer(accountDB, clusterInfo)
	fmt.Scanln()
}

func startServer(db *db.AccountDB, clusterInfo *consensus.RaftClusterInfo) {
	rpcServer := coordinator.CreateRPCServer(db, clusterInfo)

	go func() {
		log.Printf("Starting RPC Server: %v \n", 1)
		rpcServer.Start()
	}()
}
