package main

import (
	"flag"
	"log"
	"mas/consensus"
	"mas/coordinator"
	"mas/db"
)

func makeConfig() consensus.RaftClusterConfig {
	cluster := flag.String("cluster", "http://127.0.0.1:9021", "comma separated cluster peers")
	clusterid := flag.Int("clusterid", 1, "node ID")
	id := flag.Int("id", 1, "node ID")
	rpcPort := flag.String("port", ":9121", "rpc Port")
	join := flag.Bool("join", false, "join an existing cluster")
	dbName := flag.String("db", "mas", "db Name")
	flag.Parse()

	clusterConfig := consensus.RaftClusterConfig{cluster, id, rpcPort, join, dbName, clusterid}
	return clusterConfig
}

// this function will start a new Resource Manager which includes AccountService + RaftNode
// RPC port of RM will be called from Coordinator to write data to the raftgroup
func main() {
	clusterConfig := makeConfig()
	clusterInfo := consensus.RaftInit(clusterConfig)
	log.Println("Raft node createddddddd=================================")

	accountDB := db.CreateAccountDB("localhost", "root", "123456", *clusterConfig.DBName)
	// TODO Handle DB connection clean up
	startServer(accountDB, clusterInfo, clusterConfig.RPCPort)
	log.Println("Raft node SLEEPING=================================")

	stopC := make(chan string)

	stopC <- "block"

	clusterInfo.Node.Stop()
	log.Println("Raft node STOPPED=================================")
}

func startServer(db *db.AccountDB, clusterInfo *consensus.RaftClusterInfo, rpcPort *string) {
	rpcServer := coordinator.CreateRPCServer(db, clusterInfo, *rpcPort)
	go func() {
		log.Printf("Starting RPC Server: %v \n", *rpcPort)
		rpcServer.Start()
	}()
}
