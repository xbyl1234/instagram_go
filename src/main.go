package main

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

}
