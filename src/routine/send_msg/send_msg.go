package main

type CrawConfig struct {
	TaskName             string `json:"task_name"`
	CoroCount            int    `json:"coro_count"`
	ProxyPath            string `json:"proxy_path"`
	TargetUserDB         string `json:"target_user_db"`
	TargetUserCollection string `json:"target_user_collection"`
	CPAUrl               string `json:"cpa_url"`
}

func main() {

}
