package bybit_broker

import "testing"

func TestBybitBroker_GetOpenOrders(t *testing.T) {
	b := NewBroker("https://api-testnet.bybit.com/", "6IASD6KDBdunn5qLpT", "nXjZMUiB3aMiPaQ9EUKYFloYNd0zM39RjRWF")
	orders, err := b.GetOpenOrders("BTCUSD")
	if err != nil {
		t.Error(err)
		return
	}
	for _, v := range orders {
		t.Logf("%#v", v)
	}
}
