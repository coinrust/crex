module github.com/coinrust/gotrader

go 1.13

require (
	github.com/chuckpreslar/emission v0.0.0-20170206194824-a7ddd980baf9
	github.com/frankrap/bitmex-api v0.0.0-20200214075706-b1926307d51c
	github.com/frankrap/bybit-api v0.0.0-20200305072436-e3a4b902648e
	github.com/frankrap/deribit-api v0.0.0-20200211002849-4210ce6a675a
	github.com/go-echarts/go-echarts v0.0.0-20190915064101-cbb3b43ade5d
	github.com/gobuffalo/packr v1.30.1 // indirect
	github.com/sony/sonyflake v1.0.0
)

replace github.com/frankrap/bybit-api => ../../../github.com/frankrap/bybit-api
