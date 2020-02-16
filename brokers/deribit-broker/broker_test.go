package deribit_broker

import (
	"github.com/frankrap/deribit-api"
	"testing"
)

func TestDiribitBroker_GetOrderBook(t *testing.T) {
	apiKey := "AsJTU16U"
	secretKey := "mM5_K8LVxztN6TjjYpv_cJVGQBvk4jglrEpqkw1b87U"
	b := NewBroker(deribit.TestBaseURL, apiKey, secretKey)
	b.GetOrderBook("BTC-PERPETUAL", 10)
}
