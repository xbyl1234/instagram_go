package main

import (
	"fmt"
	"reflect"
)

type Extra2 struct {
	WaterfallId string  `json:"waterfall_id"`
	StartTime   float32 `json:"start_time"`
	ElapsedTime float32 `json:"elapsed_time"`
	Step        string  `json:"step"`
	Flow        string  `json:"flow"`
}
type name struct {
	Extra2
}

func main() {
	var n name
	t := reflect.ValueOf(n)

	b := t.FieldByName("flow")
	s := b.CanSet()

	fmt.Printf("%v", b)
	fmt.Printf("%v", s)
}
