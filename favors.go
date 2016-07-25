package favor

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
