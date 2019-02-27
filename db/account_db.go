package db

import (
	"fmt"
	"go.etcd.io/etcd/etcdserver/api/snap"
	. "mas/db/model"
	"sync"
)

type AccountDB struct {
	*MasDB
	db       *MasDB
	commitC  <-chan *string
	proposeC chan<- string
	mu       sync.RWMutex
	// snapShooter *snap.Snapshotter
}

func CreateAccountDB(host string, userName string, password string, commitC <-chan *string, proposeC chan<- string, snapshotterReady <-chan *snap.Snapshotter) *AccountDB {
	masDB := CreateDB(host, userName, password)
	return &AccountDB{db: masDB, commitC: commitC, proposeC: proposeC}
}

func (accDB *AccountDB) ProposeCommit() {

}

func (accDB *AccountDB) InsertAccountInfoToDB(accInfo AccountInfo) int64 {
	db := accDB.DB()
	var accId int64
	rows, err := db.Exec("insert into account_table(account_number, current_balance) values (?, ?) ", accInfo.Number, accInfo.Balance)
	if err != nil {
		panic(fmt.Sprintf("Query: %v", err))
	} else {
		accId, err = rows.LastInsertId()
	}

	return accId
}

func (accDB *AccountDB) GetAccountInfoFromDB(accountNumber string) *AccountInfo {

	return nil
}
