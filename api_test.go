package bossapi

import (
	"encoding/json"
	"testing"
	"time"
)

func TestCurrencyOf(t *testing.T) {
	for in, want := range map[string]string{
		"BTC-27JUN25-100000-C": "BTC", "ETH-1JAN26-3000-P": "ETH", "BTC_USDC-27JUN25-100000-C": "BTC",
	} {
		if got := CurrencyOf(in); got != want {
			t.Errorf("CurrencyOf(%q) = %q, want %q", in, got, want)
		}
	}
}

// The wire names are the contract: a renamed JSON tag breaks every
// deployed client, so the important ones are pinned here.
func TestWireNames(t *testing.T) {
	b, _ := json.Marshal(Quote{Bid: 1, Ask: 2, Mid: 1.5, Underlying: 100, FeeRate: 0.0003, FeeCap: 0.125, UpdatedAt: time.Unix(0, 0).UTC()})
	var m map[string]any
	json.Unmarshal(b, &m)
	for _, k := range []string{"bid_usd", "ask_usd", "mid_usd", "underlying_usd", "fee_rate", "fee_cap", "updated_at", "instrument", "exchange"} {
		if _, ok := m[k]; !ok {
			t.Errorf("Quote lost wire field %q", k)
		}
	}
	b, _ = json.Marshal(Error{Code: CodeNotFound, Message: "x"})
	if string(b) != `{"code":"not_found","message":"x"}` {
		t.Errorf("Error wire form = %s", b)
	}
	if Prefix != "/api/v1/" {
		t.Errorf("Prefix = %q", Prefix)
	}
}
