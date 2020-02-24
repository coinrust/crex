package deribit_broker

import (
	"github.com/frankrap/deribit-api"
	"github.com/frankrap/deribit-api/models"
	"log"
	"testing"
	"time"
)

func TestDiribitBroker_GetOrderBook(t *testing.T) {
	apiKey := "AsJTU16U"
	secretKey := "mM5_K8LVxztN6TjjYpv_cJVGQBvk4jglrEpqkw1b87U"
	b := NewBroker(deribit.TestBaseURL, apiKey, secretKey)
	b.GetOrderBook("BTC-PERPETUAL", 10)
}

func TestDiribitBroker_Subscribe(t *testing.T) {
	apiKey := "AsJTU16U"
	secretKey := "mM5_K8LVxztN6TjjYpv_cJVGQBvk4jglrEpqkw1b87U"
	b := NewBroker(deribit.TestBaseURL, apiKey, secretKey)
	//event := "book.ETH-PERPETUAL.100.1.100ms"
	event := "book.BTC-PERPETUAL.100ms"
	b.Subscribe(event, "", func(e *models.OrderBookNotification) {
		log.Printf("event ---: ")
	})

	for {
		time.Sleep(1 * time.Second)
	}
}