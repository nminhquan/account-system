package transaction

import (
	"errors"
	"log"
	"sync"
	"time"
)

type LockProperties struct {
	LockAttempts     int
	LockAttemptDelay time.Duration
	LockTimeout      time.Duration
}

type Lock interface {
	Lock() (bool, error)
	Unlock() (bool, error)
	IsLocked() (bool, error)
	HasHandle() (bool, error)
}

const (
	LockDefaultState = "1"
	// LockDefaultAttempt    = 5
	// LockDefaultRetryDelay = 1000
	// LockDefaultTimeout    = 4500
)

var (
	ErrLockUnlockFailed     = errors.New("lock unlock failed")
	ErrLockNotObtained      = errors.New("lock not obtained")
	ErrLockDurationExceeded = errors.New("lock duration exceeded")
)

type GlobalLock struct {
	lockId     string
	properties *LockProperties
	lockHandle string // only 1 thread have lock handle at one time for the same lockId
	mutex      sync.Mutex
}

/*
	Create new lock controller which will control the status of the lock
	To lock: NewLockController then call Lock
	To unclok: Obtain the lock then call unlock
*/
func NewLockController(lockId string, properties *LockProperties) Lock {
	// if properties == nil {
	// 	properties = &LockProperties{LockDefaultAttempt, LockDefaultRetryDelay, LockDefaultTimeout}
	// }
	lock := &GlobalLock{lockId: lockId, properties: properties}
	return lock
}

// short cut for NewLockController().Lock(), used when releasing lock
func Obtain(lockId string, properties *LockProperties) (Lock, error) {
	log.Printf("obtain %v", *properties)
	locker := NewLockController(lockId, properties)

	if ok, _ := locker.Lock(); !ok {
		return nil, ErrLockNotObtained
	}

	return locker, nil
}

func (l *GlobalLock) Lock() (bool, error) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	_, err := l.createLock()
	if err != nil {
		return false, err
	}

	return true, nil
}

func (l *GlobalLock) Unlock() (bool, error) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	ok1, _ := l.IsLocked()
	// ok2, _ := l.HasHandle() // Skip this line because lock will be called remotely, sometimes the lockCtl doesn't have the handle
	if !ok1 {
		return false, ErrLockUnlockFailed
	}
	_, err := l.releaseLock()
	if err != nil {
		return false, err
	}

	return true, nil
}

func (l *GlobalLock) HasHandle() (bool, error) {
	return l.lockHandle != "", nil
}

func (l *GlobalLock) IsLocked() (bool, error) {
	exists, err := txnDao.CheckLock(l.lockId)
	if err != nil {
		return false, err
	} else if exists == 0 {
		return false, nil
	}
	return true, nil
}

//helper func
func (l *GlobalLock) createLock() (bool, error) {
	log.Printf("CREATE lock %v", *l.properties)
	attempts := l.properties.LockAttempts
	delay := l.properties.LockAttemptDelay
	for {
		ok, err := txnDao.CreateLock(l.lockId, LockDefaultState, l.properties.LockTimeout)

		if err != nil {
			return false, err
		} else if ok {
			l.lockHandle = l.lockId
			log.Printf("[Lock] CREATED lock, handle = %v", l.lockHandle)
			return true, nil
		}
		log.Println("[Lock] Attempting retries on create lock retry: ", attempts)
		if attempts--; attempts <= 0 {
			log.Printf("[Lock] CREATE lock timeout, cannot get lock, return error")
			return false, ErrLockNotObtained
		}
		time.Sleep(delay)
	}
}

func (l *GlobalLock) refreshLock() (bool, error) {
	log.Printf("[Lock] REFRESH lock %p", l)
	ok, err := txnDao.RefreshLock(l.lockId, l.properties.LockTimeout)
	if err != nil {
		return false, err
	} else if !ok {
		return false, ErrLockDurationExceeded
	}

	return true, nil
}

func (l *GlobalLock) releaseLock() (bool, error) {
	log.Printf("[Lock] RELEASE lock %p", l)
	_, err := txnDao.DeleteLock(l.lockId)
	if err != nil {
		return false, ErrLockUnlockFailed
	}
	l.lockHandle = ""
	return true, nil
}
