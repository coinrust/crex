package serve

import (
	"flag"
	"fmt"
	. "github.com/coinrust/crex"
	"github.com/coinrust/crex/exchanges"
	"github.com/spf13/viper"
	"path/filepath"
	"strings"
)

var (
	configFile string
)

type ExchangeConfig struct {
	Exchanges []ExchangeItem `yaml:"exchanges"`
}

type ExchangeItem struct {
	Name       string `yaml:"name"`
	Access_Key string `yaml:"access_key"`
	Secret_Key string `yaml:"secret_key"`
	Testnet    bool   `yaml:"testnet"`
	WebSocket  bool   `yaml:"websocket"`
}

// Serve 加载策略并执行
func Serve(strategy Strategy) (err error) {
	flag.StringVar(&configFile, "c", "config.yaml", "")
	flag.Parse()

	base := filepath.Base(configFile)
	ext := filepath.Ext(configFile)
	var configType string
	if strings.HasPrefix(ext, ".") {
		configType = ext[1:]
	} else {
		err = fmt.Errorf("wrong configuration file")
		return
	}

	viper.SetConfigType(configType)
	viper.SetConfigName(base)
	viper.AddConfigPath(".")

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

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
	err = strategy.OnDeinit()
	return
}

// SetupStrategyFromConfig 根据配置文件设置策略参数
func SetupStrategyFromConfig(strategy Strategy) (err error) {
	c := ExchangeConfig{}
	err = viper.Unmarshal(&c)
	if err != nil {
		return
	}
	if len(c.Exchanges) == 0 {
		err = fmt.Errorf("no exchange found")
		return
	}
	var exs []Exchange
	for _, ex := range c.Exchanges {
		exchange := exchanges.NewExchange(ex.Name,
			ApiAccessKeyOption(ex.Access_Key),
			ApiSecretKeyOption(ex.Secret_Key),
			ApiTestnetOption(ex.Testnet),
			ApiWebSocketOption(ex.WebSocket))
		exs = append(exs, exchange)
	}
	err = strategy.Setup(TradeModeLiveTrading, exs...)
	if err != nil {
		return
	}
	options := viper.GetStringMap("options")
	//log.Printf("options: %#v", options)
	err = strategy.SetOptions(options)
	return
}
