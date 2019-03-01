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
	"testing"
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

func TestStartingRaftNode(t *testing.T) {
	consensus.RaftInit(makeConfig())
	log.Println("Raft node created")
}
func TestProposeCWhenPropose(t *testing.T) {
	consensus.RaftInit(makeConfig())
	log.Println("Raft node created")

	accountDB := db.CreateAccountDB("localhost", "root", "123456")
	accountServ := coordinator.CreateAccountService(accountDB, clusterInfo.CommitC, clusterInfo.ProposeC, clusterInfo.SnapshotterReady)
	go accountServ.ReadCommits(clusterInfo.CommitC, clusterInfo.ErrorC)
	account := &model.AccountInfo{1, "1235", 1000.0}

	//Propose
	accountServ.Propose(account)
	time.Sleep(20000 * time.Millisecond)

	//Check commitC if message passes
	// msg := <-clusterInfo.CommitC
	// if msg == nil {
	// 	t.Error("Unexpected message")
	// } else {
	// 	if reicAccObj := decodeMessage(msg); *reicAccObj != *account {
	// 		t.Error("Receive message is not the same as sent")
	// 	} else {
	// 		log.Printf("Received from CommitC: %v", reicAccObj)
	// 	}
	// }
	// log.Printf("msg: %v", msg)
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
