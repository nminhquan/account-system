package consensus

import (
	"go.etcd.io/etcd/etcdserver/api/snap"
	"go.etcd.io/etcd/raft/raftpb"
	"strings"
)

type RaftClusterConfig struct {
	Cluster *string
	Id      *int
	RPCPort *string
	Join    *bool
	DBName  *string
}

type RaftClusterInfo struct {
	Node             *RaftNode
	ProposeC         chan string
	CommitC          <-chan *string
	ConfChangeC      chan raftpb.ConfChange
	ErrorC           <-chan error
	SnapshotterReady <-chan *snap.Snapshotter
}

// TODO: Write TESTABLE functions
func RaftInit(clusterConfig RaftClusterConfig) *RaftClusterInfo {
	proposeC := make(chan string)
	confChangeC := make(chan raftpb.ConfChange)
	commitC, errorC, snapshotterReady, rc := NewRaftNode(*clusterConfig.Id, strings.Split(*clusterConfig.Cluster, ","), *clusterConfig.Join, proposeC, confChangeC)
	return &RaftClusterInfo{rc, proposeC, commitC, confChangeC, errorC, snapshotterReady}
}
