package mytest

import (
	"bytes"
	"encoding/gob"
	"flag"
	"log"
	consensus "mas/consensus"
	"mas/coordinator"
	"mas/db"
	"mas/model"
	"os"
	"testing"
	"time"
)

func makeConfig() consensus.RaftClusterConfig {
	cluster := flag.String("cluster", "http://127.0.0.1:9021", "comma separated cluster peers")
	id := flag.Int("id", 1, "node ID")
	dbport := flag.String("port", "9121", "rpc port")
	join := flag.Bool("join", false, "join an existing cluster")
	dbName := flag.String("db", "mas", "db Name")
	flag.Parse()

	clusterConfig := consensus.RaftClusterConfig{cluster, id, dbport, join, dbName}
	return clusterConfig
}

func StartRaftNode() *consensus.RaftClusterInfo {
	os.RemoveAll("./log/*")
	clusterInfo := consensus.RaftInit(makeConfig())
	log.Println("Raft node created")

	return clusterInfo
}

func CreateAccountService(clusterInfo *consensus.RaftClusterInfo) coordinator.AccountService {
	log.Println("Raft node created")
	accountDB := db.CreateAccountDB("localhost", "root", "123456", "abc")
	accountServ := coordinator.CreateAccountService(accountDB, clusterInfo.CommitC, clusterInfo.ProposeC, <-clusterInfo.SnapshotterReady, clusterInfo.ErrorC)
	go accountServ.ReadCommits(clusterInfo.CommitC, clusterInfo.ErrorC)
	return accountServ
}
func TestProposeCWhenPropose(t *testing.T) {
	clusterInfo := StartRaftNode()
	accountServ := CreateAccountService(clusterInfo)
	account := &model.AccountInfo{"1", "1235", 1000.0}

	//Propose
	accountServ.Propose(account)
	time.Sleep(20000 * time.Millisecond)
}

func TestCommitCoWhenProposeSuccess(t *testing.T) {
	// go func() {
	// 	for {
	// 		select {
	// 		// case msg := <-clusterInfo.ProposeC:
	// 		// 	log.Printf("Proposed: %v", decodeMessage(&msg))
	// 		case msg := <-clusterInfo.CommitC:
	// 			if msg == nil {
	// 				continue
	// 			}

	// 			log.Printf("Commited msg: %v", *msg)
	// 			// log.Printf("Commited: %v", decodeMessage(msg))
	// 		default:
	// 			// log.Println("waiting proposeC")
	// 		}
	// 	}
	// }()
}

func encodeMessage(data *struct{}) string {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(data); err != nil {
		log.Fatal(err)
	}
	return buf.String()
}

func decodeMessage(data *string) *model.AccountInfo {
	var acc model.AccountInfo
	dec := gob.NewDecoder(bytes.NewBufferString(*data))
	if err := dec.Decode(&acc); err != nil {
		log.Fatalf("raftexample: could not decode message (%v)", err)
	}
	return &acc
}
