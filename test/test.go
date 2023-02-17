package main

import (
	"fmt"

	"github.com/adlindo/gocom/ctrl"
)

func main() {

	fmt.Println("====>> ADL Common Lib Test <<====")

	ctrl.Add(GetTestCtrl())
	ctrl.Start()
}
