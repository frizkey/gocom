package gocom

import (
	"fmt"
	"sync"

	"github.com/adlindo/gocom/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var dbMap map[string]*gorm.DB = map[string]*gorm.DB{}
var dbMutex sync.Mutex
var dbCreatorMap map[string]DBCreatorFunc = map[string]DBCreatorFunc{}

type DBCreatorFunc func(dsn string) (*gorm.DB, error)

func RegDBCreator(typeName string, creator DBCreatorFunc) {

	dbCreatorMap[typeName] = creator
}

func DBConnByName(name string) *gorm.DB {

	if name == "" {
		name = "default"
	}

	ret, ok := dbMap[name]

	if !ok {

		dbMutex.Lock()
		defer dbMutex.Unlock()

		// check whether other thread already create the db
		ret, ok = dbMap[name]

		if !ok {

			if config.HasConfig("app.db."+name+".type") && config.HasConfig("app.db."+name+".dsn") {

				dbType := config.Get("app.db." + name + ".type")

				creator, ok := dbCreatorMap[dbType]

				if ok {

					dsn := config.Get("app.db." + name + ".dsn")
					poolSize := config.GetInt("app.db."+name+".poolSize", 20)
					debug := config.GetBool("app.db."+name+".debug", false)

					var err error
					ret, err = creator(dsn)

					if err == nil {

						dbMap[name] = ret

						if debug {
							ret.Debug()
						}

						if poolSize > 0 {

							db, errDB := ret.DB()

							if errDB == nil {

								db.SetMaxOpenConns(poolSize)
							}
						}

						fmt.Println("Conected to DB :", name)
					} else {

						fmt.Println("Error create DB :", err)
					}
				}
			}
		}
	}

	return ret
}

func DBConn() *gorm.DB {

	return DBConnByName("default")
}

func init() {

	RegDBCreator("postgresql", func(dsn string) (*gorm.DB, error) {

		return gorm.Open(postgres.Open(dsn), &gorm.Config{})
	})
}
