package dataloader

import (
	"testing"
	"time"
)

func TestMongoDBDataLoader_Setup(t *testing.T) {
	loader := NewMongoDBDataLoader("mongodb://localhost:27017",
		"tick_db", "deribit", "BTC-PERPETUAL")
	start, _ := time.Parse("2006-01-02 15:04:05", "2019-10-01 00:00:00")
	end, _ := time.Parse("2006-01-02 15:04:05", "2019-10-02 00:00:00")
	loader.Setup(start, end)

	t.Log("------------")
	s := time.Now()
	result := loader.ReadOrderBooks()
	//for _, v := range result {
	//	t.Logf("%v", v.Time)
	//}
	t.Logf("%v", len(result))
	t.Logf("%v", time.Now().Sub(s).String())

	t.Log("------------")
	s = time.Now()
	result = loader.ReadOrderBooks()
	//for _, v := range result {
	//	t.Logf("%v", v.Time)
	//}
	t.Logf("%v", len(result))
	t.Logf("%v", time.Now().Sub(s).String())
	last := result[len(result)-1]
	t.Logf("%v", last.Time.String())
}
