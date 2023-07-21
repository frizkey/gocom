package gocom

import (
	"github.com/adlindo/gocom/keyval"
	"github.com/adlindo/gocom/pubsub"
	"github.com/adlindo/gocom/queue"
)

func KeyVal(name ...string) keyval.KeyValClient {
	return keyval.Get(name...)
}

func PubSub(name ...string) pubsub.PubSubClient {
	return pubsub.Get(name...)
}

func Queue(name ...string) queue.QueueClient {
	return queue.Get(name...)
}
