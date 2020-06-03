package main

import (
	"fmt"
	"reflect"
)

type NotknownType struct {
	s1, s2, s3 string
}

func (n NotknownType) String() string {
	return n.s1 + " - " + n.s2 + " - " + n.s3
}

// variable to investigate:
var secret interface{} = NotknownType{"Ada", "Go", "Oberon"}

func main() {
	value := reflect.ValueOf(secret) // <main.NotknownType Value>
	typ := reflect.TypeOf(secret)    // main.NotknownType
	// alternative:
	//typ := value.Type()  // main.NotknownType
	fmt.Println(typ)
	knd := value.Kind() // struct
	fmt.Println(knd)

	// iterate through the fields of the struct:
	for i := 0; i < value.NumField(); i++ {
		fmt.Printf("Field %d: %v\n", i, value.Field(i))
		// error: panic: reflect.Value.SetString using value obtained using unexported field
		//value.Field(i).SetString("C#")
	}

	// call the first method, which is String():
	results := value.Method(0).Call(nil)
	fmt.Println(results) // [Ada - Go - Oberon]

	//slice
	sli := []int{1, 2, 3}
	ptr_sli := &sli
	val_sli := reflect.ValueOf(ptr_sli).Elem()
	fmt.Printf("val_sli type:%v  \n", reflect.TypeOf(val_sli))
	fmt.Printf("val_sli value:%v \n", val_sli)

	slice := val_sli.Interface().([]int)

	fmt.Printf("val_2 type:%v\n", reflect.TypeOf(slice))
	fmt.Printf("val_2 value:%v\n", slice[2])

	type t struct {
		Err string
	}

	// []*slice reflect
	t1 := new(t)
	t1.Err = "aaa"

	t2 := new(t)
	t2.Err = "bbb"

	ts := []*t{t1, t2}

	fmt.Printf("t type:%v  \n", reflect.TypeOf(ts))
	fmt.Printf("t kind:%v  \n", reflect.TypeOf(ts).Kind())

	tf := reflect.ValueOf(ts).Index(0).Elem().FieldByName("Err").Interface().(string)
	fmt.Printf("tf value:%v  \n", tf)

}

/*
API server listening at: 127.0.0.1:50419
main.NotknownType
struct
Field 0: Ada
Field 1: Go
Field 2: Oberon
[Ada - Go - Oberon]
*/
