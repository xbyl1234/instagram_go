package common

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"makemoney/config"
	math_rand "math/rand"
	"net/http"
	neturl "net/url"
	"os"
	"time"
	"unsafe"
)

const (
	volatileSeed   = "12345"
	CharSet_abc    = "abcdefghijklmnopqrstuvwxyz"
	CharSet_ABC    = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	CharSet_123    = "0123456789"
	CharSet_16_Num = "0123456789abcdef"
	CharSet_All    = CharSet_abc + CharSet_ABC + CharSet_123
)

func LoadJsonFile(path string, ret interface{}) error {
	by, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	err = json.Unmarshal(by, ret)
	return err
}

func Dumps(path string, obj interface{}) error {
	by, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	_, err = file.Write(by)
	if err != nil {
		return err
	}
	_ = file.Close()
	return nil
}

func DebugHttpClient(clinet *http.Client) {
	if config.UseCharles {
		uri, err := neturl.Parse("http://127.0.0.1:8888")
		if err == nil {
			clinet.Transport = &http.Transport{
				Proxy: http.ProxyURL(uri),
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			}
		}
	}
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}
func NewUUID() (string, error) {
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

func GenString(charSet string, length int) string {
	by := make([]byte, length)
	math_rand.Seed(time.Now().UnixNano())
	for index := 0; index < length; index++ {
		by[index] = charSet[math_rand.Intn(len(charSet))]
	}
	return B2s(by)
}
func GenerateUUID() string {
	uuid, err := NewUUID()
	if err != nil {
		return "cb479ee7-a50d-49e7-8b7b-60cc1a105e22" // default value when error occurred
	}
	return uuid
}

func GenerateMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func GenerateHMAC(text, key string) string {
	hasher := hmac.New(sha256.New, []byte(key))
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func B2s(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func Json2String(params map[string]string) (string, error) {
	data, err := json.Marshal(params)
	if err != nil {
		return "", err
	}
	return B2s(data), nil
}
