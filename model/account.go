package model

type AccountInfo struct {
	Id      int
	Number  string
	Balance float64
}

type PaymentInfo struct {
	Id     int
	From   string
	To     string
	Amount float64
}
