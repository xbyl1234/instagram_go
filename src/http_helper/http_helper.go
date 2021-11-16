package http_helper

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

type RequestOpt struct {
	Params   map[string]string //form
	Header   map[string]string
	IsPost   bool
	ReqUrl   string
	Data     string
	JsonData interface{}
}

var defaultUserAgent string = "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.87 UBrowser/6.2.4098.3 Safari/537.36"

func SetDefaultUserAgent(userAgent string) {
	defaultUserAgent = userAgent
}

func httpDo(client *http.Client, opt *RequestOpt) (*http.Response, error) {
	urlParams := url.Values{}

	if opt.Params != nil {
		for k, v := range opt.Params {
			urlParams.Set(k, v)
		}
	}

	url, _ := url.Parse(opt.ReqUrl)
	body := bytes.NewBuffer([]byte{})

	var method string
	if opt.IsPost {
		method = "POST"
		if len(urlParams) != 0 {
			body.WriteString(urlParams.Encode())
		} else if opt.JsonData != nil {
			jsonData, err := json.Marshal(opt.JsonData)
			if err != nil {
				return nil, err
			}
			body.Write(jsonData)
		} else if opt.Data != "" {
			body.WriteString(opt.Data)
		}
	} else {
		method = "GET"
		if url.RawQuery != "" && len(urlParams) != 0 {
			url.RawQuery += "&"
		}
		url.RawQuery += urlParams.Encode()
	}

	req, _ := http.NewRequest(method, url.String(), body)
	req.Header.Set("User-Agent", defaultUserAgent)
	req.Header.Set("Connection", "keep-alive")

	if opt.Header != nil {
		for key, vul := range opt.Header {
			req.Header.Set(key, vul)
		}
	}

	resp, err := client.Do(req)
	return resp, err
}

func fetchHttpText(resp *http.Response) (string, error) {
	context, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(context), err
}

func fetchHttpJson(resp *http.Response, response interface{}) error {
	context, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(context, response)
	return err
}

func HttpDo(client *http.Client, opt *RequestOpt) (string, error) {
	resp, err := httpDo(client, opt)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	return fetchHttpText(resp)
}

func HttpDoJson(client *http.Client, opt *RequestOpt, response interface{}) error {
	resp, err := httpDo(client, opt)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return fetchHttpJson(resp, response)
}

const upperhex = "0123456789ABCDEF"

type encoding int

const (
	encodePath encoding = 1 + iota
	encodePathSegment
	encodeHost
	encodeZone
	encodeUserPassword
	encodeQueryComponent
	encodeFragment
	encodeURI
)

func EncodeForm(s string) string {
	return Escape(s, encodePath)
}

func Escape(s string, mode encoding) string {
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

// Return true if the specified character should be escaped when
// appearing in a URL string, according to RFC 3986.
//
// Please be informed that for now shouldEscape does not check all
// reserved characters correctly. See golang.org/issue/5684.

func shouldEscape(c byte, mode encoding) bool {
	// §2.3 Unreserved characters (alphanum)
	if 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' || '0' <= c && c <= '9' {
		return false
	}

	//:/?~%&+=;,@()!*#$'
	switch c {
	case '-', '_', '.', '~':
		return false
	case '!', '$', '&', '\'', '(', ')', '*', '+', ',', ';', '=', ':', '/', '?', '%', '@', '#':
		return false
	}
	// Everything else must be escaped.
	return true
}
