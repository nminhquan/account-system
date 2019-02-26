package main

import "fmt"
import "mas/coordinator"
import "mas/db"
import _ "mas/db/model"

func main() {
	accountDB := db.CreateAccountDB("localhost", "root", "123456")
	defer accountDB.Close()
	var accountService coordinator.AccountService = coordinator.CreateAccountService(accountDB)

	//create account
	var accountId = accountService.CreateAccount("acc_number:1233124124", 0.0)
	fmt.Println(accountId)
}
