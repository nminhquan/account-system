package main

import (
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/nminhquan/mas/credentials"

	"github.com/nminhquan/mas/client"

	"github.com/myzhan/boomer"
)

const seed = 1000000000

// var redisClient = db.NewCacheService(config.RedisHost, "")
// var rocksDB, _ = db.NewRocksDB("./test/test_db/lock_db")
// var accounts = func() []string {
// 	log.Println("Load all accounts")
// 	read := gorocksdb.NewDefaultReadOptions()
// 	am, _ := rocksDB.DB().MultiGet(read, []byte("buckets:*"))
// 	keys := make([]string, 0, 0)
// 	for _, slice := range am {
// 		keys = append(keys, string(slice.Data()))

// 	}
// 	log.Println("keys = ", keys)
// 	return keys
// }()

var accClient = client.CreateAccountClient("localhost:9008",
	credentials.JWT_TOKEN_FILE,
	credentials.SSL_SERVER_CERT)

func pingPong() {

	start := time.Now()
	_, err := accClient.CreatePingPongRequest("ping")
	elapsed := time.Since(start)

	if err != nil {
		globalBoomer.RecordFailure("/proto.AccountService/TestMethod",
			"pingPong()",
			elapsed.Nanoseconds()/int64(time.Millisecond), err.Error())
	} else {
		globalBoomer.RecordSuccess("/proto.AccountService/TestMethod",
			"pingPong()",
			elapsed.Nanoseconds()/int64(time.Millisecond), int64(10))
	}
}

// func queryAccount() {
// 	start := time.Now()
// 	max := len(accounts)
// 	acc := accounts[rand.Intn(max)]
// 	_, err := accClient.GetAccountRequest(acc)
// 	elapsed := time.Since(start)
// 	// time.Sleep(500 * time.Millisecond)
// 	if err != nil {
// 		globalBoomer.RecordFailure("/proto.AccountService/GetAccount",
// 			"queryAccount()",
// 			elapsed.Nanoseconds()/int64(time.Millisecond), err.Error())
// 		os.Exit(1)
// 	} else {
// 		globalBoomer.RecordSuccess("/proto.AccountService/GetAccount",
// 			"queryAccount()",
// 			elapsed.Nanoseconds()/int64(time.Millisecond), int64(10))
// 	}
// }

// func createPayment1() {
// 	start := time.Now()
// 	max := len(accounts)

// 	fromAcc := accounts[rand.Intn(max)]
// 	toAcc := accounts[rand.Intn(max)]
// 	log.Println("from ", fromAcc, " to ", toAcc)
// 	_, err := accClient.CreatePaymentRequest(fromAcc, toAcc, 1)
// 	elapsed := time.Since(start)
// 	// time.Sleep(500 * time.Millisecond)
// 	if err != nil {
// 		globalBoomer.RecordFailure("/proto.AccountService/CreatePayment",
// 			"createPayment1()",
// 			elapsed.Nanoseconds()/int64(time.Millisecond), err.Error())
// 		os.Exit(1)
// 	} else {
// 		globalBoomer.RecordSuccess("/proto.AccountService/CreatePayment",
// 			"createPayment1()",
// 			elapsed.Nanoseconds()/int64(time.Millisecond), int64(10))
// 	}
// }

func createAccount() {
	start := time.Now()
	accountRandom := rand.Intn(seed)
	_, err := accClient.CreateAccountRequest(strconv.Itoa(accountRandom), 10000)
	elapsed := time.Since(start)
	time.Sleep(500 * time.Millisecond)
	if err != nil {
		globalBoomer.RecordFailure("/proto.AccountService/createAccount",
			"createAccount()",
			elapsed.Nanoseconds()/int64(time.Millisecond), err.Error())
		os.Exit(1)
	} else {
		globalBoomer.RecordSuccess("/proto.AccountService/createAccount",
			"createAccount()",
			elapsed.Nanoseconds()/int64(time.Millisecond), int64(10))
	}
}

func waitForQuit() {
	wg := sync.WaitGroup{}
	wg.Add(1)

	quitByMe := false
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		<-c
		quitByMe = true
		globalBoomer.Quit()
		wg.Done()
	}()

	boomer.Events.Subscribe("boomer:quit", func() {
		if !quitByMe {
			wg.Done()
		}
	})

	wg.Wait()
}

var globalBoomer = boomer.NewBoomer("127.0.0.1", 5557)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	task1 := &boomer.Task{
		Name:   "createAccount",
		Weight: 10,
		Fn:     createAccount,
	}

	globalBoomer.Run(task1)

	waitForQuit()
	log.Println("shut down")
}
