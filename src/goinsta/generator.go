package goinsta

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	math_rand "math/rand"
	"strconv"
	"time"
)

const (
	volatileSeed   = "12345"
	charSet_abc    = "abcdefghijklmnopqrstuvwxyz"
	charSet_ABC    = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	charSet_123    = "0123456789"
	charSet_16_Num = "0123456789abcdef"
	charSet_All    = charSet_abc + charSet_ABC + charSet_123
)

func generateMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func generateHMAC(text, key string) string {
	hasher := hmac.New(sha256.New, []byte(key))
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

//func generateDeviceID(seed string) string {
//	hash := generateMD5Hash(seed + volatileSeed)
//	return "android-" + hash[:16]
//}

func generateDeviceID() string {
	return "android-" + genString(charSet_16_Num, 16)
}

func newUUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}

func generateUUID() string {
	uuid, err := newUUID()
	if err != nil {
		return "cb479ee7-a50d-49e7-8b7b-60cc1a105e22" // default value when error occurred
	}
	return uuid
}

func generateSignature(data string) map[string]string {
	m := make(map[string]string)
	m["ig_sig_key_version"] = goInstaSigKeyVersion
	m["signed_body"] = fmt.Sprintf(
		"%s.%s", generateHMAC(data, goInstaIGSigKey), data,
	)
	return m
}

func genString(charSet string, length int) string {
	by := make([]byte, length)
	math_rand.Seed(time.Now().UnixNano())
	for index := 0; index < length; index++ {
		by[index] = charSet[math_rand.Intn(len(charSet))]
	}
	return b2s(by)
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
		str + "|" + strconv.Itoa(int(time.Now().Unix())) + "|" + genString(charSet_123+charSet_abc, 24)))
}
