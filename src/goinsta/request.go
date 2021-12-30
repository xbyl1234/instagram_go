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
	url  string
	inst *Instagram

	Status     string `json:"status"`
	ErrorType  string `json:"error_type"`
	Message    string `json:"message"`
	ErrorTitle string `json:"error_title"`
}

func (this *BaseApiResp) SetInfo(url string, inst *Instagram) {
	this.inst = inst
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
			this.inst.User,
			this.url,
			this.ErrorType+":"+this.Message)
		if this.Message == InsAccountError_ChallengeRequired {
			this.inst.Status = InsAccountError_ChallengeRequired
			return &common.MakeMoneyError{ErrStr: this.Message, ErrType: common.ChallengeRequiredError}
		} else if this.Message == InsAccountError_LoginRequired {
			this.inst.Status = InsAccountError_LoginRequired
			return &common.MakeMoneyError{ErrStr: this.Message, ErrType: common.LoginRequiredError}
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
	//req.Header.Set("accept-encoding", "zstd, gzip, deflate")
	//SetHeader(req, "x-ig-android-id", this.androidID)
	//igwwwClaim := this.ReadHeader(IGHeader_igwwwClaim)
	//if igwwwClaim == "" {
	//	igwwwClaim = "0"
	//}
	//SetHeader(req, IGHeader_igwwwClaim, igwwwClaim)
	//SetHeader(req, "x-bloks-is-layout-rtl", "false")
	//SetHeader(req, "x-ig-connection-speed", fmt.Sprintf("%dkbps", common.GenNumber(1000, 3700)))
	//SetHeader(req, "x-ig-bandwidth-totalbytes-b", "0")
	//SetHeader(req, "x-ig-bandwidth-totaltime-ms", "0")

	req.Header.Set("connection", "keep-alive")
	SetHeader(req, "ig-intended-user-id", strconv.FormatInt(this.ID, 10))
	SetHeader(req, "x-ig-connection-speed", "38kbps")
	SetHeader(req, "x-ig-device-id", this.deviceID)
	SetHeader(req, "x-ig-timezone-offset", "-28800")
	SetHeader(req, "x-ig-capabilities", "36r/Fx8=")
	SetHeader(req, "x-pigeon-rawclienttime", strconv.FormatInt(time.Now().Unix(), 10)+".000000")
	SetHeader(req, "x-ig-device-locale", InstagramLocation)
	SetHeader(req, "X-Ig-Abr-Connection-Speed-Kbps", "0")
	SetHeader(req, "x-ig-family-device-id", this.familyID)
	req.Header.Set("accept-language", "en-US;q=1.0")
	if req.Header.Get("content-type") == "" {
		SetHeader(req, "content-type", "application/x-www-form-urlencoded; charset=UTF-8")
	}
	req.Header.Set("user-agent", this.UserAgent)
	SetHeader(req, "x-ig-app-locale", InstagramLocation)
	SetHeader(req, "x-ig-bandwidth-speed-kbps", "0.000")
	SetHeader(req, "x-ig-mapped-locale", InstagramLocation)
	SetHeader(req, IGHeader_XMid, this.ReadHeader(IGHeader_XMid))
	SetHeader(req, "x-bloks-is-panorama-enabled", "true")
	SetHeader(req, "x-bloks-version-id", InstagramBloksVersionID)
	SetHeader(req, "x-pigeon-session-id", this.sessionID)
	SetHeader(req, "x-ig-app-id", InstagramAppID)
	SetHeader(req, "x-ig-connection-type", "WIFI")
	SetHeader(req, "X-Tigon-Is-Retry", "False")
	req.Header.Set("accept-encoding", "zstd, gzip, deflate")
	SetHeader(req, "x-fb-http-engine", "Liger")
	SetHeader(req, "x-fb-client-ip", "True")
	SetHeader(req, "x-fb-server-cluster", "True")
}

func (this *Instagram) setLoginHeader(req *http.Request) {
	//SetHeader(req, IGHeader_udsUserID, strconv.FormatInt(this.ID, 10))
	//SetHeader(req, IGHeader_iguRur, this.ReadHeader(IGHeader_iguRur))
	//SetHeader(req, IGHeader_Authorization, this.ReadHeader(IGHeader_Authorization))
	//SetHeader(req, "x-ig-app-startup-country", "OR")
}

