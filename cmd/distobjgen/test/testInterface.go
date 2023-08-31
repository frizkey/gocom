package main

//go:generate distobjgen -src Test

type Test interface {
	TestString(a, b string) string
	TestInt(a int, b int) int
	TestError() error
	TestList() []string
	TestIntList() []int
	TestInterface(a interface{}) interface{}
	TestInterfaceList(a []interface{}) []interface{}
	TestMap() map[string]string
}
