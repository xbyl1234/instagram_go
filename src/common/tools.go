package common

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"crypto/tls"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/klauspost/compress/gzip"
	"golang.org/x/net/http2"
	"golang.org/x/net/proxy"
	"io"
	"io/ioutil"
	"makemoney/common/log"
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
		var tr *http.Transport
		if false {
			uri, _ := neturl.Parse("http://127.0.0.1:8888")
			tr = &http.Transport{
				Proxy: http.ProxyURL(uri),
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			}
		} else {
			dialer, _ := proxy.SOCKS5("tcp", "127.0.0.1:8889", nil, proxy.Direct)
			tr = &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
			tr.Dial = dialer.Dial
		}

		err := http2.ConfigureTransport(tr)
		if err != nil {
			log.Error("ConfigureTransport error: %v", err)
		}
		//tr.TLSClientConfig = &tls.Config{
		//	NextProtos: []string{"h2", "h2-fb", "http/1.1"},
		//	MinVersion: tls.VersionTLS10,
		//	MaxVersion: tls.VersionTLS13,
		//	CipherSuites: []uint16{
		//		tls.TLS_AES_128_GCM_SHA256,
		//		tls.TLS_AES_256_GCM_SHA384,
		//		tls.TLS_CHACHA20_POLY1305_SHA256,
		//	},
		//}

		clinet.Transport = tr
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
func GenUUID() string {
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

func GenNumber(min, max int) int {
	math_rand.Seed(time.Now().Unix())
	return math_rand.Intn(max-min) + min
}

func GZipCompress(in []byte) []byte {
	var input bytes.Buffer
	_gzip := gzip.NewWriter(&input)
	_gzip.Write(in)
	_gzip.Close()
	return input.Bytes()
}

func GZipDecompress(in []byte) ([]byte, error) {
	input := bytes.NewReader(in)
	_gzip, err := gzip.NewReader(input)
	if err != nil {
		return nil, err
	}

	var output bytes.Buffer
	_, err = _gzip.WriteTo(&output)
	return output.Bytes(), err
}

func Base64Encode(in []byte) string {
	return base64.StdEncoding.EncodeToString(in)
}
