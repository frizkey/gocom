package main

import (
	"sync"

	"github.com/adlindo/gocom"
)

type KeyValCtrl struct {
}

func (o *KeyValCtrl) Init() {

	gocom.POST("/keyval", o.postKeyVal)
	gocom.GET("/keyval", o.getKeyVal)
	gocom.DELETE("/keyval", o.delKeyVal)

	gocom.POST("/keyvallist", o.postKeyValList)
	gocom.GET("/keyvallist", o.getKeyValList)
}

func (o *KeyValCtrl) getKeyValList(ctx gocom.Context) error {

	val := gocom.KeyVal().Range("TestKeyVal", 0, gocom.KeyVal().Len("TestKeyVal"))

	return ctx.SendResult(val)
}

func (o *KeyValCtrl) postKeyValList(ctx gocom.Context) error {

	keyVal := gocom.KeyVal()
	data := &TestDTO{}
	ctx.Bind(&data)

	if keyVal != nil {

		err := keyVal.LPush("TestKeyValList", data.DataString)

		if err == nil {
			return ctx.SendResult(true)
		}

		return ctx.SendError(gocom.NewError(101, "Add list error : "+err.Error()))
	} else {
		return ctx.SendError(gocom.NewError(101, "invalid KeyVal conn"))
	}
}

func (o *KeyValCtrl) delKeyVal(ctx gocom.Context) error {
	gocom.KeyVal().Del("TestKeyVal")

	return ctx.SendResult(true)
}

func (o *KeyValCtrl) postKeyVal(ctx gocom.Context) error {

	keyVal := gocom.KeyVal()
	data := &TestDTO{}
	ctx.Bind(&data)

	if keyVal != nil {

		err := keyVal.Set("TestKeyVal", data.DataString)

		if err == nil {
			return ctx.SendResult(true)
		}

		keyVal.Set("TestKeyValBool", data.DataBool)

		return ctx.SendError(gocom.NewError(101, "Set error : "+err.Error()))
	} else {
		return ctx.SendError(gocom.NewError(101, "invalid KeyVal conn"))
	}
}

func (o *KeyValCtrl) getKeyVal(ctx gocom.Context) error {
	val := gocom.KeyVal().Get("TestKeyVal") + " ==> " + gocom.KeyVal().Get("TestKeyValBool")

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
