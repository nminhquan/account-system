package transaction_test

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/nminhquan/mas/db"
	. "github.com/nminhquan/mas/transaction"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const testRedisKey = "__test_key__"

func TestLockBDD(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Lock test")
}

var redisClient = db.NewCacheService("localhost", "")
var lockProperties = &LockProperties{
	LockAttempts:     3,
	LockAttemptDelay: time.Duration(1000 * time.Millisecond),
	LockTimeout:      time.Duration(4000 * time.Millisecond),
}
var newLock = func() Lock {
	return NewLockController(testRedisKey, lockProperties)
}

var _ = Describe("Global lock", func() {
	var (
		subject Lock
	)
	BeforeEach(func() {
		subject = newLock()
		Expect(subject.IsLocked()).To(BeFalse())
	})
	AfterEach(func() {
		Expect(redisClient.Del(testRedisKey).Err()).NotTo(HaveOccurred())
	})

	It("should get new lock and unlock successfully", func() {
		subject.Lock()
		lockStatus, _ := redisClient.Get(testRedisKey).Result()
		fmt.Println("lockStatus: ", lockStatus)
		Expect(lockStatus).To(Equal(LockDefaultState))
		Expect(subject.IsLocked()).To(BeTrue())
		subject.Unlock()
		lockStatus, _ = redisClient.Get(testRedisKey).Result()
		fmt.Println("lockStatus after : ", lockStatus)
		Expect(lockStatus).To(Equal(""))
		Expect(subject.IsLocked()).To(BeFalse())
	})
	Specify("lock could not take other handle", func() {
		subject.Lock()
		Expect(subject.IsLocked()).To(BeTrue())
		Expect(subject.HasHandle()).To(BeTrue())
		lockCtl := NewLockController(testRedisKey, lockProperties)
		Expect(lockCtl.HasHandle()).To(BeFalse())
		Expect(lockCtl.IsLocked()).To(BeTrue())
	})
	It("should have lock status updated", func() {
		Expect(subject.Lock()).To(BeTrue())
		Expect(subject.IsLocked()).To(BeTrue())
		Expect(subject.Unlock()).To(BeTrue())
		Expect(subject.IsLocked()).To(BeFalse())
	})
	It("should work when using Obtain lock", func() {
		lockCtl, err := Obtain(testRedisKey, lockProperties)
		Expect(lockCtl.IsLocked()).To(BeTrue())
		Expect(err).To(BeNil())
	})
	It("should have error obtaining already-locked lock", func() {
		subject.Lock()
		_, err := Obtain(testRedisKey, lockProperties)
		Expect(err).To(Equal(ErrLockNotObtained))
	})
	It("cannot lock the already-locked resource", func() {
		subject.Lock()
		lockCtl := NewLockController(testRedisKey, lockProperties)
		stt, err := lockCtl.Lock()
		Expect(stt).To(BeFalse())
		Expect(err).To(Equal(ErrLockNotObtained))
	})
	It("will wait for the lock timeout and get the lock successfully", func() {
		lockProperties = &LockProperties{
			LockAttempts:     5,
			LockAttemptDelay: time.Duration(1000 * time.Millisecond),
			LockTimeout:      time.Duration(4000 * time.Millisecond),
		}
		subject = newLock()
		subject.Lock()
		lockCtl := NewLockController(testRedisKey, lockProperties)
		stt, err := lockCtl.Lock()
		Expect(stt).To(BeTrue())
		Expect(err).To(BeNil())
	})
	It("cannot release other's lock", func() {
		subject.Lock()
		lockCtl := NewLockController(testRedisKey, lockProperties)
		stt, err := lockCtl.Unlock()
		Expect(stt).To(BeFalse())
		Expect(err).To(Equal(ErrLockUnlockFailed))
		Expect(subject.HasHandle()).To(BeTrue())
		Expect(subject.IsLocked()).To(BeTrue())

	})
	It("should work correct when N go routines updates same log", func() {
		var wg sync.WaitGroup
		N := 3
		for i := 0; i < N; i++ {
			wg.Add(1)
			go func(index int) {
				lockCtl1 := NewLockController(testRedisKey, lockProperties)
				fmt.Printf("Go routine %d try to get lock %p \n", index, lockCtl1)
				lockCtl1.Lock()
				fmt.Printf("Go routine %d got the lock\n", index)
				Expect(lockCtl1.IsLocked()).To(BeTrue())
				lockCtl1.Unlock()
				fmt.Printf("Go routine %d released the lock\n", index)
				wg.Done()
			}(i)
		}
		wg.Wait()
	})
	// It("should release have error when call Unlock/IsLocked/Refresh when lock timeout", func() {
	// 	lockProperties.LockTimeout = time.Duration(3000 * time.Millisecond)
	// 	subject.Lock()
	// 	time.Sleep(5000 * time.Millisecond)
	// 	stt, err := subject.Unlock()
	// 	Expect(stt).To(BeFalse())
	// 	Expect(err).To(Equal(ErrLockUnlockFailed))
	// })
	// It("Should refresh the key with new timeout value", func() {
	// 	lockCtl, _ := Obtain(testRedisKey, lockProperties)

	// 	Expect(Obtain(testRedisKey, lockProperties))
	// 	Expect(lockCtl.Refresh()).To(BeTrue())

	// })
})
