package dao

import (
	"log"
	. "mas/db"
	"mas/model"
)

type AccountDAO interface {
	GetAccount(accountNumber string) *model.AccountInfo
	CreateAccount(accInfo model.AccountInfo) bool
	UpdateAccountBalance(accountNumber string, amount float64) bool
}

type AccountDAOImpl struct {
	db *AccountDB
}

func NewAccountDAO(db *AccountDB) *AccountDAOImpl {
	return &AccountDAOImpl{db}
}

func (accDao *AccountDAOImpl) GetAccount(accountNumber string) *model.AccountInfo {
	return accDao.db.GetAccountInfoFromDB(accountNumber)
}

func (accDao *AccountDAOImpl) CreateAccount(accInfo model.AccountInfo) bool {
	_, err := accDao.db.InsertAccountInfoToDB(accInfo)
	if err != nil {
		log.Println("[AccountDAO] CreateAccount Error: ", err)
		return false
	}

	return true
}

func (accDao *AccountDAOImpl) UpdateAccountBalance(accountNumber string, amount float64) bool {
	db := accDao.db.DB()

	_, err := db.Query("UPDATE account SET account_balance = account_balance + ? WHERE account_number = ?", amount, accountNumber)
	if err != nil {
		log.Println("[AccountDAO] UpdateAccountBalance Error: ", err)
		return false
	}
	return true
}
