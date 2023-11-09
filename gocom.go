package gocom

import (
	"github.com/frizkey/gocom/keyval"
	"github.com/frizkey/gocom/pubsub"
	"github.com/frizkey/gocom/queue"
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
