package gocom

import (
	"fmt"
	"sync"
	"time"

	"github.com/adlindo/gocom/config"
)

// --------------------------

type KV interface {
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

var kvMap map[string]KV = map[string]KV{}
var kvOnce sync.Once
var kvCreatorMap map[string]KVCreatorFunc = map[string]KVCreatorFunc{}

type KVCreatorFunc func(url string) (KV, error)

func RegKVCreator(typeName string, creator KVCreatorFunc) {

	kvCreatorMap[typeName] = creator
}

func KVConnByName(name string) KV {

	if name == "" {
		name = "default"
	}

	ret, ok := kvMap[name]

	if !ok {
		fmt.Println("aaaaaaaaaaaaa")
		kvOnce.Do(func() {

			fmt.Println("bbbbbbbbbbbbbbbbb")
			if config.HasConfig("app.kv."+name+".type") && config.HasConfig("app.kv."+name+".url") {

				kvType := config.Get("app.kv." + name + ".type")

				creator, ok := kvCreatorMap[kvType]

				fmt.Println("cccccccccccc")
				if ok {

					url := config.Get("app.kv." + name + ".url")

					fmt.Println("ddddddddddddd")
					var err error
					ret, err = creator(url)

					fmt.Println("eeeeeeeeeeeeeee")

					if err == nil {

						kvMap[name] = ret

						fmt.Println("Conected to KV :", name)
					} else {

						fmt.Println("Error create KV :", err)
					}
				}
			}
		})
	}

	return ret
}

func KVConn() KV {
	return KVConnByName("default")
}
