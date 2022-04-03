package proxys

import (
	"fmt"
	"makemoney/common"
	"makemoney/common/log"
)

type ProxyImpl interface {
	GetNoRisk(busy bool, used bool) *common.Proxy
	Get(id string) *common.Proxy
	Remove(proxy *common.Proxy)
	Dumps()
}

type ProxyConfigt struct {
	Providers []struct {
		ProviderName string           `json:"provider"`
		Url          string           `json:"url"`
		Country      string           `json:"country"`
		ProxyType    common.ProxyType `json:"proxy_type"`
	} `json:"providers"`
}

type ProxyPoolt struct {
	proxys map[string]ProxyImpl
}

func (this *ProxyPoolt) Get(country string, id string) *common.Proxy {
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

func (this *ProxyPoolt) GetNoRisk(country string, busy bool, used bool) *common.Proxy {
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
		//case "luminati":
		//	_proxy, err = InitLuminatiPool(provider.Url)
		//	break
		case "idea":
			_proxy, err = InitIdeaPool(provider.Url)
			break
		case "rola":
			_proxy, err = InitRolaPool(provider.Url)
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
