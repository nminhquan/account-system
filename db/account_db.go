package db

import (
	"database/sql"
	"fmt"
	"mas/db/model"
)

type AccountDB struct {
	db *MasDB
}

func CreateAccountDB(host string, userName string, password string) *AccountDB {
	masDB := CreateDB(host, userName, password)
	return &AccountDB{masDB}
}

func (accDB *AccountDB) Close() {
	accDB.db.Close()
}

func (accDB *AccountDB) DB() *sql.DB {
	return accDB.db.DB()
}

func (accDB *AccountDB) InsertAccountInfoToDB(accInfo model.AccountInfo) int64 {
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

func (accDB *AccountDB) GetAccountInfoFromDB(accountNumber string) *model.AccountInfo {

	return nil
}
