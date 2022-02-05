package utils

import (
	"reflect"
	"runtime"
)

type StepFn func() bool

// This could be slow!
func GetFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}
