package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/adlindo/gocom"
)

var lock *gocom.DistLock

type DistLockCtrl struct {
}

func (o *DistLockCtrl) Init() {

	gocom.POST("/distlock/trylock", o.tryLock)
	gocom.POST("/distlock/release", o.releaseLock)
}

func (o *DistLockCtrl) tryLock(ctx gocom.Context) error {

	oldLock := lock
	newLock := gocom.TryLock("test", 5*time.Minute)

	if newLock != nil {

		lock = newLock
		if oldLock != nil {
			oldLock.Release()
		}

		return ctx.SendResult(lock)
	}

	return ctx.SendError(gocom.NewError(100, "Unable to get lock"))
}

func (o *DistLockCtrl) releaseLock(ctx gocom.Context) error {

	fmt.Println("lock : ", lock)
	if lock == nil {
		return ctx.SendError(gocom.NewError(100, "Not locked"))
	}

	err := lock.Release()

	if err == nil {

		return ctx.SendResult("Released")
	}

	return ctx.SendError(gocom.NewError(100, "Error release "+err.Error()))
}

//-------------------------------------------------------

var distLockCtrl *DistLockCtrl
var distLockCtrlOnce sync.Once

func GetDistLockCtrl() *DistLockCtrl {

	distLockCtrlOnce.Do(func() {

		distLockCtrl = &DistLockCtrl{}
	})

	return distLockCtrl
}
