package favor

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"strings"
)

// ServerFavorResponse is a simple container struct for the server's
// response to a favor request
type ServerFavorResponse struct {
	Count  int     `json:"count"`
	Favors []Favor `json:"favors"`
}

// Rating represents the rating requested from the customer
type Rating struct {
	RatingFood   string `json:"rating_food"`
	RatingDriver string `json:"rating_driver"`
	Comment      string `json:"comment"`
	UpdatedAt    string `json:"updated_at"`
}

// Receipt represents the receit for a Favor order
type Receipt struct {
	Paid           string `json:"paid"`
	Price          string `json:"price"`
	Tip            string `json:"tip"`
	SuggestedTip   string `json:"suggested_tip"`
	MinimumTip     int    `json:"minimum_tip"`
	DeliveryCharge string `json:"delivery_charge"`
	CcFeeAmount    string `json:"cc_fee_amount"`
	RebatePrice    string `json:"rebate_price"`
}

// RequestFavor represents what we need to send to the Favor server to place a Favor
type RequestFavor struct {
	Title            string  `json:"title"`
	Wants            string  `json:"wants"`
	Lat              float64 `json:"lat"`
	Lng              float64 `json:"lng"`
	Street           string  `json:"street"`
	Zipcode          string  `json:"zipcode"`
	MarketID         int     `json:"market_id"`
	PrimetimeAck     int     `json:"primetime_ack"`
	Apt              string  `json:"apt,omitempty"`
	Notes            string  `json:"notes"`
	MerchantID       int     `json:"merchant_id"`
	MealID           int     `json:"meal_id,omitempty"`
	OriginMealID     int     `json:"origin_meal_id,omitempty"`
	OriginCategoryID int     `json:"origin_category_id,omitempty"`
	OriginOrderType  int     `json:"origin_order_type,omitempty"`
}

// Favor represents our requested task to the service
type Favor struct {
	ID              string   `json:"id"`
	Title           string   `json:"title"`
	Items           []string `json:"items"`
	MerchantID      string   `json:"merchant_id"`
	Stage           string   `json:"stage"`
	LastStatus      string   `json:"last_status"`
	Ratings         Rating   `json:"ratings"`
	CreatedAt       int      `json:"created_at"`
	Customer        User     `json:"customer"`
	DeliveryAddress Address  `json:"delivery_address"`
	Runner          User     `json:"runner"`
	Merchant        Merchant `json:"merchant,omitempty"`
	Receipt         Receipt  `json:"receipt"`
	// Unsure of what AssignStatus is, as none of my intercepted
	// requests seemed to have it populated.
	// AssignStatus    interface{}  `json:"assign_status"`
}

// GetFavor is used to retrieve a single favor from the Favor API.
func (c Client) GetFavor(id string) (Favor, error) {
	urlParams := map[string]string{}
	uri := c.BuildURL(fmt.Sprintf("favors/%v", id), urlParams)
	favorData, err := c.makeAPIRequest("get", uri)
	if err != nil {
		return Favor{}, err
	}
	f := struct {
		Favor Favor `json:"favor"`
	}{}

	err = json.Unmarshal(favorData, &f)
	if err != nil {
		return Favor{}, err
	}
	return f.Favor, nil
}

// GetFavors is used to retrieve a list of favors from the Favor API.
func (c Client) GetFavors() ([]Favor, error) {
	// knownParams := []string{"count", "include_cancelled", "location_source"}
	uri := c.BuildURL("favors/", map[string]string{})
	favorData, err := c.makeAPIRequest("get", uri)
	if err != nil {
		return []Favor{}, err
	}
	f := struct {
		Count  int     `json:"count"`
		Favors []Favor `json:"favors"`
	}{}

	err = json.Unmarshal(favorData, &f)
	if err != nil {
		return []Favor{}, err
	}
	return f.Favors, nil
}

// CreateFormString turns a RequestFavor struct into an appropriate POST form payload
func (rf RequestFavor) CreateFormString() url.Values {
	t := reflect.TypeOf(rf)
	v := reflect.ValueOf(rf)
	u := url.Values{}

	for i := 0; i < v.NumField(); i++ {
		field := strings.Split(t.Field(i).Tag.Get("json"), ",")[0]
		value := fmt.Sprint(v.Field(i).Interface())
		u.Add(field, value)
	}

	return u
}

// PlaceFavor places a Favor order with the Favor API
func (c Client) PlaceFavor(rf RequestFavor) (Favor, error) {
	requestBody := rf.CreateFormString()
	uri := c.BuildURL("favors/", map[string]string{})
	responseData, err := c.makeAPIRequestWithBody("post", uri, requestBody)
	if err != nil {
		return Favor{}, err
	}
	f := struct {
		Favor Favor `json:"favor"`
	}{}

	err = json.Unmarshal(responseData, &f)
	if err != nil {
		return Favor{}, err
	}
	return f.Favor, nil
}
