package gocom

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
)

type NatsPubSubClient struct {
	conn *nats.Conn
}

func (o *NatsPubSubClient) Publish(subject string, msg interface{}) error {

	var msgByte []byte
	var err error

	switch msg.(type) {
	case int, int16, int32, int64, string, float32, float64, bool:
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

func (o *NatsPubSubClient) Request(subject string, msg interface{}, timeOut ...time.Duration) (string, error) {

	var msgByte []byte
	var err error

	switch msg.(type) {
	case int, int16, int32, int64, string, float32, float64, bool:
		msgString := fmt.Sprintf("%v", msg)
		msgByte = []byte(msgString)
	default:
		msgByte, err = json.Marshal(msg)
		if err != nil {
			return "", err
		}
	}

	targetTimeout := 20 * time.Second

	if len(timeOut) > 0 {
		targetTimeout = timeOut[0]
	}

	ret, err := o.conn.Request(subject, msgByte, targetTimeout)

	if err != nil {
		return "", err
	}

	return string(ret.Data), nil
}

func (o *NatsPubSubClient) Subscribe(subject string, eventHandler PubSubEventHandler) {

	o.conn.Subscribe(subject, func(msg *nats.Msg) {

		defer func() {

			err := recover()

			if err != nil {

				fmt.Println("=====> SYSTEM PANIC WHEN PROCESS NATS MSG :", subject, " : ", err)
			}
		}()

		eventHandler(subject, string(msg.Data))
	})
}

func (o *NatsPubSubClient) RequestSubscribe(subject string, eventHandler PubSubReqEventHandler) {

	o.conn.Subscribe(subject, func(msg *nats.Msg) {

		defer func() {

			err := recover()

			if err != nil {

				fmt.Println("=====> SYSTEM PANIC WHEN PROCESS NATS REPLY MSG :", subject, " : ", err)
			}
		}()

		ret := eventHandler(subject, string(msg.Data))
		msg.Respond([]byte(ret))
	})
}

func (o *NatsPubSubClient) QueueSubscribe(subject string, queue string, eventHandler PubSubEventHandler) {

	o.conn.QueueSubscribe(subject, queue, func(msg *nats.Msg) {

		defer func() {

			err := recover()

			if err != nil {

				fmt.Println("=====> SYSTEM PANIC WHEN PROCESS NATS QUEUE MSG :", subject, ", Queue : ", queue, ", Error :", err)
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
