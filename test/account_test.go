package mytest

import (
	"mas/coordinator"
	"mas/db"
	"sync"
	"testing"

	"go.etcd.io/etcd/etcdserver/api/snap"
)

func TestAccountServiceImpl_CreateAccount(t *testing.T) {
	accountDB := db.CreateAccountDB("localhost", "root", "123456", "mas_test")
	type fields struct {
		accDb       *db.AccountDB
		commitC     <-chan *string
		proposeC    chan<- string
		mu          sync.RWMutex
		snapshotter *snap.Snapshotter
		resultC     chan interface{}
		errorC      <-chan error
	}
	type args struct {
		accountNumber string
		balance       float64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		// TODO: Add test cases.
		{
			name:   "test1",
			fields: fields{accDb: accountDB},
			args:   args{"test-account-number", 0.0},
			want:   "account-id-fake",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accServ := coordinator.CreateAccountService(tt.fields.accDb, nil, nil, nil, nil)
			if got := accServ.CreateAccount(tt.args.accountNumber, tt.args.balance); got != tt.want {
				t.Errorf("AccountServiceImpl.CreateAccount() = %v, want %v", got, tt.want)
			}
		})
	}
}
