package common

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/klauspost/compress/gzip"
	"io"
	"io/ioutil"
	math_rand "math/rand"
	"os"
	"strings"
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
	//math_rand.Seed(time.Now().UnixNano())
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
	//math_rand.Seed(time.Now().Unix())
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

func GetCode(msg string) string {
	var index = 0
	find := false
	for index = range msg {
		if msg[index] >= '0' && msg[index] <= '9' {
			find = true
			break
		}
	}
	if find {
		code := strings.ReplaceAll(msg[index:index+7], " ", "")
		if len(code) != 6 {
			return ""
		}
		return code
	} else {
		return ""
	}
}

type encoding int

const (
	encodePath encoding = 1 + iota
	encodePathSegment
	encodeHost
	encodeZone
	encodeUserPassword
	encodeQueryComponent
	encodeFragment
	encodeInstPost
)

const upperhex = "0123456789ABCDEF"

func shouldEscape(c byte, mode encoding) bool {
	// §2.3 Unreserved characters (alphanum)
	if 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' || '0' <= c && c <= '9' {
		return false
	}

	if mode == encodeHost || mode == encodeZone {
		// §3.2.2 Host allows
		//	sub-delims = "!" / "$" / "&" / "'" / "(" / ")" / "*" / "+" / "," / ";" / "="
		// as part of reg-name.
		// We add : because we include :port as part of host.
		// We add [ ] because we include [ipv6]:port as part of host.
		// We add < > because they're the only characters left that
		// we could possibly allow, and Parse will reject them if we
		// escape them (because hosts can't use %-encoding for
		// ASCII bytes).
		switch c {
		case '!', '$', '&', '\'', '(', ')', '*', '+', ',', ';', '=', ':', '[', ']', '<', '>', '"':
			return false
		}
	}

	switch c {
	case '-', '_', '.', '~': // §2.3 Unreserved characters (mark)
		return false

	case '$', '&', '+', ',', '/', ':', ';', '=', '?', '@': // §2.2 Reserved characters (reserved)
		// Different sections of the URL allow a few of
		// the reserved characters to appear unescaped.
		switch mode {
		case encodePath: // §3.3
			// The RFC allows : @ & = + $ but saves / ; , for assigning
			// meaning to individual path segments. This package
			// only manipulates the path as a whole, so we allow those
			// last three as well. That leaves only ? to escape.
			return c == '?'

		case encodePathSegment: // §3.3
			// The RFC allows : @ & = + $ but saves / ; , for assigning
			// meaning to individual path segments.
			return c == '/' || c == ';' || c == ',' || c == '?'

		case encodeUserPassword: // §3.2.1
			// The RFC allows ';', ':', '&', '=', '+', '$', and ',' in
			// userinfo, so we must escape only '@', '/', and '?'.
			// The parsing of userinfo treats ':' as special so we must escape
			// that too.
			return c == '@' || c == '/' || c == '?' || c == ':'

		case encodeQueryComponent: // §3.4
			// The RFC reserves (so we must escape) everything.
			return true

		case encodeInstPost: // §3.4
			// The RFC reserves (so we must escape) everything.
			return true

		case encodeFragment: // §4.1
			// The RFC text is silent but the grammar allows
			// everything, so escape nothing.
			return false
		}
	}

	if mode == encodeFragment {
		// RFC 3986 §2.2 allows not escaping sub-delims. A subset of sub-delims are
		// included in reserved from RFC 2396 §2.2. The remaining sub-delims do not
		// need to be escaped. To minimize potential breakage, we apply two restrictions:
		// (1) we always escape sub-delims outside of the fragment, and (2) we always
		// escape single quote to avoid breaking callers that had previously assumed that
		// single quotes would be escaped. See issue #19917.
		switch c {
		case '!', '(', ')', '*':
			return false
		}
	}

	// Everything else must be escaped.
	return true
}

func escape(s string, mode encoding) string {
	spaceCount, hexCount := 0, 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		if shouldEscape(c, mode) {
			if c == ' ' && mode == encodeQueryComponent {
				spaceCount++
			} else {
				hexCount++
			}
		}
	}

	if spaceCount == 0 && hexCount == 0 {
		return s
	}

	var buf [64]byte
	var t []byte

	required := len(s) + 2*hexCount
	if required <= len(buf) {
		t = buf[:required]
	} else {
		t = make([]byte, required)
	}

	if hexCount == 0 {
		copy(t, s)
		for i := 0; i < len(s); i++ {
			if s[i] == ' ' {
				t[i] = '+'
			}
		}
		return string(t)
	}

	j := 0
	for i := 0; i < len(s); i++ {
		switch c := s[i]; {
		case c == ' ' && mode == encodeQueryComponent:
			t[j] = '+'
			j++
		case shouldEscape(c, mode):
			t[j] = '%'
			t[j+1] = upperhex[c>>4]
			t[j+2] = upperhex[c&15]
			j += 3
		default:
			t[j] = s[i]
			j++
		}
	}
	return string(t)
}

func InstagramQueryEscape(s string) string {
	return escape(s, encodeInstPost)
}

const TimeLayout = "2006-01-02 15:04:05"

func GetNewYorkTimeString() string {
	location, _ := time.LoadLocation("America/New_York")
	return time.Now().In(location).Format(TimeLayout)
}

func GetShanghaiTimeString() string {
	location, _ := time.LoadLocation("Asia/Shanghai")
	return time.Now().In(location).Format(TimeLayout)
}
