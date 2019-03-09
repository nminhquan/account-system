package dao

import (
	. "mas/db"
	. "mas/model"
)

type AccountDAO interface {
	GetAccount(accountNumber string) *AccountInfo
	CreateAccount(accInfo AccountInfo) int64
	DeductMoney(accountNumber string, amount float64) bool
	DepositMoney(accountNumber string, amount float64) bool
}

type AccountDAOImpl struct {
	db *AccountDB
}

func CreateAccountDAO(db *AccountDB) *AccountDAOImpl {
	return &AccountDAOImpl{db}
}

func (accDao *AccountDAOImpl) GetAccount(accountNumber string) *AccountInfo {
	return accDao.db.GetAccountInfoFromDB(accountNumber)
}

func (accDao *AccountDAOImpl) CreateAccount(accInfo AccountInfo) int64 {
	return accDao.db.InsertAccountInfoToDB(accInfo)
}

func (accDao *AccountDAOImpl) DeductMoney(accountNumber string, amount float64) bool {
	panic("not yet impl")
}

func (accDao *AccountDAOImpl) DepositMoney(accountNumber string, amount float64) bool {
	panic("not yet impl")
}
