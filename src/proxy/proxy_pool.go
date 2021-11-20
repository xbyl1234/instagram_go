package proxy

import (
	"container/list"
	"crypto/tls"
	"golang.org/x/net/proxy"
	"makemoney/log"
	"makemoney/tools"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

type Proxy struct {
	Id          string    `json:"id"`
	Ip          string    `json:"ip"`
	Port        string    `json:"port"`
	Username    string    `json:"username"`
	Passwd      string    `json:"passwd"`
	Rip         string    `json:"rip"`
	ProxyType   ProxyType `json:"proxy_type"`
	NeedAuth    bool      `json:"need_auth"`
	Country     string    `json:"country"`
	IsUsed      bool      `json:"is_used"`
	IsConnError bool      `json:"is_conn_error"`
	IsRisk      bool      `json:"is_risk"`
}

func (this *Proxy) GetProxy() *http.Transport {
	if this.ProxyType == 0 {
		var proxyUrl string
		if this.NeedAuth {
			proxyUrl = "http://" + this.Username + ":" + this.Passwd + "@" + this.Ip + ":" + this.Port
		} else {
			proxyUrl = "http://" + this.Ip + ":" + this.Port
		}
		_url, _ := url.Parse(proxyUrl)
		return &http.Transport{
			Proxy:           http.ProxyURL(_url),
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	} else {
		var auth *proxy.Auth
		if this.NeedAuth {
			auth.User = this.Username
			auth.Password = this.Passwd
		} else {
			auth = nil
		}
		dialer, _ := proxy.SOCKS5("tcp", this.Ip+":"+this.Port, auth, proxy.Direct)
		var httpTran = &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
		httpTran.Dial = dialer.Dial
		return httpTran
	}
}

type ProxyType int

var (
	ProxyHttp   ProxyType = 0
	ProxySocket ProxyType = 1
)

type ProxyPoolt struct {
	allCount          int
	proxysAvailable   *list.List
	proxyNotAvailable *list.List
	allProxys         map[string]Proxy
	repeIndex         int
	proxyLock         sync.Mutex
	path              string
	dumpsPath         string
}

var ProxyPool ProxyPoolt

func InitProxyPool(path string) error {
	ProxyPool.path = path
	ProxyPool.dumpsPath = strings.ReplaceAll(path, ".json", "_dumps.json")
	if tools.PathExists(ProxyPool.dumpsPath) {
		path = ProxyPool.dumpsPath
	}

	var ProxyList map[string]Proxy
	err := tools.LoadJsonFile(path, &ProxyList)
	if err != nil {
		return err
	}

	ProxyPool.proxysAvailable = list.New()
	for _, vul := range ProxyList {
		if vul.IsRisk || vul.IsConnError || vul.IsUsed {
			continue
		}
		ProxyPool.proxysAvailable.PushBack(&vul)
	}

	ProxyPool.allCount = len(ProxyList)
	ProxyPool.allProxys = ProxyList
	return nil
}

func (this *ProxyPoolt) GetOne() (*Proxy, error) {
	this.proxyLock.Lock()
	defer this.proxyLock.Unlock()

	if this.proxysAvailable.Len() != 0 {
		ret := this.proxysAvailable.Front().Value.(*Proxy)
		this.proxysAvailable.Remove(this.proxysAvailable.Front())
		return ret, nil
	}
	return nil, &tools.MakeMoneyError{"no available proxy", 0}
}

//func (this *ProxyPoolt) GetRepeOne() (*Proxy, error) {
//
//}

func (this *ProxyPoolt) BlackConnErrorProxy(proxy *Proxy) {
	proxy.IsConnError = true
}

func (this *ProxyPoolt) BlackRiskErrorProxy(proxy *Proxy) {
	proxy.IsRisk = true
}

func (this *ProxyPoolt) SetUsedProxy(proxy *Proxy) {
	proxy.IsUsed = true
}

func (this *ProxyPoolt) Dumps() {
	err := tools.Dumps(this.dumpsPath, this.allProxys)
	if err != nil {
		log.Error("dumps proxy pool error:%v", err)
	}
}
