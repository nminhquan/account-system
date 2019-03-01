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
	rpcPort := flag.String("port", "9121", "rpc Port")
	join := flag.Bool("join", false, "join an existing cluster")
	dbName := flag.String("db", "mas", "db Name")
	flag.Parse()

	clusterConfig := consensus.RaftClusterConfig{cluster, id, rpcPort, join, dbName}
	return clusterConfig
}

func main() {
	clusterConfig := makeConfig()
	clusterInfo := consensus.RaftInit(clusterConfig)
	log.Println("Raft node createddddddd=================================")

	accountDB := db.CreateAccountDB("localhost", "root", "123456", *clusterConfig.DBName)
	startServer(accountDB, clusterInfo, clusterConfig.RPCPort)
	time.Sleep(100000 * time.Second)
	fmt.Scanln()
}

func startServer(db *db.AccountDB, clusterInfo *consensus.RaftClusterInfo, rpcPort *string) {
	rpcServer := coordinator.CreateRPCServer(db, clusterInfo, *rpcPort)
	go func() {
		log.Printf("Starting RPC Server: %v \n", *rpcPort)
		rpcServer.Start()
	}()
}
