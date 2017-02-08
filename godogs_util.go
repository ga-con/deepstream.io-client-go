package main

import (
	"reflect"
	"time"
)

func someMsLater(ms int) func() error {
	return func() error {
		time.Sleep(time.Duration(ms) * time.Millisecond)
		return nil
	}
}

func compareData(data1, data2 []interface{}) bool {
	return reflect.DeepEqual(data1, data2)
}
