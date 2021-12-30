package main

import (
	"encoding/json"
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

	imageCompression, _ := json.Marshal(map[string]interface{}{
		"lib_name":    "uikit",
		"lib_version": "1575.230000",
		"quality":     99.9,
		"colorspace":  "kCGColorSpaceDeviceRGB",
		"ssim":        1,
	})
	fmt.Printf("%v", string(imageCompression))
	fmt.Printf("%v", s)
}
