package configtest

import (
	"github.com/spf13/viper"
	"log"
)

type TestConfig struct {
	AccessKey  string
	SecretKey  string
	Passphrase string
	Testnet    bool
	ProxyURL   string
}

func LoadTestConfig(name string) *TestConfig {
	viper.SetConfigName("configtest")
	viper.AddConfigPath("../../testdata")
	err := viper.ReadInConfig()
	if err != nil {
		log.Panic(err)
	}
	cfg := &TestConfig{}
	items := viper.GetStringMapString(name)
	if len(items) == 0 {
		log.Panic("test config not found [configtest.yaml]")
	}
	for key, value := range items {
		switch key {
		case "access_key":
			cfg.AccessKey = value
		case "secret_key":
			cfg.SecretKey = value
		case "passphrase":
			cfg.Passphrase = value
		case "testnet":
			cfg.Testnet = value == "true"
		case "proxy_url":
			cfg.ProxyURL = value
		}
	}
	return cfg
}
