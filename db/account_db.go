package db

import (
	"fmt"
	"log"
	"time"

	. "gitlab.zalopay.vn/quannm4/mas/model"
	"gitlab.zalopay.vn/quannm4/mas/utils"
)

type AccountDB struct {
	*MasDB
}

func CreateAccountDB(host string, userName string, password string, dbName string) *AccountDB {
	masDB := CreateMySQLDB(host, userName, password, dbName)
	return &AccountDB{masDB}
}

func (accDB *AccountDB) InsertAccountInfoToDB(accInfo AccountInfo) (bool, error) {
	log.Printf("AccountDB::InsertAccountInfoToDB")
	db := accDB.DB()

	var accId int64
	rows, err := db.Exec("insert into account(account_id, account_number, account_balance) values (?, ?, ?) ", accInfo.Id, accInfo.Number, accInfo.Balance)
	if err != nil {
		log.Printf(fmt.Sprintf("Query: %v", err))
		return false, err
	} else {
		accId, err = rows.LastInsertId()
		log.Printf("AccountDB::InsertAccountInfoToDB success return ID = %v", accId)
		return true, err
	}
}

func (accDB *AccountDB) InsertTransactionToDB(pmInfo PaymentInfo) int64 {
	log.Printf("AccountDB::InsertTransactionToDB")
	db := accDB.DB()

	var accId int64
	timeNow := time.Now().Format(time.RFC850)
	transIdFrom := utils.NewSHAHash(pmInfo.From, pmInfo.To, timeNow)
	transIdTo := utils.NewSHAHash(pmInfo.To, pmInfo.From, timeNow)

	rows, err := db.Exec("insert into transaction(transactionId, type, accountNumber, amount) values (?, ?, ?, ?) ", transIdFrom, "debit", pmInfo.From, pmInfo.Amount)
	_, err = db.Exec("insert into transaction(transactionId, type, accountNumber, amount) values (?, ?, ?, ?) ", transIdTo, "credit", pmInfo.To, pmInfo.Amount)
	if err != nil {
		log.Printf(fmt.Sprintf("Query: %v", err))
	} else {
		accId, err = rows.LastInsertId()
	}
	log.Printf("AccountDB::InsertAccountInfoToDB success return ID = %v", accId)
	return accId
}

func (accDB *AccountDB) InsertPaymentToDB(pmInfo PaymentInfo) string {
	log.Printf("AccountDB::InsertPaymentToDB")
	db := accDB.DB()

	var accId string
	timeNow := time.Now().Format(time.RFC850)
	pmId := utils.NewSHAHash(pmInfo.From, pmInfo.To, timeNow)
	_, err := db.Exec("insert into payment(pm_id, from_acc, to_acc, amount, state) values (?, ?, ?, ?, ?) ", pmId, pmInfo.From, pmInfo.To, pmInfo.Amount, "PENDING")
	if err != nil {
		log.Printf(fmt.Sprintf("Query: %v", err))
	} else {
		accId = pmId
	}
	log.Printf("AccountDB::InsertPaymentToDB success return ID = %v", accId)
	return accId
}

func (accDB *AccountDB) GetAccountInfoFromDB(accountNumber string) *AccountInfo {
	// read from any node, may encounter case where change hasn't been applied to state machine
	// TODO ReadIndex()
	db := accDB.DB()
	var accountInfo AccountInfo
	err := db.QueryRow("select * from account where account_number = ?", accountNumber).Scan(&accountInfo.Id, &accountInfo.Number, &accountInfo.Balance)
	if err != nil {
		log.Printf("Cannot get account")
	}
	return &accountInfo
}

func (accDB *AccountDB) UpdateAccountBalance(accountNumber string, amount float64) bool {
	db := accDB.DB()

	_, err := db.Query("UPDATE account SET account_balance = account_balance + ? WHERE account_number = ?", amount, accountNumber)
	if err != nil {
		log.Println("[AccountDAO] UpdateAccountBalance Error: ", err)
		return false
	}
	return true
}
