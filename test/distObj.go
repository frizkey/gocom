package main

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/adlindo/gocom"
)

type TestObj interface {
	Satu() string
	Dua(data string) string
	Tiga(data string, dataInt int) string
	Empat(data string, dataInt int, dataBool bool) string
	Lima(data string, dataInt int, dataBool bool, dataFloat float64) string
}

var testClient TestObj
var testClientOnce sync.Once

type TestObjClient struct {
	ClassName string
}

func (o *TestObjClient) Satu() string {

	ret, _ := gocom.DistObj().Invoke("TestObj", "Satu")

	if len(ret) > 0 {
		return ret[0].(string)
	}

	return ""
}

func (o *TestObjClient) Dua(data string) string {

	ret, _ := gocom.DistObj().Invoke("TestObj", "Dua", data)

	if len(ret) > 0 {
		return ret[0].(string)
	}

	return ""
}

func (o *TestObjClient) Tiga(data string, dataInt int) string {

	ret, _ := gocom.DistObj().Invoke("TestObj", "Tiga", data, dataInt)

	if len(ret) > 0 {
		return ret[0].(string)
	}

	return ""
}

func (o *TestObjClient) Empat(data string, dataInt int, dataBool bool) string {

	ret, _ := gocom.DistObj().Invoke("TestObj", "Empat", data, dataInt, dataBool)

	if len(ret) > 0 {
		return ret[0].(string)
	}

	return ""
}

func (o *TestObjClient) Lima(data string, dataInt int, dataBool bool, dataFloat float64) string {

	fmt.Println("sebelum 5555")
	ret, _ := gocom.DistObj().Invoke("TestObj", "Lima", data, dataInt, dataBool, dataFloat)

	if len(ret) > 0 {
		return ret[0].(string)
	}

	return ""
}

func GetTestClient() TestObj {

	testClientOnce.Do(func() {

		testClient = &TestObjClient{}
	})

	return testClient
}

type DistObjCtrl struct {
}

func (o *DistObjCtrl) Init() {

	gocom.POST("/distobj/satu", o.postSatu)
	gocom.POST("/distobj/dua", o.postDua)
	gocom.POST("/distobj/tiga", o.postTiga)
	gocom.POST("/distobj/empat", o.postEmpat)
	gocom.POST("/distobj/lima", o.postLima)
}

func (o *DistObjCtrl) postSatu(ctx gocom.Context) error {

	data := &TestDTO{}
	ctx.Bind(&data)

	ret := GetTestClient().Satu()

	return ctx.SendResult(ret)
}

func (o *DistObjCtrl) postDua(ctx gocom.Context) error {

	data := &TestDTO{}
	ctx.Bind(&data)

	ret := GetTestClient().Dua(data.DataString)

	return ctx.SendResult(ret)
}

func (o *DistObjCtrl) postTiga(ctx gocom.Context) error {

	data := &TestDTO{}
	ctx.Bind(&data)

	ret := GetTestClient().Tiga(data.DataString, data.DataInt)

	return ctx.SendResult(ret)
}

func (o *DistObjCtrl) postEmpat(ctx gocom.Context) error {

	data := &TestDTO{}
	ctx.Bind(&data)

	ret := GetTestClient().Empat(data.DataString, data.DataInt, data.DataBool)

	return ctx.SendResult(ret)
}

func (o *DistObjCtrl) postLima(ctx gocom.Context) error {

	data := &TestDTO{}
	ctx.Bind(&data)

	ret := GetTestClient().Lima(data.DataString, data.DataInt, data.DataBool, data.DataFloat)

	return ctx.SendResult(ret)
}

//-------------------------------------------------------------------------------------

var distObjCtrl *DistObjCtrl
var distObjOnce sync.Once

type TestObjImpl struct {
}

func (o *TestObjImpl) Satu() string {

	fmt.Println("MASUK IMPL SATU")
	return "aa"
}

func (o *TestObjImpl) Dua(data string) string {

	fmt.Println("MASUK IMPL DUA")
	return "hello " + data
}

func (o *TestObjImpl) Tiga(data string, dataInt int) string {

	fmt.Println("MASUK IMPL TIGA")
	return "double " + strconv.Itoa(dataInt*2)
}

func (o *TestObjImpl) Empat(data string, dataInt int, dataBool bool) string {

	fmt.Println("MASUK IMPL EMPAT")

	if dataBool {
		return "double " + strconv.Itoa(dataInt*2)
	} else {
		return "triple " + strconv.Itoa(dataInt*3)
	}
}

func (o *TestObjImpl) Lima(data string, dataInt int, dataBool bool, dataFloat float64) string {

	fmt.Println("MASUK IMPL LIMA")
	if dataBool {
		return "double " + strconv.Itoa(dataInt*2)
	} else {
		return "triple " + fmt.Sprintf("%.3f", dataFloat*3)
	}
}

func GetDistObjCtrl() *DistObjCtrl {

	distObjOnce.Do(func() {

		distObjCtrl = &DistObjCtrl{}

		var impl TestObj = &TestObjImpl{}

		gocom.DistObj().AddImpl("TestObj", impl)
	})

	return distObjCtrl
}
