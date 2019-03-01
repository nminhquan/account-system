package db

import (
	"fmt"
	"log"
	. "mas/model"
)

type AccountDB struct {
	*MasDB
}

func CreateAccountDB(host string, userName string, password string, dbName string) *AccountDB {
	masDB := CreateSQLite3DB(host, userName, password, dbName)
	return &AccountDB{masDB}
}

func (accDB *AccountDB) InsertAccountInfoToDB(accInfo AccountInfo) int64 {
	log.Printf("AccountDB::InsertAccountInfoToDB")
	db := accDB.DB()
	var accId int64
	rows, err := db.Exec("insert into account(accountNumber, accountBalance) values (?, ?) ", accInfo.Number, accInfo.Balance)
	if err != nil {
		log.Printf(fmt.Sprintf("Query: %v", err))
	} else {
		accId, err = rows.LastInsertId()
	}
	log.Printf("AccountDB::InsertAccountInfoToDB success return ID = %v", accId)
	return accId
}

func (accDB *AccountDB) GetAccountInfoFromDB(accountNumber string) *AccountInfo {

	return nil
}
