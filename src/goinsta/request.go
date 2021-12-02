package goinsta

import (
	"bytes"
	"encoding/json"
	"fmt"
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
	ApiPath   string
	IsPost    bool
	IsApiB    bool
	Signed    bool
	Query     map[string]interface{}
	Body      *bytes.Buffer
	HeaderKey []string
	Header    map[string]string
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
		return &common.MakeMoneyError{ErrStr: this.Message, ErrType: common.ApiError}
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
)

func (this *Instagram) setBaseHeader(req *http.Request) {
	req.Header.Set("connection", "keep-alive")
	req.Header.Set("accept-language", "en-US")
	req.Header.Set("user-agent", goInstaUserAgent)
	req.Header.Set("x-ig-app-id", fbAnalytics)
	req.Header.Set("x-ig-capabilities", igCapabilities)
	req.Header.Set("x-ig-connection-type", connType)
	req.Header.Set("x-fb-client-ip", "True")
	req.Header.Set("x-fb-http-engine", "Liger")
	req.Header.Set("x-fb-server-cluster", "True")
	req.Header.Set("accept-encoding", "deflate")
	req.Header.Set("x-ig-family-device-id", this.familyID)

	if req.Header.Get("content-type") == "" {
		req.Header.Set("content-type", "application/x-www-form-urlencoded; charset=UTF-8")
	}
	if this.ReadHeader(IGHeader_Authorization) != "" {
		req.Header.Set(IGHeader_Authorization, this.ReadHeader(IGHeader_Authorization))
	}
	if this.ReadHeader(IGHeader_iguRur) != "" {
		req.Header.Set(IGHeader_iguRur, this.ReadHeader(IGHeader_iguRur))
	}
	//if this.ReadHeader(IGHeader_XMid) != "" {
	//	req.Header.Set(IGHeader_iguRur, this.ReadHeader(IGHeader_iguRur))
	//}

	if this.IsLogin {
		req.Header.Set("ig-intended-user-id", strconv.FormatInt(this.ID, 10))
		req.Header.Set("ig-u-ds-user-id", strconv.FormatInt(this.ID, 10))
	}
}

func (this *Instagram) setHeader(reqOpt *reqOptions, req *http.Request) {
	this.setBaseHeader(req)
	req.Header.Set("x-ig-connection-speed", fmt.Sprintf("%dkbps", common.GenNumber(1000, 3700)))
	req.Header.Set("x-ig-bandwidth-speed-kbps", "-1.000")
	req.Header.Set("x-ig-bandwidth-totalbytes-b", "0")
	req.Header.Set("x-ig-bandwidth-totaltime-ms", "0")

	req.Header.Set("x-ads-opt-out", "0")
	req.Header.Set("x-cm-latency", "-1.000")
	req.Header.Set("x-ig-app-locale", "en_US")
	req.Header.Set("x-ig-device-locale", "en_US")
	req.Header.Set("x-pigeon-session-id", common.GenUUID())
	req.Header.Set("x-pigeon-rawclienttime", strconv.FormatInt(time.Now().Unix(), 10))
	req.Header.Set("x-ig-extended-cdn-thumbnail-cache-busting-value", "1000")
	req.Header.Set("x-ig-device-id", this.uuid)
	req.Header.Set("x-ig-android-id", this.androidID)

	for index := range reqOpt.HeaderKey {
		key := reqOpt.HeaderKey[index]
		req.Header.Set(key, this.ReadHeader(key))
	}
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

			if IGHeader_udsUserID == setting[len("ig-set-"):] {
				this.ID, _ = strconv.ParseInt(resp.Header.Get(key), 10, 64)
			}
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
		req.Header.Set(key, vul)
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

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, &common.MakeMoneyError{
			ErrType:   common.RequestError,
			ExternErr: err,
		}
	}

	return body, err
}

func truncation(body []byte) []byte {
	if body == nil {
		return []byte("body is nil!")
	}
	if len(body) > 100 {
		return body[:100]
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
