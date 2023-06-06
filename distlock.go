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
	Mode   string
	Expire time.Time
}

func (o *DistLock) IsExpired() bool {

	return time.Now().After(o.Expire)
}

func (o *DistLock) Release() error {

	if o.Mode == "keyval" {

		// check if KeyVal avail
		keyVal := KeyVal()

		if keyVal != nil {

			if keyVal.Get(o.Name) == o.ID {

				return keyVal.Del(o.Name)
			}

			return errors.New("not active lock")
		}

		return errors.New("keyval not avail")
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

		return errors.New("db not avail")
	}
}

func TryLock(name string, ttl time.Duration) *DistLock {
	ret := &DistLock{
		ID:   ulid.Make().String(),
		Name: "lock_" + name,
	}

	if ttl == 0 {
		ttl = 10 * time.Minute
	}

	// check if KeyVal avail
	keyVal := KeyVal()

	if keyVal != nil {

		ret.Mode = "keyval"
		ret.Expire = time.Now().Add(ttl)
		lockStat := keyVal.SetNX(ret.Name, ret.ID, ttl)

		if !lockStat {
			return nil
		}

		keyVal.Expire(ret.Name, ttl)
	} else {

		// check db if avail
		db := DB()

		if db != nil {

			distLockOnce.Do(func() {
				db.AutoMigrate(&DistLock{})
			})

			ret.Mode = "db"
			ret.Expire = time.Now().Add(ttl)
			existing := &DistLock{}

			db.Exec("delete from dist_locks where expire < ?", time.Now())

			if db.Create(ret).Error == nil {

				// check if we get the ownership
				db.Where("name = ?", ret.Name).Find(existing)

				if existing.ID != ret.ID {
					ret = nil
				}
			} else {
				ret = nil
			}
		} else {

			ret = nil
		}
	}

	return ret
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

	// check if KeyVal avail
	keyVal := KeyVal()

	maxTime := time.Now().Add(maxWait)

	if keyVal != nil {

		ret.Mode = "keyval"

		for {

			lockStat := keyVal.SetNX(ret.Name, ret.ID, ttl)

			if !lockStat {
				if time.Now().After(maxTime) {
					ret = nil
					break
				}
				time.Sleep(1 * time.Second)
			} else {

				ret.Expire = time.Now().Add(ttl)

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

			ret.Mode = "db"
			existing := &DistLock{}

			for {
				db.Exec("delete from dist_locks where expire < ?", time.Now())

				ret.Expire = time.Now().Add(ttl)
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
