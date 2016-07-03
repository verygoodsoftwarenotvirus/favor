package solid

import (
	"encoding/json"
	"strconv"
)

// MerchantHours describes the hours a merchant is open, or will otherwise
// allow Favor runners to place orders. The response object is structured
// kind of oddly, in that Days as a string returns indices that represent
// days that those hours apply to. So I'm guessing that ideally, for a
// merchant open from 7 AM to 9 PM Monday through Friday, the response would
// look like this:
//      "days": ["0", "1", "2", "3", "4", "5", "6"],
//      "open": [{
//          "start": "0700",
//          "end": "2100"
//      }]
// However, in my limited experience, it ends up being that Days is rarely
// more than one index long.
type MerchantHours struct {
	Days []string `json:"days"`
	Open []struct {
		Start string `json:"start"`
		End   string `json:"end"`
	} `json:"open"`
}

// Merchant describes restauraunts or places of business that cooperate
// with Favor, to my understanding. Favor will process any request you
// make with them, regardless of partnership, but the app uses Merchants
// to provide you more detailed information about what that merchant offers
type Merchant struct {
	ID              string          `json:"id"`
	Name            string          `json:"name"`
	HasExpandedMenu string          `json:"has_expanded_menu"`
	Distance        float64         `json:"distance"`
	Address         string          `json:"address"`
	Zipcode         string          `json:"zipcode"`
	Hours           []MerchantHours `json:"hours"`
}

// Merchants is a container struct set up so that we
// can call sort.Sort() on any given slice of merchants.
type Merchants []Merchant

func (m Merchants) Len() int {
	return len(m)
}

func (m Merchants) Less(i, j int) bool {
	return m[i].Name < m[j].Name
}

func (m Merchants) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

// GetMerchants is used to retrieve Merchants from the Favor API.
func (s Solid) GetMerchants(lat, long float64) ([]Merchant, error) {
	urlParams := map[string]string{
		"lat":             strconv.FormatFloat(lat, 'f', -1, 64),
		"lng":             strconv.FormatFloat(long, 'f', -1, 64),
		"location_source": "gps",
	}
	url := s.BuildURL("merchants", urlParams)
	merchantData, err := s.makeAPIRequest("get", url)
	mr := struct {
		Merchants Merchants `json:"merchants"`
	}{}

	err = json.Unmarshal(merchantData, &mr)
	if err != nil {
		return nil, err
	}
	return mr.Merchants, nil
}
