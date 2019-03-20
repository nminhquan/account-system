package main

import (
	"flag"
	"log"
	"mas/consensus"
	"mas/db"
	"mas/model"
	rm "mas/resource_manager"
	"mas/transaction"
)

func makeConfig() consensus.RaftClusterConfig {
	mode := flag.String("mode", "tc", "servermode/resourcemanagermode")
	cluster := flag.String("cluster", "http://127.0.0.1:9021", "comma separated cluster peers")
	clusterid := flag.Int("clusterid", 1, "node ID")
	id := flag.Int("id", 1, "node ID")
	rpcPort := flag.String("port", ":9121", "rpc Port")
	join := flag.Bool("join", false, "join an existing cluster")
	dbName := flag.String("db", "mas", "db Name")
	flag.Parse()

	clusterConfig := consensus.RaftClusterConfig{cluster, id, rpcPort, join, dbName, clusterid, mode}
	return clusterConfig
}

// this function will start a new Resource Manager which includes AccountService + RaftNode
// RPC port of RM will be called from Coordinator to write data to the raftgroup
func main() {
	clusterConfig := makeConfig()
	if *clusterConfig.Mode == "tc" {
		log.Println("Running server mode")
		go transaction.CreateTCServer(model.TC_SERVICE_HOST).Start()
		go transaction.CreateLockServer(model.LOCK_SERVICE_HOST).Start()
		stopC := make(chan string)
		stopC <- "block"
	} else {
		clusterInfo := consensus.RaftInit(clusterConfig)
		log.Println("Raft node createddddddd=================================")

		accountDB := db.CreateAccountDB("localhost", "root", "123456", *clusterConfig.DBName)

		var accountService rm.AccountService = rm.NewAccountService(accountDB, clusterInfo.CommitC, clusterInfo.ProposeC, <-clusterInfo.SnapshotterReady, clusterInfo.ErrorC)
		startServer(accountService, clusterConfig.RPCPort)
		log.Println("Raft node SLEEPING=================================")

		stopC := make(chan string)
		stopC <- "block"

		clusterInfo.Node.Stop()
		log.Println("Raft node STOPPED=================================")
	}

}

func startServer(service rm.AccountService, rpcPort *string) {
	resourceManager := rm.CreateResourceManager(service, *rpcPort)
	go func() {
		log.Printf("Starting RPC Server: %v \n", *rpcPort)
		resourceManager.StartRPCServer()
	}()
}
