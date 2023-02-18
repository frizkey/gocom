package main

import (
	"fmt"
	"sync"

	"github.com/adlindo/gocom"
	"github.com/jinzhu/copier"
)

type TestCtrl struct {
}

func (o *TestCtrl) Init() {

	gocom.GET("/test/hello", o.getTestHello)

	gocom.GET("/test", o.TestMultiHandler, o.TestGet)
	gocom.GET("/test/:id", o.TestGetOne)
	gocom.POST("/test", o.TestPost)
	gocom.PUT("/test/:id", o.TestPut)
	gocom.PATCH("/test/:id", o.TestPut)

	gocom.POST("/kv", o.postKV)
	gocom.GET("/kv", o.getKV)
	gocom.DELETE("/kv", o.delKV)
}

func (o *TestCtrl) getTestHello(ctx gocom.Context) error {

	return ctx.SendResult("Hello World !")
}

func (o *TestCtrl) TestGet(ctx gocom.Context) error {

	// lock := gocom.GetLock("Test", 0, 20*time.Second)

	// if lock != nil {
	// 	fmt.Println("Berhasil dapat lock")
	// } else {
	// 	fmt.Println("Gek Berhasil dapat lock")
	// }

	ret := GetTestRepo().GetAll()

	dtoRet := []TestDTO{}
	copier.Copy(&dtoRet, &ret)

	fmt.Println("Dalam TestGet")

	return ctx.SendResult(dtoRet)
}

func (o *TestCtrl) TestMultiHandler(ctx gocom.Context) error {

	fmt.Println("Dalam TestGet 2")

	return ctx.Next()
}

func (o *TestCtrl) TestGetOne(ctx gocom.Context) error {

	id := ctx.Param("id")

	ret := GetTestRepo().GetOne(id)

	if ret != nil {
		dtoRet := TestDTO{}
		copier.Copy(&dtoRet, &ret)

		return ctx.SendResult(dtoRet)
	}

	return ctx.SendError(1001, "Data not found")
}

func (o *TestCtrl) TestPost(ctx gocom.Context) error {

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

func (o *TestCtrl) TestPut(ctx gocom.Context) error {

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

func (o *TestCtrl) delKV(ctx gocom.Context) error {
	gocom.KV().Del("TestKV")

	return ctx.SendResult(true)
}

func (o *TestCtrl) postKV(ctx gocom.Context) error {

	kv := gocom.KV()
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

func (o *TestCtrl) getKV(ctx gocom.Context) error {
	val := gocom.KV().Get("TestKV")

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
