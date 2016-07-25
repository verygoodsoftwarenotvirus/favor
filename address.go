package favor

import (
// "fmt"
)

type Address struct {
	ID         string `json:"id"`
	CustomerID string `json:"customer_id"`
	Lat        string `json:"lat"`
	Lng        string `json:"lng"`
	Street     string `json:"street"`
	Zipcode    string `json:"zipcode"`
	Apartment  string `json:"apartment"`
	Notes      string `json:"notes"`
}
