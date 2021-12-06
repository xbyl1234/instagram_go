package goinsta

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/klauspost/compress/gzip"
	"github.com/klauspost/compress/zstd"
	"io"
	"io/ioutil"
	"makemoney/common"
	"makemoney/common/log"
	"makemoney/config"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type reqOptions struct {
	ApiPath       string
	IsPost        bool
	IsApiB        bool
	IsApiGraph    bool
	Signed        bool
	Query         map[string]interface{}
	Body          *bytes.Buffer
	Header        map[string]string
	DisAutoHeader bool
}

type BaseApiResp struct {
	url      string
	username string

	Status    string `json:"status"`
	ErrorType string `json:"error_type"`
	Message   string `json:"message"`
}

func (this *BaseApiResp) SetInfo(url string, username string) {
	this.username = username
	this.url = url
}

func (this *BaseApiResp) isError() bool {
	return this.Status != "ok"
}

func (this *BaseApiResp) CheckError(err error) error {
	if err != nil {
		return err
	}
	if this == nil {
		return &common.MakeMoneyError{ErrStr: this.Message, ErrType: common.OtherError}
	}
	if this.Status != "ok" {
		log.Warn("account: %s, url: %s, api error: %s",
			this.username,
			this.url,
			this.ErrorType+":"+this.Message)
		if this.Message == InsAccountError_ChallengeRequired {
			return &common.MakeMoneyError{ErrStr: this.Message, ErrType: common.ChallengeRequiredError}
		} else {
			return &common.MakeMoneyError{ErrStr: this.Message, ErrType: common.ApiError}
		}
	}
	return nil
}

var (
	IGHeader_EncryptionId           string = "password-encryption-key-id"
	IGHeader_EncryptionKey          string = "password-encryption-pub-key"
	IGHeader_Authorization          string = "authorization"
	IGHeader_udsUserID              string = "ig-u-ds-user-id"
	IGHeader_iguiggDirectRegionHint string = "ig-u-ig-direct-region-hint"
	IGHeader_iguShbid               string = "ig-u-shbid"
	IGHeader_iguShbts               string = "ig-u-shbts"
	IGHeader_iguRur                 string = "ig-u-rur"
	IGHeader_UseAuthHeaderForSso    string = "use-auth-header-for-sso"
	IGHeader_XMid                   string = "x-mid"
	IGHeader_igwwwClaim             string = "x-ig-www-claim"
)

func SetHeader(req *http.Request, key string, vul string) {
	req.Header[key] = []string{vul}
}

func (this *Instagram) setBaseHeader(req *http.Request) {
	req.Header.Set("connection", "keep-alive")
	req.Header.Set("accept-language", "zh-CN, en-US")
	req.Header.Set("user-agent", this.UserAgent)
	req.Header.Set("accept-encoding", "zstd, gzip, deflate")

	if req.Header.Get("content-type") == "" {
		SetHeader(req, "content-type", "application/x-www-form-urlencoded; charset=UTF-8")
	}

	SetHeader(req, "x-ig-family-device-id", this.familyID)
	SetHeader(req, "ig-intended-user-id", strconv.FormatInt(this.ID, 10))
	SetHeader(req, "x-ig-app-id", AppID)
	SetHeader(req, "x-ig-capabilities", "3brTvx0=")
	SetHeader(req, "x-ig-connection-type", "WIFI")
	SetHeader(req, "x-ig-device-id", this.uuid)
	SetHeader(req, "x-ig-android-id", this.androidID)

	igwwwClaim := this.ReadHeader(IGHeader_igwwwClaim)
	if igwwwClaim == "" {
		igwwwClaim = "0"
	}
	SetHeader(req, IGHeader_igwwwClaim, igwwwClaim)

	SetHeader(req, "x-ig-timezone-offset", "true")

	SetHeader(req, "x-bloks-version-id", BloksVersionID)
	SetHeader(req, "x-bloks-is-layout-rtl", "false")
	SetHeader(req, "x-bloks-is-panorama-enabled", "true")

	SetHeader(req, "x-ig-app-locale", goInstaLocation)
	SetHeader(req, "x-ig-device-locale", goInstaLocation)
	SetHeader(req, "x-ig-mapped-locale", goInstaLocation)

	SetHeader(req, "x-ig-bandwidth-speed-kbps", "-1.000")
	SetHeader(req, "x-ig-bandwidth-totalbytes-b", "0")
	SetHeader(req, "x-ig-bandwidth-totaltime-ms", "0")

	SetHeader(req, "x-fb-client-ip", "True")
	SetHeader(req, "x-fb-http-engine", "Liger")
	SetHeader(req, "x-fb-server-cluster", "True")

	SetHeader(req, "x-pigeon-session-id", this.sessionID)
}

func (this *Instagram) setLoginHeader(req *http.Request) {
	SetHeader(req, IGHeader_udsUserID, strconv.FormatInt(this.ID, 10))
	SetHeader(req, IGHeader_iguRur, this.ReadHeader(IGHeader_iguRur))
	SetHeader(req, IGHeader_XMid, this.ReadHeader(IGHeader_XMid))
	SetHeader(req, IGHeader_Authorization, this.ReadHeader(IGHeader_Authorization))
	SetHeader(req, "x-ig-app-startup-country", "OR")
	SetHeader(req, "x-pigeon-rawclienttime", strconv.FormatInt(time.Now().Unix(), 10))
}

