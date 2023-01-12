package modules

import (
	"Proxy/components/utils"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"h12.io/socks"
)

func GetHttpTransport(Proxy string) *http.Transport {
	ProxyUrl, err := url.Parse(Proxy)
	if utils.HandleError(err) {
		return &http.Transport{}
	}

	return &http.Transport{
		Proxy: http.ProxyURL(ProxyUrl),
	}
}

func GetSocksTransport(Proxy string) *http.Transport {
	return &http.Transport{
		Dial: socks.Dial(fmt.Sprintf("%s?timeout=%ds", Proxy, utils.Config.Filter.Timeout)),
	}
}

func GetTransport(Proxy string) *http.Transport {
	if strings.Contains(Proxy, "http://") {
		return GetHttpTransport(Proxy)
	} else {
		return GetSocksTransport(Proxy)
	}
}

func ProxyReq(req string, proxy string) (res *http.Response, err error) {
	ReqUrl, err := url.Parse(req)
	if utils.HandleError(err) {
		return nil, err
	}

	client := &http.Client{
		Timeout:   time.Duration(time.Duration(utils.Config.Filter.Timeout) * time.Second),
		Transport: GetTransport(proxy),
	}

	res, err = client.Get(ReqUrl.String())
	return res, err
}

func CheckProxy(Proxy string) {
	response, err := ProxyReq("https://api.ipify.org", Proxy)

	if err != nil {
		if utils.Config.Options.ShowDeadProxies {
			utils.Log(fmt.Sprintf("[DEAD]  %s", Proxy))
		}
		utils.Dead++
		return
	}

	defer response.Body.Close()

	content, err := ioutil.ReadAll(response.Body)
	if utils.HandleError(err) {
		return
	}

	// Check if the proxy is "transparent"
	is_elite := string(content) != utils.ActualIp

	prox := strings.Split(Proxy, "://")

	utils.Log(fmt.Sprintf("[ALIVE] [ELITE: %v] [%s] %s", is_elite, prox[0], prox[1]))
	utils.Valid++

	switch prox[0] {
		case "http":
			utils.Http++
		case "socks4":
			utils.Socks4++
		case "socks5":
			utils.Socks5++
	}

	if !is_elite && !utils.Config.Options.SaveTransparent {
		return
	}

	utils.AppendFile("checked.txt", Proxy)
}
