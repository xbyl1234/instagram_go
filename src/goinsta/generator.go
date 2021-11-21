package goinsta

import (
	"encoding/base64"
	"makemoney/tools"
	"strconv"
	"time"
)

//func generateDeviceID(seed string) string {
//	hash := generateMD5Hash(seed + volatileSeed)
//	return "android-" + hash[:16]
//}

func generateDeviceID() string {
	return "android-" + tools.GenString(tools.CharSet_16_Num, 16)
}

func generateSignature(data string) string {
	return "signed_body=SIGNATURE." + data
}

func genJazoest(pid string) string {
	var sum int = 0
	for ch := range pid {
		sum += ch
	}
	return "2" + strconv.Itoa(sum)
}

func genSnNonce(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(
		str + "|" + strconv.Itoa(int(time.Now().Unix())) + "|" +
			tools.GenString(tools.CharSet_123+tools.CharSet_abc, 24)))
}
