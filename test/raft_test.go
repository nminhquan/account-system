package mas_test

import (
	"bytes"
	"encoding/gob"
	"flag"
	"log"
	"os"
	"testing"
	"time"

	consensus "github.com/nminhquan/mas/consensus"
	"github.com/nminhquan/mas/model"
	. "github.com/nminhquan/mas/resource_manager"
)

func makeConfig() consensus.RaftClusterConfig {
	mode := flag.String("mode", "server", "servermode/resourcemanagermode")
	cluster := flag.String("cluster", "http://127.0.0.1:9021", "comma separated cluster peers")
	clusterid := flag.Int("clusterid", 1, "node ID")
	id := flag.Int("id", 1, "node ID")
	rpcPort := flag.String("rmPort", ":9121", "rpc Port")
	join := flag.Bool("join", false, "join an existing cluster")
	dbName := flag.String("db", "mas", "db Name")
	flag.Parse()

	clusterConfig := consensus.RaftClusterConfig{cluster, id, rpcPort, join, dbName, clusterid, mode}
	return clusterConfig
}

func StartRaftNode() *consensus.RaftClusterInfo {
	os.RemoveAll("./log/*")
	clusterInfo := consensus.RaftInit(makeConfig())
	log.Println("Raft node created")

	return clusterInfo
}

func CreateAccountService(clusterInfo *consensus.RaftClusterInfo) AccountService {
	log.Println("Raft node created")
	accountServ := NewAccountService(nil, clusterInfo.CommitC, clusterInfo.ProposeC, <-clusterInfo.SnapshotterReady, clusterInfo.ErrorC, nil)
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
