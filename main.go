package main //import abc

import (
	"flag"
	"fmt"
	"log"
	"strconv"

	"github.com/nminhquan/mas/config"
	"github.com/nminhquan/mas/node"
	"github.com/nminhquan/mas/transaction"
	"go.etcd.io/etcd/raft/raftpb"

	"github.com/nminhquan/mas/client"

	"github.com/nminhquan/mas/consensus"
	"github.com/nminhquan/mas/db"
	rm "github.com/nminhquan/mas/resource_manager"

	"github.com/nminhquan/mas/utils"
)

func makeConfig() consensus.RaftClusterConfig {
	mode := flag.String("mode", "tc", "servermode/resourcemanagermode")
	cluster := flag.String("cluster", "http://127.0.0.1:9021", "comma separated cluster peers")
	clusterid := flag.Int("clusterid", 1, "node ID")
	id := flag.Int("id", 1, "node ID")
	rmPort := flag.String("rmport", ":9121", "rm Port")
	raftPort := flag.String("raftport", ":9122", "raft Port")
	join := flag.Bool("join", false, "join an existing cluster")
	dbName := flag.String("db", "mas", "db Name")
	flag.Parse()

	clusterConfig := consensus.RaftClusterConfig{cluster, id, raftPort, rmPort, join, dbName, clusterid, mode}
	return clusterConfig
}

// this function will start a new Resource Manager which includes AccountService + RaftNode
// RPC port of RM will be called from Coordinator to write data to the raftgroup
func main() {
	clusterConfig := makeConfig()
	if *clusterConfig.Mode == "tc" {
		log.Println("Running server mode")
		go node.CreateTCNodeManagerServer(config.TCNodeServHost).Start()
		go transaction.CreateTCServer(config.TCServHost).Start()
		go transaction.CreateLockServer(config.LockServHost).Start()
		stopC := make(chan string)
		stopC <- "block"
	} else {
		clusterID := strconv.Itoa(*clusterConfig.IdCluster)
		nodeID := strconv.Itoa(*clusterConfig.Id)
		hostIP := utils.GetHostIPv4()
		rmAddress := hostIP + *clusterConfig.RMPort
		raftAddress := hostIP + *clusterConfig.RaftPort
		nodeClient := client.CreateNodeClient(config.TCNodeServHost, "")
		ok := nodeClient.RegisterNodeToCluster(clusterID, nodeID, rmAddress, raftAddress, *clusterConfig.Join)

		if !ok {
			log.Fatalln("Cannot register new node to raft")
		}

		clusterInfo := consensus.RaftInit(clusterConfig)
		log.Println("Raft node created=================================")

		accountDB, err := db.NewRocksDB(fmt.Sprintf("%v/%v", config.RocksDBDir, *clusterConfig.DBName))
		if err != nil {
			log.Fatalln("Cannot create account db, ", err)
		}
		var accountService rm.AccountService = rm.NewAccountService(accountDB, clusterInfo.CommitC, clusterInfo.ProposeC, <-clusterInfo.SnapshotterReady, clusterInfo.ErrorC, clusterInfo.Node)
		StartServer(accountService, clusterConfig.RMPort, clusterInfo.ConfChangeC)
		log.Println("Raft node is running=================================")
		stopC := make(chan string)
		stopC <- "block"

		clusterInfo.Node.Stop()
		log.Println("Raft node STOPPED=================================")
	}

}

func StartServer(service rm.AccountService, rpcPort *string, confChangeC chan<- raftpb.ConfChange) {
	resourceManager := rm.CreateResourceManager(service, *rpcPort, confChangeC)
	go func() {
		log.Printf("Starting RPC Server: %v \n", *rpcPort)
		resourceManager.StartRPCServer()
	}()
}
