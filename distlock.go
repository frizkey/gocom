package gocom

import (
	"errors"
	"sync"
	"time"

	"github.com/oklog/ulid/v2"
)

var distLockOnce sync.Once

type DistLock struct {
	ID     string
	Name   string `gorm:"primarykey"`
	Expire time.Time
}

func GetLock(name string, maxWait time.Duration, ttl time.Duration) *DistLock {

	ret := &DistLock{
		ID:   ulid.Make().String(),
		Name: "lock_" + name,
	}

	if maxWait == 0 {
		maxWait = 10 * time.Minute
	}

	if ttl == 0 {
		ttl = 10 * time.Minute
	}

	// check if KV avail
	kv := KV()

	maxTime := time.Now().Add(maxWait)

	if kv != nil {

		for {

			lockStat := kv.SetNX(ret.Name, ret.ID, ttl)

			if !lockStat {
				if time.Now().After(maxTime) {
					ret = nil
					break
				}
				time.Sleep(1 * time.Second)
			} else {
				break
			}
		}

		return ret
	} else {

		// check db if avail
		db := DB()

		if db != nil {

			distLockOnce.Do(func() {
				db.AutoMigrate(&DistLock{})
			})

			ret.Expire = time.Now().Add(ttl)
			existing := &DistLock{}

			for {
				db.Exec("delete from dist_locks where expire < ?", time.Now())

				if db.Create(ret).Error == nil {

					// check if we get the ownership
					db.Where("name = ?", ret.Name).Find(existing)

					if existing.ID == ret.ID {
						break
					}
				}

				if time.Now().After(maxTime) {
					ret = nil
					break
				}

				time.Sleep(1 * time.Second)
			}

			return ret
		}

		return nil
	}
}

func ReleaseLock(lock *DistLock) error {

	return lock.Release()
}

func (o *DistLock) Release() error {

	// check if KV avail
	kv := KV()

	if kv != nil {

		if kv.Get(o.Name) == o.ID {

			return kv.Del(o.Name)
		}

		return errors.New("not active lock")
	} else {

		// check db if avail
		db := DB()

		if db != nil {

			distLockOnce.Do(func() {
				db.AutoMigrate(&DistLock{})
			})

			existing := &DistLock{}
			db.Where("name = ?", o.Name).Find(existing)

			if existing.ID == o.ID {
				return db.Exec("delete from dist_locks where name = ? and id = ?", o.Name, o.ID).Error
			}

			return errors.New("not active lock")
		}

		return nil
	}
}
