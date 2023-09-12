package queue

import (
	"encoding/json"
	"fmt"

	"github.com/nats-io/nats.go"
)

type NatsQueueClient struct {
	conn *nats.Conn
}

func (o *NatsQueueClient) Publish(name string, payload interface{}) error {

	var msgByte []byte
	var err error

	switch payload.(type) {
	case int, int16, int32, int64, string, float32, float64, bool:
		msgString := fmt.Sprintf("%v", payload)
		msgByte = []byte(msgString)
	default:
		msgByte, err = json.Marshal(payload)

		if err != nil {
			return err
		}
	}

	return o.conn.Publish(name, msgByte)
}

func (o *NatsQueueClient) Consume(name string, consumer QueueConsumerFunc) {

	sub, err := o.conn.QueueSubscribe(name, name, func(msg *nats.Msg) {

		defer func() {

			err := recover()

			if err != nil {

				fmt.Println("=====> SYSTEM PANIC WHEN PROCESS NATS QUEUE MSG :", name, ", Error :", err)
			}
		}()

		consumer(name, string(msg.Data))
	})

	if err == nil {

		no := 10000000
		sub.SetPendingLimits(no, no*1024)
	}
}

func init() {
	RegQueueCreator("nats", func(url string) (QueueClient, error) {

		ret := &NatsQueueClient{}

		var err error
		ret.conn, err = nats.Connect(url)

		if err != nil {
			return nil, err
		}

		return ret, nil
	})
}
