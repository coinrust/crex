package serve

import (
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	. "github.com/coinrust/crex"
	"github.com/coinrust/crex/exchanges"
)

var (
	configFile string
)

type SConfig struct {
	Exchanges []SExchange            `toml:"exchange"`
	Options   map[string]interface{} `toml:"option"`
}

type SExchange struct {
	Name      string `toml:"name"`
	DebugMode bool   `toml:"debug_mode"`
	AccessKey string `toml:"access_key"`
	SecretKey string `toml:"secret_key"`
	Testnet   bool   `toml:"testnet"`
	WebSocket bool   `toml:"websocket"`
}

// Serve 加载策略并执行
func Serve(strategy Strategy) (err error) {
	flag.StringVar(&configFile, "c", "config.toml", "")
	flag.Parse()

	err = strategy.SetSelf(strategy)
	if err != nil {
		return
	}

	err = SetupStrategyFromConfig(strategy)
	if err != nil {
		return
	}

	err = strategy.OnInit()
	if err != nil {
		return
	}
	err = strategy.Run()
	if err != nil {
		return
	}
	err = strategy.OnExit()
	return
}

// SetupStrategyFromConfig 根据配置文件设置策略参数
func SetupStrategyFromConfig(strategy Strategy) (err error) {
	c := SConfig{}
	if _, err = toml.DecodeFile(configFile, &c); err != nil {
		return
	}
	if len(c.Exchanges) == 0 {
		err = fmt.Errorf("no exchange found")
		return
	}
	var exs []Exchange
	for _, ex := range c.Exchanges {
		exchange := exchanges.NewExchange(ex.Name,
			ApiDebugModeOption(ex.DebugMode),
			ApiAccessKeyOption(ex.AccessKey),
			ApiSecretKeyOption(ex.SecretKey),
			ApiTestnetOption(ex.Testnet),
			ApiWebSocketOption(ex.WebSocket))
		exs = append(exs, exchange)
	}
	if err = strategy.Setup(TradeModeLiveTrading, exs...); err != nil {
		return
	}
	//log.Printf("options: %#v", options)
	err = strategy.SetOptions(c.Options)
	return
}
