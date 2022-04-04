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
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

type reqOptions struct {
	ApiPath        string
	IsPost         bool
	IsApiB         bool
	IsApiGraph     bool
	Signed         bool
	Query          map[string]interface{}
	Body           *bytes.Buffer
	Json           interface{}
	Header         map[string]string
	HeaderSequence *HeaderSequence
	DisAutoHeader  bool
	RawApiPath     string
}

func SetHeader(req *http.Request, key string, vul string) {
	//req.Header[key] = []string{vul}
	req.Header.Set(key, vul)
}

func (this *Instagram) setHeader(reqOpt *reqOptions, req *http.Request) {
	var seq *HeaderSequence
	if reqOpt.HeaderSequence == nil {
		var headerMap map[string]*HeaderSequence

		if this.IsLogin {
			headerMap = LoginHeaderMap
		} else {
			headerMap = NoLoginHeaderMap
		}
		ApiPathKey := reqOpt.ApiPath
		if reqOpt.RawApiPath != "" {
			ApiPathKey = reqOpt.RawApiPath
		}
		seq = headerMap[ApiPathKey]
		if seq == nil {
			log.Error("api path: %s has no header map!", reqOpt.ApiPath)
		}
		req.HeaderSequence = seq.HeaderSeq
		req.OnlySequence = true
	} else {
		seq = reqOpt.HeaderSequence
		req.HeaderSequence = reqOpt.HeaderSequence.HeaderSeq
		req.OnlySequence = true
	}

	for _, fun := range seq.HeaderFun {
		fun(this, reqOpt, req)
	}

	if common.IsDebug {
		for _, header := range seq.HeaderSeq {
			if req.Header.Get(header) == "" && header != "Content-Length" {
				//log.Warn("api path: %s, header: %s is null", reqOpt.ApiPath, header)
			}
		}
	}
}

func (this *Instagram) afterRequest(reqUrl *url.URL, resp *http.Response) {
	_url, _ := url.Parse(InstagramHost)
	for _, value := range this.c.Jar.Cookies(_url) {
		if strings.Contains(value.Name, "csrftoken") {
			this.token = value.Value
		}
	}

	for key := range resp.Header {
		header := ""
		value := ""
		if strings.Index(key, "Ig-Set-") == 0 {
			header = key[len("Ig-Set-"):]
			value = resp.Header.Get(key)
		} else if strings.Index(key, "X-Ig-Set-") == 0 {
			header = "X-Ig-" + key[len("X-Ig-Set-"):]
			value = resp.Header.Get(key)
		}
		if header != "" {
			this.httpHeader[header] = value
			if common.IsDebug {
				//log.Info("account: %s set header %s = %s", this.User, header, value)
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
			//reqOpt.Query["_uuid"] = this.Device.DeviceID
			//reqOpt.Query["_uid"] = this.ID
		}

		if reqOpt.Signed {
			_query, err := json.Marshal(reqOpt.Query)
			if err != nil {
				return nil, err
			}
			//query = strings.ReplaceAll(common.B2s(_query), "\\\\", "\\") //for password
			//query = "signed_body=SIGNATURE." + common.InstagramQueryEscape(query)
			query = "signed_body=SIGNATURE." + common.InstagramQueryEscape(common.B2s(_query))
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
	} else if reqOpt.Json != nil {
		bf = bytes.NewBuffer([]byte{})
		var query string
		_query, err := json.Marshal(reqOpt.Json)
		if err != nil {
			return nil, err
		}
		query = common.B2s(_query)
		if reqOpt.Signed {
			query = strings.ReplaceAll(query, "\\\\", "\\") //for password
			query = "signed_body=SIGNATURE." + common.InstagramQueryEscape(query)
		}

		if reqOpt.IsPost {
			bf.WriteString(query)
		} else {
			log.Warn("url %s get not allow Json", reqOpt.ApiPath)
		}
	} else {
		bf = bytes.NewBuffer([]byte{})
	}

	var req *http.Request
	req, err = http.NewRequest(method, _url.String(), bf)
	if err != nil {
		return nil, err
	}

	if !reqOpt.DisAutoHeader && !reqOpt.IsApiGraph {
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
	if !reqOpt.DisAutoHeader && !reqOpt.IsApiGraph {
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
	if common.UseTruncation {
		if len(body) > 100 {
			return body[:100]
		}
	}
	return body
}

func (this *Instagram) CheckInstReqError(url string, body []byte, err error) {
	if err != nil {
		log.Error("account: %s, url: %s, request error: %v, resp: %s", this.User, url, err, truncation(body))
		//this.ReqErrorCount += 1
		if common.IsError(err, common.RequestError) {
			//this.ReqContError++
		}
	} else {
		if common.IsDebug {
			log.Info("account: %s, url: %s, api resp %s", this.User, url, truncation(body))
		}
		//this.ReqContError = 0
	}
}

func (this *Instagram) PrepareProxy() error {
	var id = ""
	if this.Proxy != nil {
		id = this.Proxy.ID
	}
	proxy, errProxy := ProxyCallBack(this.AccountInfo.Register.RegisterIpCountry, id)
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

	//for true {
	for i := 0; i < 3; i++ {
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

	//for true {
	for i := 0; i < 3; i++ {
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
