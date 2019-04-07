package dao

import (
	"log"
	"strconv"

	. "gitlab.zalopay.vn/quannm4/mas/db"
	"gitlab.zalopay.vn/quannm4/mas/model"
)

type AccountDAO interface {
	GetAccount(accountNumber string) (*model.AccountInfo, error)
	CreateAccount(accInfo model.AccountInfo) (bool, error)
	UpdateAccountBalance(accountNumber string, amount float64) (bool, error)
}

type AccountDAOImpl struct {
	db *RocksDB
}

func NewAccountDAO(db *RocksDB) *AccountDAOImpl {
	return &AccountDAOImpl{db}
}

func (accDao *AccountDAOImpl) GetAccount(accountNumber string) (*model.AccountInfo, error) {
	bal, err := accDao.db.GetAccountBalance(accountNumber)
	if err != nil {
		log.Println("cannot GetAccount")
		return nil, err
	}
	balNumber, _ := strconv.ParseFloat(bal, 64)
	return &model.AccountInfo{Number: accountNumber, Balance: balNumber}, nil
}

func (accDao *AccountDAOImpl) CreateAccount(accInfo model.AccountInfo) (bool, error) {
	err := accDao.db.SetAccountBalance(accInfo.Number, accInfo.Balance)
	if err != nil {
		log.Println("[AccountDAO] CreateAccount Error: ", err)
		return false, err
	}

	return true, nil
}

func (accDao *AccountDAOImpl) UpdateAccountBalance(accountNumber string, amount float64) (bool, error) {
	bal, _ := accDao.db.GetAccountBalance(accountNumber)
	currentBal, _ := strconv.ParseFloat(bal, 64)
	err := accDao.db.SetAccountBalance(accountNumber, (currentBal + amount))
	if err != nil {
		log.Println("[AccountDAO] UpdateAccountBalance Error: ", err)
		return false, err
	}

	return true, nil
}
