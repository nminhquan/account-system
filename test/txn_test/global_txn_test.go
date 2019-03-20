package transaction_test

import (
	"fmt"
	mytest_mock "mas/test/mock"
	"mas/transaction"
	"testing"

	"github.com/golang/mock/gomock"
	"gotest.tools/assert"
)

func TestMockBasic(t *testing.T) {
	// use the mocked object

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockLocalTxn1 := mytest_mock.NewMockTransaction(mockCtrl)
	mockLocalTxn1.EXPECT().Prepare().Return(false)
	mockLocalTxn1.EXPECT().Begin().Return(false)
	mockLocalTxn1.EXPECT().Commit().Return(false)
	mockLocalTxn1.EXPECT().Rollback().Return(false)

	mockLocalTxn2 := mytest_mock.NewMockTransaction(mockCtrl)
	mockLocalTxn2.EXPECT().Prepare().Return(false)
	mockLocalTxn2.EXPECT().Begin().Return(false)
	mockLocalTxn2.EXPECT().Commit()
	mockLocalTxn2.EXPECT().Rollback()

	subTxns := []transaction.Transaction{
		mockLocalTxn1,
		mockLocalTxn2,
	}
	var txn transaction.Transaction = transaction.NewGlobalTransaction(subTxns)
	ret := txn.Begin()
	fmt.Println(ret)
	assert.Equal(t, ret, false)
}
