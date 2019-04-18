package main

import (
	"flag"
	"log"

	"gitlab.zalopay.vn/quannm4/mas/client"

	"gitlab.zalopay.vn/quannm4/mas/credentials"
)

func main() {
	getAcc := flag.Bool("getAcc", false, "createAcc/createPmt")
	createAcc := flag.Bool("createAcc", false, "createAcc/createPmt")
	createPmt := flag.Bool("createPmt", false, "createAcc/createPmt")
	accNum := flag.String("number", "", "Account number")
	accBal := flag.Float64("balance", 0.0, "Account balance")
	accFrom := flag.String("from", "", "From Account number of Payment")
	accTo := flag.String("to", "", "To Account number of Payment")
	amount := flag.Float64("amount", 0.0, "Amount of money for Payment")
	apiToken := flag.String("token", credentials.JWT_TOKEN_FILE, "Path to JWT token file")
	apiCert := flag.String("cert", credentials.SSL_SERVER_CERT, "Path to Server certificate")
	tcServiceHost := flag.String("tc", "localhost:9008", "TC_SERVICE_HOST")

	flag.Parse()
	accClient := client.CreateAccountClient(*tcServiceHost, *apiToken, *apiCert)
	switch {
	case *getAcc:
		if *accNum == "" {
			log.Fatalln("Account number must not be empty")
		}
		accInfo, err := accClient.GetAccountRequest(*accNum)
		if err != nil {
			panic(err)
		}
		log.Println("AccountInfo:", accInfo)
	case *createAcc:
		if *accNum == "" {
			log.Fatalln("Account number must not be empty")
		}
		msg, err := accClient.CreateAccountRequest(*accNum, *accBal)
		if err != nil {
			panic(err)
		}
		log.Println("CreateAccount status:", msg)
	case *createPmt:
		if *accFrom == "" || *accTo == "" {
			log.Fatalln("Account number must not be empty")
		} else if *amount <= 0.0 {
			log.Fatalln("Amount transfer must be > 0.0")
		}
		msg, err := accClient.CreatePaymentRequest(*accFrom, *accTo, *amount)
		if err != nil {
			panic(err)
		}
		log.Println("CreatePayment status:", msg)

		msg, err = accClient.GetAccountRequest(*accFrom)
		if err != nil {
			panic(err)
		}
		log.Println("From Account status:", msg)

		msg, err = accClient.GetAccountRequest(*accTo)
		if err != nil {
			panic(err)
		}
		log.Println("To Account status:", msg)
	default:
		log.Fatalln("invalid action")
	}
}
