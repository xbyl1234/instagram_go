package goinsta

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type reqOptions struct {
	// Connection is connection header. Default is "close".
	//Connection string

	// Login process
	Login bool

	// Endpoint is the request path of instagram api
	Endpoint string

	// IsPost set to true will send request with POST method.
	//
	// By default this option is false.
	IsPost bool

	// UseV2 is set when API endpoint uses v2 url.
	UseV2 bool

	Signed bool
	// Query is the parameters of the request
	//
	// This parameters are independents of the request method (POST|GET)
	Query map[string]string
}

func (insta *Instagram) sendSimpleRequest(uri string, a ...interface{}) (body []byte, err error) {
	return insta.sendRequest(
		&reqOptions{
			Endpoint: fmt.Sprintf(uri, a...),
		},
	)
}

func (insta *Instagram) setHeader(reqOpt *reqOptions, req *http.Request) {
	req.Header.Set("connection", "keep-alive")
	req.Header.Set("content-type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("accept-language", "en-US")
	req.Header.Set("user-agent", goInstaUserAgent)
	req.Header.Set("x-ig-app-id", fbAnalytics)
	req.Header.Set("x-ig-capabilities", igCapabilities)
	req.Header.Set("x-ig-connection-type", connType)
	req.Header.Set("x-ig-connection-speed", fmt.Sprintf("%dkbps", acquireRand(1000, 3700)))
	req.Header.Set("x-ig-bandwidth-speed-kbps", "-1.000")
	req.Header.Set("x-ig-bandwidth-totalbytes-b", "0")
	req.Header.Set("x-ig-bandwidth-totaltime-ms", "0")

	req.Header.Set("x-ads-opt-out", "0")
	req.Header.Set("x-cm-latency", "-1.000")
	req.Header.Set("x-ig-app-locale", "en_US")
	req.Header.Set("x-ig-device-locale", "en_US")
	req.Header.Set("x-pigeon-session-id", generateUUID())
	req.Header.Set("x-pigeon-rawclienttime", strconv.FormatInt(time.Now().Unix(), 10))
	req.Header.Set("x-ig-extended-cdn-thumbnail-cache-busting-value", "1000")
	req.Header.Set("x-ig-device-id", insta.uuid)
	req.Header.Set("x-ig-android-id", insta.dID)
	req.Header.Set("x-fb-http-engine", "Liger")
}

func (insta *Instagram) afterRequest(reqUrl *url.URL, resp *http.Response) {
	_url, _ := url.Parse(goInstaAPIUrl)
	for _, value := range insta.c.Jar.Cookies(_url) {
		if strings.Contains(value.Name, "csrftoken") {
			insta.token = value.Value
		}
	}

	encryptionId := resp.Header.Get("ig-set-password-encryption-key-id")
	if encryptionId != "" {
		insta.encryptionId = encryptionId
	}
	encryptionKey := resp.Header.Get("ig-set-password-encryption-pub-key")
	if encryptionKey != "" {
		insta.encryptionKey = encryptionKey
	}
	authorization := resp.Header.Get("ig-set-authorization")
	if authorization != "" {
		insta.authorization = authorization
	}
}

func (insta *Instagram) SendRequest(reqOpt *reqOptions, response interface{}) (err error) {
	method := "GET"
	if reqOpt.IsPost {
		method = "POST"
	}

	nu := goInstaAPIUrl
	if reqOpt.UseV2 {
		nu = goInstaAPIUrlv2
	}

	_url, err := url.Parse(nu + reqOpt.Endpoint)
	if err != nil {
		return err
	}

	vs := url.Values{}
	bf := bytes.NewBuffer([]byte{})

	for k, v := range reqOpt.Query {
		vs.Add(k, v)
	}

	if reqOpt.IsPost {
		if reqOpt.Signed {
			bf.WriteString(generateSignature(vs.Encode()))
		} else {
			bf.WriteString(vs.Encode())
		}
	} else {
		for k, v := range _url.Query() {
			vs.Add(k, strings.Join(v, " "))
		}
		_url.RawQuery = vs.Encode()
	}

	var req *http.Request
	req, err = http.NewRequest(method, _url.String(), bf)
	if err != nil {
		return
	}

	insta.setHeader(reqOpt, req)

	resp, err := insta.c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	insta.afterRequest(_url, resp)

	body, err := ioutil.ReadAll(resp.Body)
	if err == nil {
		err = isError(resp.StatusCode, body)
	}
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &response)
	return err
}

func (insta *Instagram) sendRequest(o *reqOptions) (body []byte, err error) {
	method := "GET"
	if o.IsPost {
		method = "POST"
	}

	nu := goInstaAPIUrl
	if o.UseV2 {
		nu = goInstaAPIUrlv2
	}

	u, err := url.Parse(nu + o.Endpoint)
	if err != nil {
		return nil, err
	}

	vs := url.Values{}
	bf := bytes.NewBuffer([]byte{})

	for k, v := range o.Query {
		vs.Add(k, v)
	}

	if o.IsPost {
		bf.WriteString(vs.Encode())
	} else {
		for k, v := range u.Query() {
			vs.Add(k, strings.Join(v, " "))
		}

		u.RawQuery = vs.Encode()
	}

	var req *http.Request
	req, err = http.NewRequest(method, u.String(), bf)
	if err != nil {
		return
	}

	insta.setHeader(o, req)

	resp, err := insta.c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	u, _ = url.Parse(goInstaAPIUrl)
	for _, value := range insta.c.Jar.Cookies(u) {
		if strings.Contains(value.Name, "csrftoken") {
			insta.token = value.Value
		}
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err == nil {
		err = isError(resp.StatusCode, body)
	}
	return body, err
}

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
		return b2s(b), err
	}
	return "", err
}

func (insta *Instagram) prepareDataQuery(other ...map[string]interface{}) map[string]string {
	data := map[string]string{
		"_uuid":      insta.uuid,
		"_csrftoken": insta.token,
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
