package model

import (
	"go.etcd.io/etcd/raft"
)

type Instruction struct {
	Type string // Createaccount/createpayment
	Data interface{}
}

type TCTransactionEntry struct {
	txn_id        string
	ts            int
	state         string
	parent_txn_id string
}

type RaftData struct {
	EntryData []byte
	Role      raft.StateType
}

type AccountInfo struct {
	Id      string
	Number  string
	Balance float64
}

type PaymentInfo struct {
	Id     string
	From   string
	To     string
	Amount float64
}
