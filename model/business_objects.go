package model

type Instruction struct {
	Type string //Can be create account, create payment, update account, SQL query
	Data interface{}
}

type SQLTransaction struct {
	Query  string
	Params interface{}
}

type AccountInfo struct {
	Id      string
	Number  string
	Balance float64
}

type PaymentInfo struct {
	Id     string
	From   string
	To     string
	Amount float64
}
