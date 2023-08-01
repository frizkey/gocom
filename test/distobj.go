package main

import (
	"fmt"
	"sync"

	"github.com/adlindo/gocom"
	"github.com/adlindo/gocom/distobj"
)

type TestDist interface {
	TestString(a, b string) string
	TestInt(a, b int) int
}

// PROXY ----------------------------------------------------------------------

type TestDistProxy struct {
	className string
	prefix    string
}

var __TestDistProxyMap map[string]*TestDistProxy = map[string]*TestDistProxy{}

func GetTestDistProxy(prefix ...string) TestDist {

	targetPrefix := ""
	if len(prefix) > 0 {
		targetPrefix = prefix[0]
	}

	ret, ok := __TestDistProxyMap[targetPrefix]

	if !ok {
		ret = &TestDistProxy{}
		ret.className = "TestDist"
		ret.prefix = targetPrefix

		__TestDistProxyMap[targetPrefix] = ret
	}

	return ret
}

func (o *TestDistProxy) TestString(a, b string) string {

	retList, err := distobj.Invoke(o.prefix, o.className, "TestString", a, b)

	if err != nil {
		return ""
	}

	ret0 := distobj.ToStr(retList[0])
	return ret0
}

func (o *TestDistProxy) TestInt(a, b int) int {

	retList, err := distobj.Invoke(o.prefix, o.className, "TestInt", a, b)

	if err != nil {
		fmt.Println("error TestInt : ", err)
		return 0
	}

	fmt.Println("mau masuk : ", retList[0])
	ret0 := distobj.ToInt(retList[0])
	return ret0
}

// IMPL ----------------------------------------------------------------------

type TestDistImpl struct {
}

func (o *TestDistImpl) TestString(a, b string) string {

	return "Merge : " + a + "=>" + b
}

func (o *TestDistImpl) TestInt(a, b int) int {

	fmt.Println("masuk testint ", a, b)
	return a + b
}

//-------------------------------------------------------

type DistObjCtrl struct {
}

func (o *DistObjCtrl) Init() {

	gocom.POST("/distobj/string", o.tryString)
	gocom.POST("/distobj/int", o.tryInt)
}

type TestPayload struct {
	StrA string
	StrB string
	IntA int
	IntB int
}

func (o *DistObjCtrl) tryString(ctx gocom.Context) error {

	payload := TestPayload{}
	err := ctx.Bind(&payload)

	if err != nil {
		return ctx.SendError(gocom.NewError(1, err.Error()))
	}

	ret := GetTestDistProxy().TestString(payload.StrA, payload.StrB)
	return ctx.SendResult(ret)
}

func (o *DistObjCtrl) tryInt(ctx gocom.Context) error {

	payload := TestPayload{}
	err := ctx.Bind(&payload)

	if err != nil {
		return ctx.SendError(gocom.NewError(1, err.Error()))
	}

	ret := GetTestDistProxy().TestInt(payload.IntA, payload.IntB)

	return ctx.SendResult(ret)
}

var distObjCtrl *DistObjCtrl
var distObjOnce sync.Once

func GetDistObjCtrl() *DistObjCtrl {

	distObjOnce.Do(func() {

		distObjCtrl = &DistObjCtrl{}

		distobj.AddImpl("", "TestDist", &TestDistImpl{})
	})

	return distObjCtrl
}
