package favor

import (
	"github.com/stretchr/testify/assert"
	"net/url"
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
	expected := url.Values{"market_id": []string{"0"}, "title": []string{"Salty Greg's Frog House"}, "wants": []string{"One order of frog's legs, please."}, "zipcode": []string{"2000"}, "origin_category_id": []string{"0"}, "origin_order_type": []string{"0"}, "lat": []string{"-33.865143"}, "apt": []string{"123"}, "notes": []string{"In west Philadelphia, born and raised."}, "street": []string{"42 Wallaby Way"}, "merchant_id": []string{"12345"}, "origin_meal_id": []string{"0"}, "lng": []string{"151.2099"}, "primetime_ack": []string{"0"}, "meal_id": []string{"0"}}
	actual := x.CreateFormString()
	assert.Equal(t, expected, actual)
}
