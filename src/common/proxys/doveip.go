package proxys

import (
	"fmt"
	"makemoney/common"
	"makemoney/common/log"
	"net/http"
	url2 "net/url"
	"sync"
	"time"
)

type DovePool struct {
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

//https://dvapi.doveip.com/cmapi.php?rq=distribute&user=safrg534524&token=RklaOXZUbFp6WTFZQjBSUzVJTkt4QT09&auth=1&geo=US&city=all&agreement=0&timeout=35&num=10
func InitDovePool(url string) (ProxyImpl, error) {
	var pool = &DovePool{}
	purl, err := url2.Parse(url)
	if err != nil {
		return nil, err
	}

	for key, value := range purl.Query() {
		if key == "agreement" {
			if value[0] == "0" {
				pool.proxyType = common.ProxySocket
			} else {
				pool.proxyType = common.ProxyHttp
			}
		} else if key == "geo" {
			pool.Country = value[0]
		}
	}

	pool.url = url
	pool.client = &http.Client{}
	//common.DebugHttpClient(pool.client)
	return pool, nil
}

type DoveResp struct {
	Errno int    `json:"errno"`
	Msg   string `json:"msg"`
	Data  []struct {
		Geo      string `json:"geo"`
		Ip       string `json:"ip"`
		Port     int    `json:"port"`
		DIp      string `json:"d_ip"`
		Timeout  int    `json:"timeout"`
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"data"`
}

func (this *DovePool) RequestProxy() bool {
	resp := &DoveResp{}
	for true {
		err := common.HttpDoJson(this.client, &common.RequestOpt{
			ReqUrl: this.url,
			IsPost: false,
		}, resp)
		if err != nil {
			log.Error("dove proxy request error: %v", err)
			time.Sleep(1 * time.Second)
			continue
		}
		if resp.Errno == 409 {
			log.Error("dove proxy request frequently...")
			time.Sleep(3 * time.Second)
			continue
		}
		if resp.Errno == 403 {
			log.Error("dove proxy no money...")
			return false
		}
		if resp.Errno != 200 {
			log.Error("dove proxy request error: %d", resp.Errno)
			return false
		}
		break
	}
	if len(resp.Data) == 0 {
		log.Error("doveip request proxy list is null!")
		return false
	}
	this.proxyList = make([]*common.Proxy, len(resp.Data))
	this.proxyMask = make([]bool, len(resp.Data))

	for index := range resp.Data {
		dp := &resp.Data[index]
		this.proxyList[index] = &common.Proxy{
			ID:        "dove",
			Ip:        dp.Ip,
			Port:      fmt.Sprintf("%d", dp.Port),
			Username:  dp.Username,
			Passwd:    dp.Password,
			Rip:       dp.DIp,
			ProxyType: this.proxyType,
			NeedAuth:  true,
			Country:   this.Country,
		}
		this.proxyMask[index] = true
	}

	log.Info("doveip request proxy list success!")
	return true
}

func (this *DovePool) get() *common.Proxy {
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

func (this *DovePool) GetNoRisk(busy bool, used bool) *common.Proxy {
	return this.get()
}

func (this *DovePool) Get(id string) *common.Proxy {
	return this.get()
}
