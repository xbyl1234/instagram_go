package proxy

import (
	"crypto/tls"
	"fmt"
	"golang.org/x/net/http2"
	"golang.org/x/net/proxy"
	"makemoney/common"
	"makemoney/common/log"
	"net/http"
	"net/url"
)

type ProxyType int
type BlackType int

var (
	ProxyHttp   ProxyType = 0
	ProxySocket ProxyType = 1

	BlacktypeNoblack      BlackType = 0
	BlacktypeRisk         BlackType = 1
	BlacktypeConn         BlackType = 2
	BlacktypeRegisterrisk BlackType = 3
)

type Proxy struct {
	ID              string    `json:"id"`
	Ip              string    `json:"ip"`
	Port            string    `json:"port"`
	Username        string    `json:"username"`
	Passwd          string    `json:"passwd"`
	Rip             string    `json:"rip"`
	ProxyType       ProxyType `json:"proxy_type"`
	NeedAuth        bool      `json:"need_auth"`
	Country         string    `json:"country"`
	IsUsed          bool      `json:"is_used"`
	IsBusy          bool      `json:"is_busy"`
	RegisterSuccess int       `json:"register_success"`
	RegisterError   int       `json:"register_error"`
	BlackType       BlackType `json:"black_type"`
}

func (this *Proxy) GetProxy() *http.Transport {
	var tr *http.Transport

	if this.ProxyType == 0 {
		var proxyUrl string
		if this.NeedAuth {
			proxyUrl = "http://" + this.Username + ":" + this.Passwd + "@" + this.Ip + ":" + this.Port
		} else {
			proxyUrl = "http://" + this.Ip + ":" + this.Port
		}
		_url, _ := url.Parse(proxyUrl)
		tr = &http.Transport{
			Proxy:           http.ProxyURL(_url),
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	} else {
		var auth = &proxy.Auth{}
		if this.NeedAuth {
			auth.User = this.Username
			auth.Password = this.Passwd
		} else {
			auth = nil
		}
		dialer, _ := proxy.SOCKS5("tcp", this.Ip+":"+this.Port, auth, proxy.Direct)
		tr = &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
		tr.Dial = dialer.Dial
	}

	tr.TLSClientConfig = &tls.Config{
		NextProtos: []string{"h2", "h2-fb", "http/1.1"},
		MinVersion: tls.VersionTLS13,
		MaxVersion: tls.VersionTLS13,
		CipherSuites: []uint16{
			tls.TLS_AES_128_GCM_SHA256,
			tls.TLS_AES_256_GCM_SHA384,
			tls.TLS_CHACHA20_POLY1305_SHA256,
		},
	}
	err := http2.ConfigureTransport(tr)
	if err != nil {
		log.Error("ConfigureTransport error: %v", err)
		return nil
	}
	return tr
}

type ProxyImpl interface {
	GetNoRisk(busy bool, used bool) *Proxy
	Get(id string) *Proxy
	Black(proxy *Proxy, _type BlackType)
	Remove(proxy *Proxy)
	Dumps()
}

type ProxyConfigt struct {
	Providers []struct {
		ProviderName string    `json:"provider"`
		Url          string    `json:"url"`
		Country      string    `json:"country"`
		ProxyType    ProxyType `json:"proxy_type"`
	} `json:"providers"`
}

type ProxyPoolt struct {
	proxys map[string]ProxyImpl
}

func (this *ProxyPoolt) Get(country string, id string) *Proxy {
	if country == "" {
		for key := range this.proxys {
			country = key
			break
		}
	}

	p := this.proxys[country]
	if p == nil {
		return nil
	}
	return p.Get(id)
}

func (this *ProxyPoolt) GetNoRisk(country string, busy bool, used bool) *Proxy {
	if country == "" {
		for key := range this.proxys {
			country = key
			break
		}
	}
	p := this.proxys[country]
	if p == nil {
		return nil
	}
	return p.GetNoRisk(busy, used)
}

var ProxyPool ProxyPoolt
var ProxyConfig ProxyConfigt

func InitProxyPool(configPath string) error {
	err := common.LoadJsonFile(configPath, &ProxyConfig)
	if err != nil || len(ProxyConfig.Providers) == 0 {
		log.Error("load proxy config error: %v", err)
		return err
	}

	ProxyPool.proxys = make(map[string]ProxyImpl)
	for _, provider := range ProxyConfig.Providers {
		var _proxy ProxyImpl
		var err error

		switch provider.ProviderName {
		case "dove":
			_proxy, err = InitDovePool(provider.Url)
			break
		case "luminati":
			_proxy, err = InitLuminatiPool(provider.Url)
			break
		default:
			return &common.MakeMoneyError{ErrStr: fmt.Sprintf("proxy config provider error: %s",
				provider.ProviderName), ErrType: common.OtherError}
		}
		if err != nil {
			return &common.MakeMoneyError{ErrStr: fmt.Sprintf("proxy config provider error: %s",
				provider.ProviderName), ErrType: common.OtherError}
		}

		ProxyPool.proxys[provider.Country] = _proxy
	}
	return err
}
