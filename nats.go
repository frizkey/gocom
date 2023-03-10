package gocom

import (
	"encoding/json"
	"fmt"

	"github.com/nats-io/nats.go"
)

type NatsPubSubClient struct {
	conn *nats.Conn
}

func (o *NatsPubSubClient) Publish(subject string, msg interface{}) error {

	var msgByte []byte
	var err error

	switch msg.(type) {
	case int:
	case int16:
	case int32:
	case int64:
	case string:
	case float32:
	case float64:
	case bool:
		msgString := fmt.Sprintf("%v", msg)
		msgByte = []byte(msgString)
	default:
		msgByte, err = json.Marshal(msg)

		if err != nil {
			return err
		}
	}

	return o.conn.Publish(subject, msgByte)
}

func (o *NatsPubSubClient) Subscribe(subject string, eventHandler PubSubEventHandler) {

	o.conn.Subscribe(subject, func(msg *nats.Msg) {

		defer func() {

			err := recover()

			if err != nil {

				fmt.Println("=====> SYSTEM PANIC WHEN PROCESS NATS MSG :", err)
			}
		}()

		eventHandler(subject, string(msg.Data))
	})
}

func (o *NatsPubSubClient) QueueSubscribe(subject string, queue string, eventHandler PubSubEventHandler) {

	o.conn.QueueSubscribe(subject, queue, func(msg *nats.Msg) {

		defer func() {

			err := recover()

			if err != nil {

				fmt.Println("=====> SYSTEM PANIC WHEN PROCESS NATS QUEUE MSG :", err)
			}
		}()

		eventHandler(subject, string(msg.Data))
	})
}

func init() {
	RegPubSubCreator("nats", func(connString string) (PubSubClient, error) {
		ret := &NatsPubSubClient{}

		var err error
		ret.conn, err = nats.Connect(connString)

		if err != nil {
			return nil, err
		}

		return ret, nil
	})
}
