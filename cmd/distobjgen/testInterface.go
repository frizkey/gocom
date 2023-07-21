package main

//go:generate distobjgen -src ITest

type ITest interface {
	TestString(a, b string) string
	TestInt(a int, b int) int
}
