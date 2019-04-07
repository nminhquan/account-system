package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"gitlab.zalopay.vn/quannm4/mas/credentials"

	"gitlab.zalopay.vn/quannm4/mas/client"

	"github.com/myzhan/boomer"
)
var accClient = client.CreateAccountClient("localhost:9008", credentials.JWT_TOKEN_FILE, credentials.SSL_SERVER_CERT)
func queryAccount() {
	start := time.Now()
	time.Sleep(100 * time.Millisecond)
	elapsed := time.Since(start)

	_, err := accClient.GetAccountRequest("1234")
	// log.Println("==================\n=================\nerr", err)
	if err != nil {
		globalBoomer.RecordFailure("grpc", "CreateAccountClient", elapsed.Nanoseconds()/int64(time.Millisecond), "udp error")
	} else {
		globalBoomer.RecordSuccess("grpc", "CreateAccountClient", elapsed.Nanoseconds()/int64(time.Millisecond), int64(10))
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
		Name:   "queryAccount",
		Weight: 10,
		Fn:     queryAccount,
	}

	globalBoomer.Run(task1)

	waitForQuit()
	log.Println("shut down")
}
