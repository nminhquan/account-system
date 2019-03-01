package coordinator

import (
	"bytes"
	"encoding/gob"
	"go.etcd.io/etcd/etcdserver/api/snap"
	"log"
	"mas/db"
	"mas/model"
	"sync"
)

const DEFAULT_BALANCE float64 = 0.0

type AccountService interface {
	Start()
	CreateAccount(string, float64) int64
	Propose(interface{})
	ReadCommits(<-chan *string, <-chan error)
	GetAccount(string) *model.AccountInfo
}

type AccountServiceImpl struct {
	accDb       *db.AccountDB
	commitC     <-chan *string
	proposeC    chan<- string
	mu          sync.RWMutex
	snapshotter *snap.Snapshotter
	resultC     chan interface{}
}

func (accServ *AccountServiceImpl) Start() {
	log.Printf("AccountService::Waiting for commits")
	go accServ.ReadCommits(accServ.commitC, nil)
}

func CreateAccountService(accDb *db.AccountDB, commitC <-chan *string, proposeC chan<- string, snapshotterReady <-chan *snap.Snapshotter) *AccountServiceImpl {
	return &AccountServiceImpl{accDb: accDb, commitC: commitC, proposeC: proposeC}
}

func (accServ *AccountServiceImpl) CreateAccount(accountNumber string, balance float64) int64 {
	accInfo := model.AccountInfo{Number: accountNumber, Balance: DEFAULT_BALANCE}

	accServ.Propose(accInfo)

	log.Printf("Waiting for result message")
	select {
	case accId := <-accServ.resultC:
		log.Printf("Received result message %v", accId)
		iAccId := accId.(int64)
		return iAccId
	}
}

func (accServ *AccountServiceImpl) GetAccount(accountNumber string) *model.AccountInfo {
	accInfo := accServ.accDb.GetAccountInfoFromDB(accountNumber)
	return accInfo
}

func (accServ *AccountServiceImpl) Propose(data interface{}) {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(data); err != nil {
		log.Fatal(err)
	}
	accServ.proposeC <- buf.String()
}

func (accServ *AccountServiceImpl) ReadCommits(commitC <-chan *string, errorC <-chan error) {
	log.Printf("AccountService::ReadCommits")
	for data := range commitC {
		if data == nil {
			// // done replaying log; new data incoming
			// // OR signaled to load snapshot
			// TODO snapshot, err := accServ.snapshotter.Load()
			// if err == snap.ErrNoSnapshot {
			// 	return
			// }
			// if err != nil {
			// 	log.Panic(err)
			// }
			// log.Printf("loading snapshot at term %d and index %d", snapshot.Metadata.Term, snapshot.Metadata.Index)
			// if err := accServ.recoverFromSnapshot(snapshot.Data); err != nil {
			// 	log.Panic(err)
			// }
			// continue
			continue
		}

		var accountData model.AccountInfo
		dec := gob.NewDecoder(bytes.NewBufferString(*data))
		if err := dec.Decode(&accountData); err != nil {
			log.Fatalf("raftexample: could not decode message (%v)", err)
		}
		accServ.mu.Lock()
		log.Printf("AccountService::ReadCommits Apply change to state machine %v", accountData)
		newAccId := accServ.accDb.InsertAccountInfoToDB(accountData)
		accServ.mu.Unlock()
		accServ.resultC <- newAccId
	}
	if err, ok := <-errorC; ok {
		log.Fatal(err)
	}
}

// TODO func (s *kvstore) getSnapshot() ([]byte, error) {
// 	s.mu.RLock()
// 	defer s.mu.RUnlock()
// 	return json.Marshal(s.kvStore)
// }

// func (s *kvstore) recoverFromSnapshot(snapshot []byte) error {
// 	var store map[string]string
// 	if err := json.Unmarshal(snapshot, &store); err != nil {
// 		return err
// 	}
// 	s.mu.Lock()
// 	defer s.mu.Unlock()
// 	s.kvStore = store
// 	return nil
// }
