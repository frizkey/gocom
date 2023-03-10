package gocom

import (
	"fmt"
	"sync"

	"github.com/adlindo/gocom/config"
)

var pubSubMap map[string]PubSubClient = map[string]PubSubClient{}
var pubSubMutex sync.Mutex
var pubSubCreatorMap map[string]PubSubCreatorFunc = map[string]PubSubCreatorFunc{}

type PubSubClient interface {
	Publish(subject string, msg interface{}) error
	Subscribe(subject string, eventHandler PubSubEventHandler)
	QueueSubscribe(subject string, queue string, eventHandler PubSubEventHandler)
}

type PubSubCreatorFunc func(connString string) (PubSubClient, error)
type PubSubEventHandler func(name string, msg string)

func RegPubSubCreator(typeName string, creator PubSubCreatorFunc) {

	pubSubCreatorMap[typeName] = creator
}

func PubSub(name ...string) PubSubClient {

	targetName := "default"

	if len(name) > 0 {
		targetName = name[0]
	}

	ret, ok := pubSubMap[targetName]

	if !ok {
		pubSubMutex.Lock()
		defer pubSubMutex.Unlock()

		if config.HasConfig("app.pubsub."+targetName+".type") && config.HasConfig("app.pubsub."+targetName+".url") {

			pubSubType := config.Get("app.pubsub." + targetName + ".type")

			creator, ok := pubSubCreatorMap[pubSubType]

			if ok {

				url := config.Get("app.pubsub." + targetName + ".url")

				var err error
				ret, err = creator(url)

				if err == nil {

					pubSubMap[targetName] = ret

					fmt.Println("Conected to PubSub :", targetName)
				} else {

					fmt.Println("Error create PubSub :", err)
				}
			}
		}
	}

	return ret
}
