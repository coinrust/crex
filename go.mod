module github.com/coinrust/gotrader

go 1.13

require (
	github.com/chuckpreslar/emission v0.0.0-20170206194824-a7ddd980baf9
	github.com/frankrap/bitmex-api v0.0.0-20200225011535-9765a3267676
	github.com/frankrap/bybit-api v0.0.0-20200316091251-fc44b96f58e8
	github.com/frankrap/deribit-api v0.0.0-20200211002849-4210ce6a675a
	github.com/frankrap/huobi-api v0.0.0-00010101000000-000000000000
	github.com/go-echarts/go-echarts v0.0.0-20190915064101-cbb3b43ade5d
	github.com/gobuffalo/packr v1.30.1 // indirect
	github.com/nntaoli-project/GoEx v1.0.7
	github.com/pkg/errors v0.9.1 // indirect
	github.com/sony/sonyflake v1.0.0
	github.com/spf13/viper v1.6.2
	golang.org/x/sys v0.0.0-20200124204421-9fbb57f87de9 // indirect
)

replace github.com/frankrap/huobi-api => ../../../github.com/frankrap/huobi-api
