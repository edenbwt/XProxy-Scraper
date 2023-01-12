package utils

import (
	"io/ioutil"
	"net/http"

	"github.com/BurntSushi/toml"
)

var (
	Config   = ConfigStruct{}
	ActualIp string
	Valid    = 0

	Socks4 = 0
	Socks5 = 0
	Http   = 0
	Dead   = 0
)

type ConfigStruct struct {
	Filter struct {
		Timeout int  `toml:"timeout"`
		Http    bool `toml:"http"`
		Socks4  bool `toml:"socks4"`
		Socks5  bool `toml:"socks5"`
	} `toml:"filter"`

	Dev struct {
		Debug bool `toml:"debug"`
	} `toml:"dev"`

	Options struct {
		Scrape              bool `toml:"scrape"`
		Threads             int  `toml:"threads"`
		ScrapeThreads       int  `toml:"scrape_threads"`
		SaveTransparent     bool `toml:"save_transparent"`
		ShowDeadProxies     bool `toml:"show_dead_proxies"`
		RemoveUrlOnError    bool `toml:"remove_url_on_error"`
		ScrapeTimeout       int  `toml:"scrape_timeout"`
		CheckScrapedProxies bool `toml:"check_scraped_proxies"`
	} `toml:"options"`
}

func GetActualIp() string {
	res, err := http.Get("https://api.ipify.org")
	if HandleError(err) {
		return ""
	}

	defer res.Body.Close()

	content, err := ioutil.ReadAll(res.Body)
	if HandleError(err) {
		return ""
	}

	return string(content)
}

func LoadConfig() {
	if _, err := toml.DecodeFile("script/config.toml", &Config); err != nil {
		panic(err)
	}
	ActualIp = GetActualIp()
}
