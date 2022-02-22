package goinsta

import (
	"encoding/base64"
	"makemoney/common"
	"strconv"
	"time"
)

func generateSignature(data string) string {
	return "signed_body=SIGNATURE." + data
}

func genJazoest(pid string) string {
	var sum = 0
	for ch := range pid {
		sum += ch
	}
	return "2" + strconv.Itoa(sum)
}

func genSnNonce(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(
		str + "|" + strconv.Itoa(int(time.Now().Unix())) + "|" +
			common.GenString(common.CharSet_123+common.CharSet_abc, 24)))
}