func (this *Instagram) setHeader(reqOpt *reqOptions, req *http.Request) {
	this.setBaseHeader(req)
	if this.IsLogin {
		this.setLoginHeader(req)
	}

	//SetHeader(req,"x-ig-connection-speed", fmt.Sprintf("%dkbps", common.GenNumber(1000, 3700)))
	//SetHeader(req,"x-ads-opt-out", "0")
	//SetHeader(req,"x-cm-latency", "-1.000")
	//SetHeader(req,"x-ig-extended-cdn-thumbnail-cache-busting-value", "1000")
}

func (this *Instagram) afterRequest(reqUrl *url.URL, resp *http.Response) {
	_url, _ := url.Parse(goInstaHost)
	for _, value := range this.c.Jar.Cookies(_url) {
		if strings.Contains(value.Name, "csrftoken") {
			this.token = value.Value
		}
	}

	for key := range resp.Header {
		setting := strings.ToLower(key)
		if strings.Index(setting, "ig-set-") == 0 {
			this.httpHeader[setting[len("ig-set-"):]] = resp.Header.Get(key)
		}
	}
}

func (this *Instagram) httpDo(reqOpt *reqOptions) ([]byte, error) {
	method := "GET"
	if reqOpt.IsPost {
		method = "POST"
	}

	var baseUrl string
	if reqOpt.IsApiB {
		baseUrl = goInstaHost_B
	} else if reqOpt.IsApiGraph {
		baseUrl = goInstaHost_Graph
	} else {
		baseUrl = goInstaHost
	}

	_url, err := url.Parse(baseUrl + reqOpt.ApiPath)
	if err != nil {
		return nil, err
	}

	var bf *bytes.Buffer
	if reqOpt.Query != nil {
		bf = bytes.NewBuffer([]byte{})
		var query string
		if reqOpt.Signed {
			_query, err := json.Marshal(reqOpt.Query)
			if err != nil {
				return nil, err
			}
			query = common.B2s(_query)
			query = "signed_body=SIGNATURE." + url.QueryEscape(query)
		} else {
			vurl := url.Values{}
			for key, vul := range reqOpt.Query {
				vurl.Set(key, fmt.Sprintf("%v", vul))
			}
			query = vurl.Encode()
		}

		if reqOpt.IsPost {
			bf.WriteString(query)
		} else {
			_url.RawQuery = query
		}
	} else if reqOpt.Body != nil {
		bf = reqOpt.Body
	} else {
		bf = bytes.NewBuffer([]byte{})
	}

	var req *http.Request
	req, err = http.NewRequest(method, _url.String(), bf)
	if err != nil {
		return nil, err
	}

	this.setHeader(reqOpt, req)
	for key, vul := range reqOpt.Header {
		SetHeader(req, key, vul)
	}

	resp, err := this.c.Do(req)
	if err != nil {
		return nil, &common.MakeMoneyError{
			ErrType:   common.RequestError,
			ExternErr: err,
		}
	}

	defer resp.Body.Close()
	this.afterRequest(_url, resp)

	encoding := resp.Header.Get("Content-Encoding")

	var body io.Reader
	switch encoding {
	case "gzip":
		body, err = gzip.NewReader(resp.Body)
		break
	case "zstd":
		body, err = zstd.NewReader(resp.Body)
		break
	case "deflate":
	case "":
		body = resp.Body
		break
	}

	if err != nil {
		return nil, &common.MakeMoneyError{
			ErrType:   common.RequestError,
			ExternErr: err,
		}
	}

	bodyData, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, &common.MakeMoneyError{
			ErrType:   common.RequestError,
			ExternErr: err,
		}
	}

	return bodyData, err
}

func truncation(body []byte) []byte {
	if body == nil {
		return []byte("body is nil!")
	}
	if config.UseTruncation {
		if len(body) > 100 {
			return body[:100]
		}
	}
	return body
}

func (this *Instagram) CheckInstReqError(url string, body []byte, err error) {
	if err != nil {
		log.Error("account: %s, url: %s, request error: %v, resp: %s", this.User, url, err, truncation(body))
		this.ReqErrorCount += 1
		if common.IsError(err, common.RequestError) {
			this.ReqContError++
		}
	} else {
		if config.IsDebug {
			log.Info("account: %s, url: %s, api resp %s", this.User, url, truncation(body))
		}
		this.ReqContError = 0
	}
}

func (this *Instagram) HttpRequest(reqOpt *reqOptions) ([]byte, error) {
	body, err := this.httpDo(reqOpt)
	this.CheckInstReqError(reqOpt.ApiPath, body, err)
	return body, err
}

func (this *Instagram) HttpRequestJson(reqOpt *reqOptions, response interface{}) (err error) {
	body, err := this.httpDo(reqOpt)
	if err == nil {
		err = json.Unmarshal(body, &response)
	}
	this.CheckInstReqError(reqOpt.ApiPath, body, err)

	if response != nil {
		value := reflect.ValueOf(response)
		if value.CanInterface() {
			setInfo := value.MethodByName("SetInfo")
			if setInfo.Kind() == reflect.Func {
				setInfo.Call([]reflect.Value{reflect.ValueOf(reqOpt.ApiPath), reflect.ValueOf(this.User)})
			} else {
				log.Warn("reflect SetInfo error! url: %s", reqOpt.ApiPath)
			}
		} else {
			log.Warn("reflect SetInfo error! url: %s", reqOpt.ApiPath)
		}
	}
	return err
}
