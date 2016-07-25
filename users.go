package favor

// User represents a user of the service
type User struct {
	ID         string `json:"id"`
	Forename   string `json:"forename"`
	Surname    string `json:"surname"`
	Phone      string `json:"phone"`
	Email      string `json:"email"`
	Countasked string `json:"countasked"`
	FbID       int    `json:"fb_id"`
	Image      string `json:"image"`
}
