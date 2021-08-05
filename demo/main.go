package main

import (
	"fmt"
	"reflect"
)

type Animal struct{}

func (*Animal) Ping() {
	fmt.Println("Ping")
}

func main() {
	animal := Animal{}
	val := reflect.ValueOf(&animal)
	f := val.MethodByName("Ping")
	f.Call([]reflect.Value{})

}
