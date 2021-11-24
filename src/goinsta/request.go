package goinsta

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"makemoney/common"
	"makemoney/common/log"
	"makemoney/config"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type reqOptions struct {
	Login     bool
	ApiPath   string
	IsPost    bool
	IsApiB    bool
	Signed    bool
	Query     map[string]interface{}
	HeaderKey []string
}

type sendOptions struct {
	Url       string
	IsPost    bool
	Body      *bytes.Buffer
	HeaderKey []string
	Header    map[string]string
}

type BaseApiResp struct {
	Status    string `json:"status"`
	ErrorType string `json:"error_type"`
	Message   string `json:"message"`
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

func (insta *Instagram) sendSimpleRequest(uri string, a ...interface{}) (body []byte, err error) {
	return insta.HttpRequest(
		&reqOptions{
			ApiPath: fmt.Sprintf(uri, a...),
		},
	)
}

func (insta *Instagram) setBaseHeader(req *http.Request) {
	req.Header.Set("connection", "keep-alive")
	req.Header.Set("content-type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("accept-language", "en-US")
	req.Header.Set("user-agent", goInstaUserAgent)
	req.Header.Set("x-ig-app-id", fbAnalytics)
	req.Header.Set("x-ig-capabilities", igCapabilities)
	req.Header.Set("x-ig-connection-type", connType)
	req.Header.Set("x-fb-client-ip", "True")
	req.Header.Set("x-fb-http-engine", "Liger")
	req.Header.Set("x-fb-server-cluster", "True")
	req.Header.Set("accept-encoding", "deflate")
	req.Header.Set("x-ig-family-device-id", insta.familyID)

	if insta.ReadHeader(IGHeader_Authorization) != "" {
		req.Header.Set(IGHeader_Authorization, insta.ReadHeader(IGHeader_Authorization))
	}
	if insta.ReadHeader(IGHeader_iguRur) != "" {
		req.Header.Set(IGHeader_iguRur, insta.ReadHeader(IGHeader_iguRur))
	}
	if insta.IsLogin {
		req.Header.Set("ig-intended-user-id", insta.id)
		req.Header.Set("ig-u-ds-user-id", insta.id)
	}
}

func (insta *Instagram) setHeader(reqOpt *reqOptions, req *http.Request) {
	insta.setBaseHeader(req)
	req.Header.Set("x-ig-connection-speed", fmt.Sprintf("%dkbps", acquireRand(1000, 3700)))
	req.Header.Set("x-ig-bandwidth-speed-kbps", "-1.000")
	req.Header.Set("x-ig-bandwidth-totalbytes-b", "0")
	req.Header.Set("x-ig-bandwidth-totaltime-ms", "0")

	req.Header.Set("x-ads-opt-out", "0")
	req.Header.Set("x-cm-latency", "-1.000")
	req.Header.Set("x-ig-app-locale", "en_US")
	req.Header.Set("x-ig-device-locale", "en_US")
	req.Header.Set("x-pigeon-session-id", common.GenerateUUID())
	req.Header.Set("x-pigeon-rawclienttime", strconv.FormatInt(time.Now().Unix(), 10))
	req.Header.Set("x-ig-extended-cdn-thumbnail-cache-busting-value", "1000")
	req.Header.Set("x-ig-device-id", insta.uuid)
	req.Header.Set("x-ig-android-id", insta.androidID)

	for index := range reqOpt.HeaderKey {
		key := reqOpt.HeaderKey[index]
		req.Header.Set(key, insta.ReadHeader(key))
	}
}

func (insta *Instagram) afterRequest(reqUrl *url.URL, resp *http.Response) {
	_url, _ := url.Parse(goInstaAPIUrl)
	for _, value := range insta.c.Jar.Cookies(_url) {
		if strings.Contains(value.Name, "csrftoken") {
			insta.token = value.Value
		}
	}

	for key := range resp.Header {
		setting := strings.ToLower(key)
		if strings.Index(setting, "ig-set-") == 0 {
			insta.httpHeader[setting[len("ig-set-"):]] = resp.Header.Get(key)

			if IGHeader_udsUserID == setting[len("ig-set-"):] {
				insta.id = resp.Header.Get(key)
			}
		}
	}
}

func (insta *Instagram) httpDo(reqOpt *reqOptions) ([]byte, error) {
	method := "GET"
	if reqOpt.IsPost {
		method = "POST"
	}

	var baseUrl string
	if reqOpt.IsApiB {
		baseUrl = goInstaAPIUrl_B
	} else {
		baseUrl = goInstaAPIUrl
	}

	_url, err := url.Parse(baseUrl + reqOpt.ApiPath)
	if err != nil {
		return nil, err
	}

	bf := bytes.NewBuffer([]byte{})

	var query string
	_query, err := json.Marshal(reqOpt.Query)
	if err != nil {
		return nil, err
	}
	query = common.B2s(_query)

	if reqOpt.Signed {
		query = "signed_body=SIGNATURE." + url.QueryEscape(query)
	} else {
		query = url.QueryEscape(query)
	}

	if reqOpt.IsPost {
		bf.WriteString(query)
	} else {
		_url.RawQuery = query
	}

	var req *http.Request
	req, err = http.NewRequest(method, _url.String(), bf)
	if err != nil {
		return nil, err
	}

	insta.setHeader(reqOpt, req)

	resp, err := insta.c.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	insta.afterRequest(_url, resp)

	body, err := ioutil.ReadAll(resp.Body)
	if err == nil {
		err = isError(resp.StatusCode, body)
	}
	return body, err
}

func (insta *Instagram) CheckInstReqError(url string, body []byte, err error) {
	var hadLog = false
	defer func() {
		if !hadLog {
			insta.ReqSuccessCount += 1
		}

		if config.IsDebug && !hadLog {
			log.Info("account: %s, url: %s, api resp %s", insta.User, url, body)
		}
	}()

	if err != nil {
		log.Warn("account: %s, url: %s, request error: %v", insta.User, url, err)
		insta.ReqErrorCount += 1
		hadLog = true
		return
	}

	resp := &BaseApiResp{}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		log.Warn("account: %s, url: %s, Unmarshal error %s", insta.User, url, body)
		insta.ReqApiErrorCount += 1
		hadLog = true
	} else {
		if resp.isError() {
			log.Warn("account: %s, url: %s, api error: %s", insta.User, url, body)
			insta.ReqApiErrorCount += 1
			hadLog = true
		}
	}
}

func (insta *Instagram) HttpRequest(reqOpt *reqOptions) ([]byte, error) {
	body, err := insta.httpDo(reqOpt)
	insta.CheckInstReqError(reqOpt.ApiPath, body, err)
	return body, err
}

func (insta *Instagram) HttpRequestJson(reqOpt *reqOptions, response interface{}) (err error) {
	body, err := insta.httpDo(reqOpt)
	insta.CheckInstReqError(reqOpt.ApiPath, body, err)

	err = json.Unmarshal(body, &response)
	return err
}

func (insta *Instagram) HttpSend(sendOpt *sendOptions, response interface{}) ([]byte, error) {
	var req *http.Request
	_url, err := url.Parse(sendOpt.Url)
	if err != nil {
		return nil, err
	}

	var body *bytes.Buffer
	var method string

	if sendOpt.IsPost {
		method = "POST"
		body = sendOpt.Body
	} else {
		method = "GET"
		body = bytes.NewBuffer([]byte{})
	}

	req, err = http.NewRequest(method, _url.String(), body)
	if err != nil {
		return nil, err
	}

	insta.setBaseHeader(req)
	for key, vul := range sendOpt.Header {
		req.Header.Set(key, vul)
	}

	resp, err := insta.c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	insta.afterRequest(_url, resp)
	if response != nil {
		err = json.Unmarshal(respBody, response)
	}
	insta.CheckInstReqError(sendOpt.Url, respBody, err)
	return respBody, err
}

//{"message":"Please wait a few minutes before you try again.","status":"fail"}
func isError(code int, body []byte) (err error) {
	switch code {
	case 200:
	case 503:
		return Error503{
			Message: "Instagram API error. Try it later.",
		}
	case 400:
		ierr := Error400{}
		err = json.Unmarshal(body, &ierr)
		if err != nil {
			return err
		}

		if ierr.Message == "challenge_required" {
			return ierr.ChallengeError
		}

		if err == nil && ierr.Message != "" {
			return ierr
		}
	default:
		ierr := ErrorN{}
		err = json.Unmarshal(body, &ierr)
		if err != nil {
			return err
		}
		return ierr
	}
	return nil
}

func (insta *Instagram) prepareData(other ...map[string]interface{}) (string, error) {
	data := map[string]interface{}{
		"_uuid":      insta.uuid,
		"_csrftoken": insta.token,
	}
	if insta.Account != nil && insta.Account.ID != 0 {
		data["_uid"] = strconv.FormatInt(insta.Account.ID, 10)
	}

	for i := range other {
		for key, value := range other[i] {
			data[key] = value
		}
	}
	b, err := json.Marshal(data)
	if err == nil {
		return common.B2s(b), err
	}
	return "", err
}

func (insta *Instagram) prepareDataQuery(other ...map[string]interface{}) map[string]interface{} {
	data := map[string]interface{}{
		"_uuid":      insta.uuid,
		"_csrftoken": insta.token,
	}
	if insta.Account != nil && insta.Account.ID != 0 {
		data["_uid"] = strconv.FormatInt(insta.Account.ID, 10)
	}
	for i := range other {
		for key, value := range other[i] {
			data[key] = toString(value)
		}
	}
	return data
}

func acquireRand(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}
