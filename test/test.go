package main

import (
	"fmt"

	"github.com/adlindo/gocom"
)

func main() {

	fmt.Println("====>> ADL Common Lib Test <<====")

	gocom.AddController(GetTestCtrl())
	gocom.Start()
}
