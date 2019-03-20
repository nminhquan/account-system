package mytest

import (
	"fmt"
	"mas/coordinator"
	"mas/db"
	"mas/model"
	mytest_mock "mas/test/mock"
	"testing"

	"gotest.tools/assert"

	"github.com/golang/mock/gomock"
)

func TestCreateAccount(t *testing.T) {
	// var wg sync.WaitGroup
	mockCtrl := gomock.NewController(t)
	commitC := make(chan *string)
	proposeC := make(chan string)
	defer mockCtrl.Finish()

	go func() {
		t := "trigger"
		commitC <- &t
	}()

	go func() {
		for {
			c := <-proposeC
			commitC <- &c
		}
	}()

	var accService = coordinator.NewAccountService(&db.AccountDB{}, commitC, proposeC, nil, nil)
	go accService.Start()
	result := accService.CreateAccount("aa", 0)
	assert.Assert(t, result == "OK")
}

func TestProccessingPayment(t *testing.T) {
	// var wg sync.WaitGroup
	mockCtrl := gomock.NewController(t)
	commitC := make(chan *string)
	proposeC := make(chan string)
	defer mockCtrl.Finish()

	// ROLLBACK CASE
	// mockAccService := mytest_mock.NewMockAccountService(mockCtrl)
	accInfo := model.AccountInfo{Number: "aa", Balance: 10000}

	go func() {
		t := "trigger"
		commitC <- &t
	}()

	go func() {
		for {
			c := <-proposeC
			commitC <- &c
		}
	}()
	var accService = coordinator.NewAccountService(&db.AccountDB{}, commitC, proposeC, nil, nil)
	go accService.Start()
	mockDao := mytest_mock.NewMockAccountDAO(mockCtrl)
	mockDao.EXPECT().GetAccount("aa").Return(&accInfo).MaxTimes(10)

	fmt.Println(mockDao.GetAccount("aa"))
	accService.SetAccDAO(mockDao)
	result := accService.ProcessPayment("aa", "bb", 100)
	assert.Equal(t, result, "OK")

	result2 := accService.ProcessPayment("aa", "bb", 100)
	assert.Assert(t, result2 == "OK")
}
