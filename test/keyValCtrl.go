package main

import (
	"sync"

	"github.com/adlindo/gocom"
)

type KeyValCtrl struct {
}

func (o *KeyValCtrl) Init() {

	gocom.POST("/kv", o.postKV)
	gocom.GET("/kv", o.getKV)
	gocom.DELETE("/kv", o.delKV)

	gocom.POST("/kvlist", o.postKVList)
	gocom.GET("/kvlist", o.getKVList)
}

func (o *KeyValCtrl) getKVList(ctx gocom.Context) error {

	val := gocom.KV().Range("TestKV", 0, gocom.KV().Len("TestKV"))

	return ctx.SendResult(val)
}

func (o *KeyValCtrl) postKVList(ctx gocom.Context) error {

	kv := gocom.KV()
	data := &TestDTO{}
	ctx.Bind(&data)

	if kv != nil {

		err := kv.LPush("TestKVList", data.DataString)

		if err == nil {
			return ctx.SendResult(true)
		}

		return ctx.SendError(101, "Add list error : "+err.Error())
	} else {
		return ctx.SendError(101, "invalid KV conn")
	}
}

func (o *KeyValCtrl) delKV(ctx gocom.Context) error {
	gocom.KV().Del("TestKV")

	return ctx.SendResult(true)
}

func (o *KeyValCtrl) postKV(ctx gocom.Context) error {

	kv := gocom.KV()
	data := &TestDTO{}
	ctx.Bind(&data)

	if kv != nil {

		err := kv.Set("TestKV", data.DataString)

		if err == nil {
			return ctx.SendResult(true)
		}

		kv.Set("TestKVBool", data.DataBool)

		return ctx.SendError(101, "Set error : "+err.Error())
	} else {
		return ctx.SendError(101, "invalid KV conn")
	}
}

func (o *KeyValCtrl) getKV(ctx gocom.Context) error {
	val := gocom.KV().Get("TestKV") + " ==> " + gocom.KV().Get("TestKVBool")

	return ctx.SendResult(val)
}

var keyValCtrl *KeyValCtrl
var keyValCtrlOnce sync.Once

func GetKeyValCtrl() *KeyValCtrl {

	keyValCtrlOnce.Do(func() {

		keyValCtrl = &KeyValCtrl{}
	})

	return keyValCtrl
}
