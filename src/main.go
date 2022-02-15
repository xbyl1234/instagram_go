package main

import (
	"fmt"
	"time"
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

func GenUploadID() string {
	upId := fmt.Sprintf("%d", time.Now().UnixMicro())
	return upId[:len(upId)-1]
}

func main() {

	print(GenUploadID())
}
