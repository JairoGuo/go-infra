package lock

import (
	"fmt"
	"time"
	"unsafe"
)

type Locker interface {
	Acquire(lockName string) bool
	TayAcquire(lockName string, waitTime int64, expireTime int64) bool
	Release(lockName string) bool
}

var lock Locker

func Use(l Locker) {
	lock = l
}

func Lock() {
	pointer := unsafe.Pointer(&lock)
	LockWithName(fmt.Sprintf("%v", pointer))
}
func LockWithName(lockName string) {
	lock.Acquire(lockName)
}

func TryLock() bool {
	pointer := unsafe.Pointer(&lock)
	return TryLockWithName(fmt.Sprintf("%v", pointer))
}
func TryLockWithName(lockName string) bool {
	return TryLockWithParam(lockName, int64(30*time.Second), int64(30*time.Second))
}
func TryLockWithWaitTime(lockName string, waitTime int64) bool {
	return TryLockWithParam(lockName, waitTime, int64(30*time.Second))
}
func TryLockWithExpireTime(lockName string, expireTime int64) bool {
	return TryLockWithParam(lockName, int64(30*time.Second), expireTime)
}
func TryLockWithWaitAndExpireTime(lockName string, waitTime int64, expireTime int64) bool {
	return TryLockWithParam(lockName, waitTime, expireTime)
}
func TryLockWithParam(lockName string, waitTime int64, expireTime int64) bool {
	return lock.TayAcquire(lockName, waitTime, expireTime)
}
func TryLockWith(options ...Option) bool {
	param := &Param{}
	for _, opt := range options {
		opt(param)
	}
	return TryLockWithParam(param.lockName, param.waitTime, param.expireTime)
}

func Unlock() {
	pointer := unsafe.Pointer(&lock)
	lock.Release(fmt.Sprintf("%v", pointer))
}

func UnlockWithName(lockName string) {
	lock.Release(lockName)
}

type Param struct {
	lockName   string
	waitTime   int64
	expireTime int64
}

type Option func(param *Param)

func WithName(name string) Option {
	return func(param *Param) {
		param.lockName = name
	}
}

func WithWaitTime(waitTime int64) Option {
	return func(param *Param) {
		param.waitTime = waitTime
	}
}

func WithExpireTime(expireTime int64) Option {
	return func(param *Param) {
		param.expireTime = expireTime
	}
}
