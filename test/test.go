package main

import (
	"fmt"

	"github.com/adlindo/gocom"
	"github.com/adlindo/gocom/test/distobj"
)

func main() {

	fmt.Println("====>> ADL Common Lib Test <<====")

	gocom.AddCtrl(GetTestCtrl())
	gocom.AddCtrl(GetKeyValCtrl())
	gocom.AddCtrl(GetSecretCtrl())
	gocom.AddCtrl(GetJWTCtrl())
	gocom.AddCtrl(GetQueueCtrl())
	gocom.AddCtrl(GetPubSubCtrl())
	gocom.AddCtrl(GetDistLockCtrl())
	gocom.AddCtrl(distobj.GetDistObjCtrl())

	gocom.Start()
}
