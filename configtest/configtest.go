package configtest

import (
	"github.com/BurntSushi/toml"
	"log"
	"path/filepath"
)

var (
	testDataDir = "../../testdata/"
)

type Config struct {
	BinanceFutures TestConfig `toml:"binancefutures"`
	Bitmex         TestConfig `toml:"bitmex"`
	Bybit          TestConfig `toml:"bybit"`
	Deribit        TestConfig `toml:"deribit"`
	Hbdm           TestConfig `toml:"hbdm"`
	HbdmSwap       TestConfig `toml:"hbdmswap"`
	OkexFutures    TestConfig `toml:"okexfutures"`
	OkexSwap       TestConfig `toml:"okexswap"`
}

type TestConfig struct {
	AccessKey  string `toml:"access_key"`
	SecretKey  string `toml:"secret_key"`
	Passphrase string `toml:"passphrase"`
	Testnet    bool   `toml:"testnet"`
	ProxyURL   string `toml:"proxy_url"`
}

func LoadTestConfig(name string) *TestConfig {
	fPath := filepath.Join(testDataDir, "configtest.toml")
	var cfg Config
	if _, err := toml.DecodeFile(fPath, &cfg); err != nil {
		log.Panic(err)
	}
	tCfg := &TestConfig{}
	switch name {
	case "binancefutures":
		tCfg = &cfg.BinanceFutures
	case "bitmex":
		tCfg = &cfg.Bitmex
	case "bybit":
		tCfg = &cfg.Bybit
	case "deribit":
		tCfg = &cfg.Deribit
	case "hbdm":
		tCfg = &cfg.Hbdm
	case "hbdmswap":
		tCfg = &cfg.HbdmSwap
	case "okexfutures":
		tCfg = &cfg.OkexFutures
	case "okexswap":
		tCfg = &cfg.OkexSwap
	}
	return tCfg
}
