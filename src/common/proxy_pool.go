package common

import (
	"crypto/tls"
	"golang.org/x/net/proxy"
	"makemoney/common/log"
	math_rand "math/rand"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

type ProxyType int
type BlackType int

var (
	ProxyHttp   ProxyType = 0
	ProxySocket ProxyType = 1

	BlackType_NoBlack      BlackType = 0
	BlackType_Risk         BlackType = 1
	BlackType_Conn         BlackType = 2
	BlackType_RegisterRisk BlackType = 3
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

type ProxyPoolt struct {
	allCount  int
	allProxys map[string]*Proxy
	ProxyList []*Proxy
	proxyLock sync.Mutex
	path      string
	dumpsPath string
}

var ProxyPool ProxyPoolt

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
		var auth *proxy.Auth = &proxy.Auth{}
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

func InitProxyPool(path string) error {
	ProxyPool.path = path
	ProxyPool.dumpsPath = strings.ReplaceAll(path, ".json", "_dumps.json")
	if PathExists(ProxyPool.dumpsPath) {
		path = ProxyPool.dumpsPath
	}

	var ProxyMap map[string]*Proxy
	err := LoadJsonFile(path, &ProxyMap)
	if err != nil {
		return err
	}

	ProxyList := make([]*Proxy, len(ProxyMap))
	var index = 0
	for _, vul := range ProxyMap {
		if vul.BlackType != BlackType_NoBlack {
			continue
		}
		ProxyList[index] = vul
		index++
	}
	if len(ProxyList) == 0 {
		return &MakeMoneyError{ErrStr: "no proxy", ErrType: PorxyError}
	}

	ProxyPool.ProxyList = ProxyList[:index]
	ProxyPool.allCount = len(ProxyMap)
	ProxyPool.allProxys = ProxyMap
	return nil
}

func (this *ProxyPoolt) GetNoRisk(busy bool, used bool) *Proxy {
	this.proxyLock.Lock()
	defer this.proxyLock.Unlock()
	for true {
		math_rand.Seed(time.Now().Unix())
		index := math_rand.Intn(len(this.ProxyList))
		if this.ProxyList[index].BlackType == BlackType_NoBlack {
			if busy {
				if this.ProxyList[index].IsBusy {
					continue
				}
				this.ProxyList[index].IsBusy = true
			}

			if used {
				if this.ProxyList[index].IsUsed {
					continue
				}
				this.ProxyList[index].IsUsed = true
			}

			return this.ProxyList[index]
		}
	}
	return nil
}

func (this *ProxyPoolt) Get(id string) *Proxy {
	this.proxyLock.Lock()
	defer this.proxyLock.Unlock()
	return this.allProxys[id]
}

func (this *ProxyPoolt) Black(proxy *Proxy, _type BlackType) {
	this.proxyLock.Lock()
	defer this.proxyLock.Unlock()
	proxy.BlackType = _type
	this.remove(proxy)
	this.Dumps()
}

func (this *ProxyPoolt) Remove(proxy *Proxy) {
	this.proxyLock.Lock()
	defer this.proxyLock.Unlock()
	this.remove(proxy)
}

func (this *ProxyPoolt) remove(proxy *Proxy) {
	find := false
	var index int
	for index = range this.ProxyList {
		if this.ProxyList[index] == proxy {
			find = true
			break
		}
	}
	if find {
		delete(this.allProxys, proxy.ID)
		this.ProxyList = append(this.ProxyList[:index], this.ProxyList[index+1:]...)
	}
}

func (this *ProxyPoolt) Dumps() {
	err := Dumps(this.dumpsPath, this.allProxys)
	if err != nil {
		log.Error("dumps proxy pool error:%v", err)
	}
}
