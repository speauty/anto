package main

import (
	"fmt"
	"reflect"
)

type Test struct {
	Id   string `json:"id"`
	Name string `json:"-"`
}

func main() {
	newTest := new(Test)
	newTest.Id = "fsds"
	testType := reflect.TypeOf(newTest)
	testVal := reflect.ValueOf(newTest).Elem()
	lenFields := testType.Elem().NumField()
	for i := 0; i < lenFields; i++ {
		fieldName := testType.Elem().Field(i)
		jsobTagVal := fieldName.Tag.Get("json")
		fmt.Println(jsobTagVal, testVal.Field(i).Interface())
	}
}
