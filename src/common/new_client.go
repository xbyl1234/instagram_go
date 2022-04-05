package common

import (
	"crypto/tls"
	"fmt"
	"golang.org/x/net/http2"
	"golang.org/x/net/proxy"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

type ProxyType int
type HttpConfigFun func(c *http.Client)

const (
	ProxyHttp   ProxyType = 0
	ProxySocket ProxyType = 1
)

type Proxy struct {
	ID        string        `json:"id"`
	Ip        string        `json:"ip"`
	Port      string        `json:"port"`
	Username  string        `json:"username"`
	Passwd    string        `json:"passwd"`
	Rip       string        `json:"rip"`
	ProxyType ProxyType     `json:"proxy_type"`
	NeedAuth  bool          `json:"need_auth"`
	Country   string        `json:"country"`
	LiveTime  time.Duration `json:"live_time"`
	StartTime time.Time
}

func (this *Proxy) IsOutLiveTime() bool {
	if time.Duration(float64(time.Since(this.StartTime))*0.5) > this.LiveTime {
		return true
	}
	return false
}

func (this *Proxy) GetProxyUrl() string {
	url := ""
	if this.ProxyType == ProxyHttp {
		url += "http://"
	} else {
		url += "socks5://"
	}
	if this.NeedAuth {
		url += this.Username + ":" + this.Passwd + "@"
	}
	url += this.Ip
	url += ":"
	url += this.Port
	return url
}

func (this *Proxy) GetProxy() HttpConfigFun {
	return func(c *http.Client) {
		if this == nil {
			return
		}
		tr := c.Transport.(*http.Transport)
		if this.ProxyType == ProxyHttp {
			var proxyUrl string
			if this.NeedAuth {
				proxyUrl = "http://" + this.Username + ":" + this.Passwd + "@" + this.Ip + ":" + this.Port
			} else {
				proxyUrl = "http://" + this.Ip + ":" + this.Port
			}
			_url, _ := url.Parse(proxyUrl)
			tr.Proxy = http.ProxyURL(_url)
		} else {
			var auth = &proxy.Auth{}
			if this.NeedAuth {
				auth.User = this.Username
				auth.Password = this.Passwd
			} else {
				auth = nil
			}
			dialer, _ := proxy.SOCKS5("tcp", this.Ip+":"+this.Port, auth, proxy.Direct)
			tr.Dial = dialer.Dial
		}

		err := http2.ConfigureTransport(tr)
		if err != nil {
			fmt.Printf("new client %v \n", err)
		}
	}
}

func DisableHttpSslPinng() HttpConfigFun {
	return func(c *http.Client) {
		tr := c.Transport.(*http.Transport)
		tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
}

func DisableRedirect() HttpConfigFun {
	return func(c *http.Client) {
		c.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}
}
func NeedJar() HttpConfigFun {
	return func(c *http.Client) {
		jar, _ := cookiejar.New(nil)
		c.Jar = jar
	}
}

func DefaultHttpTimeout() HttpConfigFun {
	return func(c *http.Client) {
		tr := c.Transport.(*http.Transport)
		tr.ForceAttemptHTTP2 = true
		tr.MaxIdleConns = 100
		tr.IdleConnTimeout = 90 * time.Second
		tr.TLSHandshakeTimeout = 10 * time.Second
		tr.ExpectContinueTimeout = 1 * time.Second
	}
}

func CreateGoHttpClient(httpConfigs ...HttpConfigFun) *http.Client {
	tr := &http.Transport{
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	httpClient := &http.Client{
		Transport: tr,
	}
	for _, config := range httpConfigs {
		config(httpClient)
	}
	if UseCharles {
		DefaultHttpProxy.GetProxy()(httpClient)
		DisableHttpSslPinng()(httpClient)
	}
	return httpClient
}
