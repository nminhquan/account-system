package model

const (
	INS_TYPE_QUERY_ACCOUNT   = "acc_query"
	INS_TYPE_CREATE_ACCOUNT  = "acc_create"
	INS_TYPE_SEND_PAYMENT    = "pm_send"
	INS_TYPE_RECEIVE_PAYMENT = "pm_receive"
	INS_TYPE_ROLLBACK        = "pm_rollback"
	RPC_MESSAGE_OK           = "OK"
	RPC_MESSAGE_FAIL         = "FAIL"
	TXN_STATE_PENDING        = "PENDING"
	TXN_STATE_COMMITTED      = "COMMITTED"
	TXN_STATE_ABORTED        = "ABORTED"
)
