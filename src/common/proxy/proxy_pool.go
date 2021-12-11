package proxy

import (
	"fmt"
	"makemoney/common"
	"makemoney/common/log"
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

type ProxyPoolt interface {
	GetNoRisk(busy bool, used bool) *Proxy
	Get(id string) *Proxy
	Black(proxy *Proxy, _type BlackType)
	Remove(proxy *Proxy)
}

var ProxyPool ProxyPoolt

type ProxyConfigt struct {
	Provider  string    `json:"provider"`
	Url       string    `json:"url"`
	ProxyType ProxyType `json:"proxy_type"`
}

var ProxyConfig ProxyConfigt

func InitProxyPool(configPath string) error {
	err := common.LoadJsonFile(configPath, &ProxyConfig)
	if err != nil {
		log.Error("load proxy config error: %v", err)
		return err
	}

	switch ProxyConfig.Provider {
	case "dove":
		ProxyPool, err = InitDovePool(ProxyConfig.Url)
		break
	case "luminati":
		ProxyPool, err = InitLuminatiPool(ProxyConfig.Url)
		break
	default:
		return &common.MakeMoneyError{ErrStr: fmt.Sprintf("proxy config provider error: %s",
			ProxyConfig.Provider), ErrType: common.OtherError}
	}
	return err
}
