package favor

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateFormString(t *testing.T) {
	x := RequestFavor{
		Title:        "Salty Greg's Frog House",
		Wants:        "One order of frog's legs, please.",
		Lat:          -33.865143,
		Lng:          151.209900,
		Street:       "42 Wallaby Way",
		Zipcode:      "2000",
		MarketID:     0,
		PrimetimeAck: 0,
		Apt:          "123",
		Notes:        "In west Philadelphia, born and raised.",
		MerchantID:   12345,
	}
	expected := "apt=123&lat=-33.865143&lng=151.2099&market_id=0&meal_id=0&merchant_id=12345&notes=In+west+Philadelphia%2C+born+and+raised.&origin_category_id=0&origin_meal_id=0&origin_order_type=0&primetime_ack=0&street=42+Wallaby+Way&title=Salty+Greg%27s+Frog+House&wants=One+order+of+frog%27s+legs%2C+please.&zipcode=2000"
	actual, err := x.CreateFormString()
	if err != nil {
		t.Errorf("Constructor allowed new Favor to be built with faulty token")
	}
	assert.Equal(t, expected, actual)
}
