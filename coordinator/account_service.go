package coordinator

import "mas/db/model"
import "mas/db"

const DEFAULT_BALANCE float64 = 0.0

type AccountService interface {
	CreateAccount(string, float64) int64
	GetAccount(string) *model.AccountInfo
}

type AccountServiceImpl struct {
	accDb *db.AccountDB
}

func CreateAccountService(accDb *db.AccountDB) *AccountServiceImpl {
	return &AccountServiceImpl{accDb}
}
func (acc *AccountServiceImpl) CreateAccount(accountNumber string, balance float64) int64 {
	accInfo := model.AccountInfo{Number: accountNumber, Balance: DEFAULT_BALANCE}
	return acc.accDb.InsertAccountInfoToDB(accInfo)
}

func (acc *AccountServiceImpl) GetAccount(accountNumber string) *model.AccountInfo {
	accInfo := acc.accDb.GetAccountInfoFromDB(accountNumber)
	return accInfo
}
