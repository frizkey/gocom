package gocom

import (
	"fmt"
	"sync"
	"time"

	"github.com/adlindo/gocom/config"
)

// --------------------------

type KVClient interface {
	Set(key string, val interface{}, ttl ...time.Duration) error
	SetNX(key string, val interface{}, ttl ...time.Duration) bool
	Get(key string) interface{}
	GetString(key string) string
	GetInt(key string) int
	Del(key string) error

	LPush(key string, val interface{}) error
	LPop(key string) interface{}
	LPopString(key string) string
	LPopInt(key string) int
	RPush(key string, val interface{}) error
	RPop(key string) interface{}
	RPopString(key string) string
	RPopInt(key string) int
	Len(key string) int
	AtIndex(key string, index int) interface{}
	AtIndexString(key string, index int) string
	AtIndexInt(key string, index int) int
}

// --------------------------

var kvMap map[string]KVClient = map[string]KVClient{}
var kvOnce sync.Once
var kvCreatorMap map[string]KVCreatorFunc = map[string]KVCreatorFunc{}

type KVCreatorFunc func(url string) (KVClient, error)

func RegKVCreator(typeName string, creator KVCreatorFunc) {

	kvCreatorMap[typeName] = creator
}

func KV(name ...string) KVClient {

	targetName := "default"

	if len(name) > 0 {
		targetName = name[0]
	}

	ret, ok := kvMap[targetName]

	if !ok {

		kvOnce.Do(func() {

			if config.HasConfig("app.kv."+targetName+".type") && config.HasConfig("app.kv."+targetName+".url") {

				kvType := config.Get("app.kv." + targetName + ".type")

				creator, ok := kvCreatorMap[kvType]

				if ok {

					url := config.Get("app.kv." + targetName + ".url")

					var err error
					ret, err = creator(url)

					if err == nil {

						kvMap[targetName] = ret

						fmt.Println("Conected to KV :", targetName)
					} else {

						fmt.Println("Error create KV :", err)
					}
				}
			}
		})
	}

	return ret
}