func (this *Instagram) setHeader(reqOpt *reqOptions, req *http.Request) {
	this.setBaseHeader(req)
	if this.IsLogin {
		this.setLoginHeader(req)
	}
	for key := range this.httpHeader {
		SetHeader(req, key, this.httpHeader[key])
	}
	//SetHeader(req,"x-ads-opt-out", "0")
	//SetHeader(req,"x-cm-latency", "-1.000")
	//SetHeader(req,"x-ig-extended-cdn-thumbnail-cache-busting-value", "1000")
}

func (this *Instagram) afterRequest(reqUrl *url.URL, resp *http.Response) {
	_url, _ := url.Parse(InstagramHost)
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
		baseUrl = InstagramHost_B
	} else if reqOpt.IsApiGraph {
		baseUrl = InstagramHost_Graph
	} else {
		baseUrl = InstagramHost
	}

	_url, err := url.Parse(baseUrl + reqOpt.ApiPath)
	if err != nil {
		return nil, err
	}

	var bf *bytes.Buffer
	if reqOpt.Query != nil {
		bf = bytes.NewBuffer([]byte{})
		var query string

		if reqOpt.IsPost && this.IsLogin && !reqOpt.IsApiGraph {
			if this.token != "" {
				reqOpt.Query["_csrftoken"] = this.token
			}
			reqOpt.Query["_uuid"] = this.deviceID
			reqOpt.Query["_uid"] = this.ID
		}

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

	if reqOpt.DisAutoHeader && !reqOpt.IsApiGraph {
		this.setHeader(reqOpt, req)
	}
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
	if reqOpt.DisAutoHeader && !reqOpt.IsApiGraph {
		this.afterRequest(_url, resp)
	}
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

func (this *Instagram) PrepareProxy() error {
	var id = ""
	if this.Proxy != nil {
		id = this.Proxy.ID
	}
	proxy, errProxy := ProxyCallBack(id)
	if errProxy != nil {
		log.Error("account: %s, get proxy error: %v", this.User, errProxy)
		return errProxy
	}
	this.SetProxy(proxy)
	return nil
}

func (this *Instagram) HttpRequest(reqOpt *reqOptions) ([]byte, error) {
	if this.Proxy == nil || this.Proxy.Ip == "" {
		err := this.PrepareProxy()
		if err != nil {
			return nil, err
		}
	}

	for true {
		result, err := this._HttpRequest(reqOpt)
		if common.IsError(err, common.RequestError) {
			log.Warn("account: %s,url: %s, request error: %v...", this.User, reqOpt.ApiPath, err)
			errProxy := this.PrepareProxy()
			if errProxy != nil {
				return nil, err
			}
			continue
		}
		return result, err
	}
	return nil, nil
}

func (this *Instagram) HttpRequestJson(reqOpt *reqOptions, response interface{}) error {
	if this.Proxy == nil || this.Proxy.Ip == "" {
		err := this.PrepareProxy()
		if err != nil {
			return err
		}
	}

	for true {
		err := this._HttpRequestJson(reqOpt, response)
		if common.IsError(err, common.RequestError) {
			//log.Warn("account: %s,url: %s, request error: %v...", this.User, reqOpt.ApiPath, err)
			errProxy := this.PrepareProxy()
			if errProxy != nil {
				return err
			}
			continue
		}
		return err
	}
	return nil
}

func (this *Instagram) _HttpRequest(reqOpt *reqOptions) ([]byte, error) {
	body, err := this.httpDo(reqOpt)
	this.CheckInstReqError(reqOpt.ApiPath, body, err)
	return body, err
}

func (this *Instagram) _HttpRequestJson(reqOpt *reqOptions, response interface{}) (err error) {
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
				setInfo.Call([]reflect.Value{reflect.ValueOf(reqOpt.ApiPath), reflect.ValueOf(this)})
			} else {
				log.Warn("reflect SetInfo error! url: %s", reqOpt.ApiPath)
			}
		} else {
			log.Warn("reflect SetInfo error! url: %s", reqOpt.ApiPath)
		}
	}
	return err
}
