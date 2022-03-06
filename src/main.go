package main

import (
	"fmt"
	"makemoney/common"
	"makemoney/goinsta"
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
	var err error
	encodePasswd, err := goinsta.EncryptPassword("xbyl1234", "27", "LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0KTUlJQklqQU5CZ2txaGtpRzl3MEJBUUVGQUFPQ0FROEFNSUlCQ2dLQ0FRRUF5dElDa25QUWRWdG1EcCtXbEQ5YQpQNExyb1JEZDd6N0NZaElqQll5K3gyQ1J6QkFJM0NHSUh4UTdEcVBnMnJhd2RpckRINm9BTGRxNCt3RWE5QVlQCndOOHY5S0RXMURZTlN0MU9VRnZXS24zMEgrbHVkdHpGTGhCLzdiZmlkcFM1cURkZDhLWDBaNTVERVBiZHhhYksKZGNnMUtCdlB2d1l5ZjNKTnhienAxUUJoMXpkMTZCVjhqZEdsK0F0R081NHdXeHhuZzQ2MG5xRjRmRE54Q2VOcApwaVFOR1dSN0ZaeVVEUVNCRWhiNGJ3Z1dubW9sQ29KekxhRWNmN2U4OC9FZmhVWFIxakM1QklrOU5SN21hR2lvCklFUEJ4YU5pcDZrRmoxU3VEYW5RcFRSR2V3eVA1L1FpTzAvNm1uQ3dZL1pWaENMRHZyNXZramhjUk1mbk9LeUkKS3dJREFRQUIKLS0tLS1FTkQgUFVCTElDIEtFWS0tLS0tCg==")
	encodePasswd = common.InstagramQueryEscape(encodePasswd)
	if err != nil {

	}
	print(encodePasswd)
	print(encodePasswd)
}
