package transaction_test

import (
	"testing"

	"github.com/nminhquan/mas/db"
	"github.com/nminhquan/mas/model"
	"github.com/nminhquan/mas/test/mock/mock_client"
	"github.com/nminhquan/mas/transaction"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const TXN_GLOBAL_TEST_ID = "TXN_GLOBAL_TEST_ID"
const TXN_LOCAL_TEST_ID = "TXN_GLOBAL_TEST_ID"

var redisClient = db.NewCacheService("localhost", "")

func TestLocalTXNBDD(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "local transaction test")
}

var _ = Describe("Local transaction test", func() {
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
		Expect(redisClient.Del("sub_txn:" + TXN_GLOBAL_TEST_ID + TXN_LOCAL_TEST_ID).Err()).NotTo(HaveOccurred())
	})

	Specify("Begin() should fail when CreatePhase1Request fails", func() {
		// use the mocked object
		mockClient := mock_client.NewMockClient(mockCtrl)
		ins := model.Instruction{}
		mockClient.EXPECT().CreatePhase1Request(ins).Return(false)

		localTxn := transaction.NewLocalTransaction(mockClient, nil, ins, TXN_LOCAL_TEST_ID, TXN_GLOBAL_TEST_ID)
		ret := localTxn.Begin()
		Expect(ret).To(BeFalse())
	})
	Specify("Begin() should succeed when CreatePhase1Request succeed", func() {
		// use the mocked object
		mockClient := mock_client.NewMockClient(mockCtrl)
		ins := model.Instruction{}
		mockClient.EXPECT().CreatePhase1Request(ins).Return(true)

		localTxn := transaction.NewLocalTransaction(mockClient, nil, ins, TXN_LOCAL_TEST_ID, TXN_GLOBAL_TEST_ID)
		ret := localTxn.Begin()
		Expect(ret).To(BeTrue())
	})

})
