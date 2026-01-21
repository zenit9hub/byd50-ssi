package logger

import (
	"fmt"
	"reflect"
	"runtime"
)

func FuncStart() {
	pc, _, _, _ := runtime.Caller(1)
	funcName := runtime.FuncForPC(pc).Name()
	fmt.Printf("\n[%v][Start] ==>\n", funcName)
}

func FuncEnd() {
	pc, _, _, _ := runtime.Caller(1)
	funcName := runtime.FuncForPC(pc).Name()
	fmt.Printf("[%v][End] <==\n\n", funcName)
}

func FuncName() string {
	pc, _, _, _ := runtime.Caller(1)
	funcName := runtime.FuncForPC(pc).Name()
	return funcName
}

func GetFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}
