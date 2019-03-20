package transaction_test

import (
	"fmt"
	"testing"

	"mas/model"
	mytest_mock "mas/test/mock"
	"mas/transaction"

	"github.com/golang/mock/gomock"
	"gotest.tools/assert"
)

func TestMockBasic(t *testing.T) {
	// use the mocked object

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockObj := mytest_mock.NewMockClient(mockCtrl)
	ins := model.Instruction{}
	mockObj.EXPECT().CreatePhase1Request().Return(false)

	ret := mockObj.CreatePaymentRequest("aa", "bb", 10.0)
	fmt.Println(ret)
	assert.Equal(t, ret, false)
}

func TestMockPrepare(t *testing.T) {
	// use the mocked object
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockClientFrom := mytest_mock.NewMockClient(mockCtrl)
	// mockClient.EXPECT().CreatePaymentRequest("aa", "bb", 10.0).Return(false)
	mockClientTo := mytest_mock.NewMockClient(mockCtrl)

	pmInfo := model.PaymentInfo{From: "aa", To: "bb", Amount: 10.0}

	instruction := model.Instruction{}
	txnSwitchFrom := transaction.NewTxnSwitch(pmInfo.From)
	txnSwitchTo := transaction.NewTxnSwitch(pmInfo.To)

	var localTxnFrom transaction.Transaction = transaction.NewLocalTransaction(mockClientFrom, txnSwitchFrom, instruction)
	var localTxnTo transaction.Transaction = transaction.NewLocalTransaction(mockClientTo, txnSwitchTo, instruction)
	txn := transaction.NewGlobalTransaction([]transaction.Transaction{localTxnFrom, localTxnTo})
	prepared := txn.Prepare()
	assert.Equal(t, prepared, true)
}

func TestMockBegin(t *testing.T) {
	// use the mocked object
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// ROLLBACK CASE
	mockClientSender := mytest_mock.NewMockClient(mockCtrl)
	mockClientSender.EXPECT().CreatePaymentRequest("aa", "bb", 10.0).Return(true).MaxTimes(10)
	mockClientReceiver := mytest_mock.NewMockClient(mockCtrl)
	mockClientReceiver.EXPECT().CreatePaymentRequest("bb", "aa", -10.0).Return(false)

	pmInfo := model.PaymentInfo{From: "aa", To: "bb", Amount: 10.0}
	txn := transaction.NewTransaction(mockClientSender, mockClientReceiver, pmInfo)
	phase1Passed := txn.Begin()
	fmt.Println("aaaaa : ", phase1Passed)
	assert.Equal(t, phase1Passed, false)

	if !phase1Passed {
		fmt.Println("Rolling back")
		txn.Rollback()
	}

	// SUCCESS CASE
	mockClientReceiver2 := mytest_mock.NewMockClient(mockCtrl)
	mockClientReceiver2.EXPECT().CreatePaymentRequest("bb", "aa", -10.0).Return(true)
	txn = transaction.NewTransaction(mockClientSender, mockClientReceiver2, pmInfo)
	fmt.Println("Committing")
	phase1Passed = txn.Begin()
	fmt.Println("aaaaa 111s: ", phase1Passed)
	assert.Equal(t, phase1Passed, true)
	txn.Commit()
}
