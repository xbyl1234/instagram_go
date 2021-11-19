package tools

import (
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"makemoney/config"
	"net/http"
	neturl "net/url"
	"os"
)

func LoadJsonFile(path string, ret interface{}) error {
	by, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	err = json.Unmarshal(by, ret)
	return err
}

func Dumps(path string, obj interface{}) error {
	by, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	_, err = file.Write(by)
	if err != nil {
		return err
	}
	_ = file.Close()
	return nil
}

func DebugHttpClient(clinet *http.Client) {
	if config.UseCharles {
		uri, err := neturl.Parse("http://127.0.0.1:8888")
		if err == nil {
			clinet.Transport = &http.Transport{
				Proxy: http.ProxyURL(uri),
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			}
		}
	}
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}
