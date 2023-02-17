package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/adlindo/gocom"
	"github.com/adlindo/gocom/ctrl"
	"github.com/jinzhu/copier"
)

type TestCtrl struct {
}

func (o *TestCtrl) Init() {

	ctrl.GET("/test/hello", o.getTestHello)

	ctrl.GET("/test", o.TestMultiHandler, o.TestGet)
	ctrl.GET("/test/:id", o.TestGetOne)
	ctrl.POST("/test", o.TestPost)
	ctrl.PUT("/test/:id", o.TestPut)
	ctrl.PATCH("/test/:id", o.TestPut)

	ctrl.POST("/kv", o.postKV)
	ctrl.GET("/kv", o.getKV)
	ctrl.DELETE("/kv", o.delKV)
}

func (o *TestCtrl) getTestHello(ctx ctrl.Context) error {

	return ctx.SendResult("Hello World !")
}

func (o *TestCtrl) TestGet(ctx ctrl.Context) error {

	lock := gocom.GetLock("Test", 0, 20*time.Second)

	if lock != nil {
		fmt.Println("Berhasil dapat lock")
	} else {
		fmt.Println("Gek Berhasil dapat lock")
	}

	ret := GetTestRepo().GetAll()

	dtoRet := []TestDTO{}
	copier.Copy(&dtoRet, &ret)

	fmt.Println("Dalam TestGet")

	return ctx.SendResult(dtoRet)
}

func (o *TestCtrl) TestMultiHandler(ctx ctrl.Context) error {

	fmt.Println("Dalam TestGet 2")

	return ctx.Next()
}

func (o *TestCtrl) TestGetOne(ctx ctrl.Context) error {

	id := ctx.Param("id")

	ret := GetTestRepo().GetOne(id)

	if ret != nil {
		dtoRet := TestDTO{}
		copier.Copy(&dtoRet, &ret)

		return ctx.SendResult(dtoRet)
	}

	return ctx.SendError(1001, "Data not found")
}

func (o *TestCtrl) TestPost(ctx ctrl.Context) error {

	fmt.Println("Dalam TestPost")

	dto := &TestDTO{}
	ctx.Bind(dto)
	fmt.Println("====>", dto)

	mdl := &Test{}
	copier.Copy(mdl, dto)

	GetTestRepo().Create(mdl)

	fmt.Println(mdl)

	copier.Copy(dto, mdl)

	return ctx.SendResult(dto)
}

func (o *TestCtrl) TestPut(ctx ctrl.Context) error {

	fmt.Println("Dalam TestPut")

	dto := &TestDTO{}
	ctx.Bind(dto)
	fmt.Println("====>", dto)

	mdl := &Test{}
	copier.Copy(mdl, dto)
	mdl.ID = ctx.Param("id")

	GetTestRepo().Update(mdl)

	fmt.Println(mdl)

	copier.Copy(dto, mdl)

	return ctx.SendResult(dto)
}

func (o *TestCtrl) delKV(ctx ctrl.Context) error {
	gocom.KVConn().Del("TestKV")

	return ctx.SendResult(true)
}

func (o *TestCtrl) postKV(ctx ctrl.Context) error {

	kv := gocom.KVConn()
	data := &TestDTO{}
	ctx.Bind(&data)

	if kv != nil {

		err := kv.Set("TestKV", data.DataString)

		if err == nil {
			return ctx.SendResult(true)
		}

		return ctx.SendError(101, "Set error : "+err.Error())
	} else {
		return ctx.SendError(101, "invalid KV conn")
	}
}

func (o *TestCtrl) getKV(ctx ctrl.Context) error {
	val := gocom.KVConn().Get("TestKV")

	return ctx.SendResult(val)
}

//-----------------------------------------------

var testCtrl *TestCtrl
var testCtrlOnce sync.Once

func GetTestCtrl() *TestCtrl {

	testCtrlOnce.Do(func() {

		testCtrl = &TestCtrl{}
	})

	return testCtrl
}
