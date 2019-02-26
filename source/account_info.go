package source

import "mas/model"

func CreateAccountInfo(account model.AccountInfo) int {
	return account.Id
}

func GetAccountInfo(id int) *model.AccountInfo {
	// connect DB and get account info
	acc := model.AccountInfo{Id: 1, Number: 2}
	return &acc
}
