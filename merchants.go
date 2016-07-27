package favor

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// HoursOpen is only a struct because if it were anonymously in MerchantHoursResponse,
// testing would be much more verbose, and I don't like it.
type HoursOpen struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

// MerchantHoursResponse describes the hours a merchant is open, or will otherwise
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
type MerchantHoursResponse struct {
	Days []string    `json:"days"`
	Open []HoursOpen `json:"open"`
}

// MerchantHours is a more usable struct than MerchantHours, but should
// reflect the same information.
type MerchantHours struct {
	Open  time.Time
	Close time.Time
}

func (m MerchantHoursResponse) buildTimesFromHours() (map[int][]MerchantHours, error) {
	hours := map[int][]MerchantHours{}
	now := time.Now()

	for _, d := range m.Days {
		newDay := MerchantHours{}
		day, err := strconv.Atoi(d)
		if err != nil {
			return nil, fmt.Errorf("Error parsing day %v into integer!:\n%v", d, err)
		}
		for _, o := range m.Open {
			closeString := o.End
			closeDay := day
			if strings.HasPrefix(o.End, "+") {
				if closeDay == 6 {
					closeDay = 0
				} else {
					closeDay++
				}
				closeString = strings.TrimPrefix(closeString, "+")
			}

			openHour, hourErr := strconv.Atoi(o.Start[:2])
			openMinute, minuteErr := strconv.Atoi(o.Start[2:])
			if hourErr != nil || minuteErr != nil {
				return nil, fmt.Errorf("Error parsing open %v into integer!:\n%v", o.Start, err)
			}
			closeHour, hourErr := strconv.Atoi(closeString[:2])
			closeMinute, minuteErr := strconv.Atoi(closeString[2:])
			if hourErr != nil || minuteErr != nil {
				return nil, fmt.Errorf("Error parsing close %v into integer!:\n%v", closeString, err)
			}

			newDay.Open = time.Date(now.Year(), now.Month(), day, openHour, openMinute, 0, 0, time.UTC)
			newDay.Close = time.Date(now.Year(), now.Month(), closeDay, closeHour, closeMinute, 0, 0, time.UTC)
			hours[day] = append(hours[day], newDay)
		}
	}
	return hours, nil
}

// IsOpenAt is a helper function to determine if a merchant
// is available for placing orders at a given time. Possible
// usage would be something like m.IsOpenAt(time.Now())
// func (m Merchant) IsOpenAt(t time.Time) bool {
// 	return false
// }

// Merchant describes restauraunts or places of business that cooperate
// with Favor, to my understanding. Favor will process any request you
// make with them, regardless of partnership, but the app uses Merchants
// to provide you more detailed information about what that merchant offers
type Merchant struct {
	ID              string          `json:"id"`
	FranchiseID     string          `json:"franchise_id,omitempty"`
	MarketID        string          `json:"market_id,omitempty"`
	Name            string          `json:"name"`
	Phone           string          `json:"phone,omitempty"`
	Address         string          `json:"address"`
	City            string          `json:"city,omitempty"`
	State           string          `json:"state,omitempty"`
	Zipcode         string          `json:"zipcode"`
	Distance        float64         `json:"distance,omitempty"`
	HasExpandedMenu string          `json:"has_expanded_menu,omitempty"`
	Hours           []MerchantHours `json:"hours,omitempty"`
	Lat             string          `json:"lat,omitempty"`
	Lng             string          `json:"lng,omitempty"`
	IsCarOnly       string          `json:"is_car_only,omitempty"`
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

// GetMerchant is used to retrieve a single merchant from the Favor API.
func (c Client) GetMerchant(id string) (Merchant, error) {
	urlParams := map[string]string{}
	uri := c.BuildURL(fmt.Sprintf("merchant/%v", id), urlParams)
	merchantData, err := c.makeAPIRequest("get", uri)
	if err != nil {
		return Merchant{}, err
	}
	mr := struct {
		Merchant Merchant `json:"merchant"`
	}{}

	err = json.Unmarshal(merchantData, &mr)
	if err != nil {
		return Merchant{}, err
	}
	return mr.Merchant, nil
}

// GetMerchants is used to retrieve Merchants from the Favor API.
func (c Client) GetMerchants(lat, long float64) ([]Merchant, error) {
	urlParams := map[string]string{
		"lat":             strconv.FormatFloat(lat, 'f', -1, 64),
		"lng":             strconv.FormatFloat(long, 'f', -1, 64),
		"location_source": "gps",
	}
	uri := c.BuildURL("merchants", urlParams)
	merchantData, err := c.makeAPIRequest("get", uri)
	mr := struct {
		Merchants Merchants `json:"merchants"`
	}{}

	err = json.Unmarshal(merchantData, &mr)
	if err != nil {
		return nil, err
	}
	return mr.Merchants, nil
}
