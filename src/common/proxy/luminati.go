package proxy

import (
	"crypto/tls"
	"golang.org/x/net/proxy"
	"makemoney/common"
	"makemoney/common/log"
	math_rand "math/rand"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

type LuminatiPool struct {
	ProxyPoolt
	allCount  int
	allProxys map[string]*Proxy
	ProxyList []*Proxy
	proxyLock sync.Mutex
	path      string
	dumpsPath string
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

func InitLuminatiPool(path string) (ProxyPoolt, error) {
	var pool = &LuminatiPool{}
	pool.path = path
	pool.dumpsPath = strings.ReplaceAll(path, ".json", "_dumps.json")
	if common.PathExists(pool.dumpsPath) {
		path = pool.dumpsPath
	}

	var ProxyMap map[string]*Proxy
	err := common.LoadJsonFile(path, &ProxyMap)
	if err != nil {
		return nil, err
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
		return nil, &common.MakeMoneyError{ErrStr: "no proxy", ErrType: common.PorxyError}
	}

	pool.ProxyList = ProxyList[:index]
	pool.allCount = len(ProxyMap)
	pool.allProxys = ProxyMap
	return pool, nil
}

func (this *LuminatiPool) GetNoRisk(busy bool, used bool) *Proxy {
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

func (this *LuminatiPool) Get(id string) *Proxy {
	this.proxyLock.Lock()
	defer this.proxyLock.Unlock()
	return this.allProxys[id]
}

func (this *LuminatiPool) Black(proxy *Proxy, _type BlackType) {
	this.proxyLock.Lock()
	defer this.proxyLock.Unlock()
	proxy.BlackType = _type
	this.remove(proxy)
	this.Dumps()
}

func (this *LuminatiPool) Remove(proxy *Proxy) {
	this.proxyLock.Lock()
	defer this.proxyLock.Unlock()
	this.remove(proxy)
}

func (this *LuminatiPool) remove(proxy *Proxy) {
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

func (this *LuminatiPool) Dumps() {
	err := common.Dumps(this.dumpsPath, this.allProxys)
	if err != nil {
		log.Error("dumps proxy pool error:%v", err)
	}
}
