package proxys

import (
	"makemoney/common"
	"makemoney/common/log"
	"net/http"
	url2 "net/url"
	"strings"
	"sync"
	"time"
)

type RolaPool struct {
	ProxyImpl
	url         string
	proxyType   common.ProxyType
	lastReqTime time.Time
	lock        sync.Mutex
	proxyList   []*common.Proxy
	proxyMask   []bool
	client      *http.Client
	Country     string
}

func InitRolaPool(url string) (ProxyImpl, error) {
	var pool = &RolaPool{}
	purl, err := url2.Parse(url)
	if err != nil {
		return nil, err
	}

	for key, value := range purl.Query() {
		if key == "protocol" {
			if value[0] == "socks5" {
				pool.proxyType = common.ProxySocket
			} else {
				pool.proxyType = common.ProxyHttp
			}
		} else if key == "country" {
			pool.Country = value[0]
		}
	}

	pool.url = url
	pool.client = &http.Client{}
	//common.DebugHttpClient(pool.client)
	return pool, nil
}

type RolaResp struct {
	Code int      `json:"code"`
	Data []string `json:"data"`
	Msg  string   `json:"msg"`
}

func (this *RolaPool) RequestProxy() bool {
	resp := &RolaResp{}
	for true {
		err := common.HttpDoJson(this.client, &common.RequestOpt{
			ReqUrl: this.url,
			IsPost: false,
		}, resp)
		if err != nil {
			log.Error("rola proxy request error: %v", err)
			time.Sleep(1 * time.Second)
			continue
		}
		if resp.Code != 0 {
			log.Error("rola proxy request error: %v", resp.Msg)
			time.Sleep(3 * time.Second)
			continue
		}

		break
	}
	if len(resp.Data) == 0 {
		log.Error("rola request proxy list is null!")
		return false
	}
	this.proxyList = make([]*common.Proxy, len(resp.Data))
	this.proxyMask = make([]bool, len(resp.Data))

	for index := range resp.Data {
		dp := resp.Data[index]
		sp := strings.Split(dp, ":")
		this.proxyList[index] = &common.Proxy{
			ID:        "rola",
			Ip:        sp[0],
			Port:      sp[1],
			Username:  "",
			Passwd:    "",
			Rip:       sp[0],
			ProxyType: this.proxyType,
			NeedAuth:  false,
			Country:   this.Country,
		}
		this.proxyMask[index] = true
	}

	log.Info("rola request proxy list success!")
	return true
}

func (this *RolaPool) get() *common.Proxy {
	this.lock.Lock()
	defer this.lock.Unlock()
	index := 0
	find := false
	for index = range this.proxyMask {
		if this.proxyMask[index] {
			find = true
			break
		}
	}

	if find {
		this.proxyMask[index] = false
		return this.proxyList[index]
	}
	if this.RequestProxy() {
		this.proxyMask[0] = false
		return this.proxyList[0]
	} else {
		return nil
	}
}

func (this *RolaPool) GetNoRisk(busy bool, used bool) *common.Proxy {
	return this.get()
}

func (this *RolaPool) Get(id string) *common.Proxy {
	return this.get()
}