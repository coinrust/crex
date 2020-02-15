package data

import (
	"testing"
	"time"
)

func TestNewDeribitData(t *testing.T) {
	filename := `D:\trading\deribit\deribit_btc_perpetual_and_futures_tick_by_tick_book_snapshots_10_levels\deribit_BTC-27MAR20_and_futures_tick_by_tick_book_snapshots_10_levels_2019-10-01_2019-11-01.csv`
	data := NewDeribitData(filename)
	data.Next()
}

func TestLoadDeribitTickByTickBookSnapshots(t *testing.T) {
	filename := `D:\trading\deribit\deribit_btc_perpetual_and_futures_tick_by_tick_book_snapshots_10_levels\deribit_BTC-PERPETUAL_and_futures_tick_by_tick_book_snapshots_10_levels_2019-10-01_2019-11-01.csv`
	data, err := LoadDeribitTickByTickBookSnapshots(filename)
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("%v", data.Len())
	for {
		time.Sleep(time.Second)
	}

	//for _, item := range data {
	//	t.Logf("%#v", item)
	//}
}
