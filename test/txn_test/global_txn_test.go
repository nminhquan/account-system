package transaction_test

import (
	"testing"

	"github.com/nminhquan/mas/test/mock/mock_transaction"
	"github.com/nminhquan/mas/transaction"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGlobalTXNBDD(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "transaction test")
}

var _ = Describe("Global test", func() {
	var (
		mockCtrl *gomock.Controller
		// mockTransaction *mytest_mock.MockTransaction
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		// mockTransaction = mytest_mock.NewMockTransaction(mockCtrl)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	It("Global Begin() should succeed if all Locals Transation succeed", func() {
		mockLocalTxn1 := mock_transaction.NewMockTransaction(mockCtrl)
		mockLocalTxn1.EXPECT().Begin().Return(true)
		mockLocalTxn2 := mock_transaction.NewMockTransaction(mockCtrl)
		mockLocalTxn2.EXPECT().Begin().Return(true)
		mockLocalTxn3 := mock_transaction.NewMockTransaction(mockCtrl)
		mockLocalTxn3.EXPECT().Begin().Return(true)
		mockLocalTxn4 := mock_transaction.NewMockTransaction(mockCtrl)
		mockLocalTxn4.EXPECT().Begin().Return(true)

		subTxns := []transaction.Transaction{
			mockLocalTxn1,
			mockLocalTxn2,
			mockLocalTxn3,
			mockLocalTxn4,
		}
		var txn transaction.Transaction = transaction.NewGlobalTransaction(subTxns)
		ret := txn.Begin()

		Expect(ret).To(BeTrue())
	})
	It("Global Begin() should fail with one of Locals Transation fails", func() {
		mockLocalTxn1 := mock_transaction.NewMockTransaction(mockCtrl)
		mockLocalTxn1.EXPECT().Begin().Return(true)
		mockLocalTxn2 := mock_transaction.NewMockTransaction(mockCtrl)
		mockLocalTxn2.EXPECT().Begin().Return(true)
		mockLocalTxn3 := mock_transaction.NewMockTransaction(mockCtrl)
		mockLocalTxn3.EXPECT().Begin().Return(false)
		mockLocalTxn4 := mock_transaction.NewMockTransaction(mockCtrl)
		mockLocalTxn4.EXPECT().Begin().Return(true)

		subTxns := []transaction.Transaction{
			mockLocalTxn1,
			mockLocalTxn2,
			mockLocalTxn3,
			mockLocalTxn4,
		}
		var txn transaction.Transaction = transaction.NewGlobalTransaction(subTxns)
		ret := txn.Begin()
		Expect(ret).To(BeFalse())
	})

})
