package favor

import (
	"reflect"
	"testing"
	"time"
)

func TestGetMerchant(t *testing.T) {
	s, err := New(dummyToken)
	if err != nil {
		t.Errorf("Constructor failed with the following error: %v", err)
		t.FailNow()
	}
	s.Secure = false

	dummyMerchantResponse := `
	{
		"merchant": {
			"id": "1234",
			"franchise_id": "1",
			"market_id": "1",
			"name": "Farts McGregor's Corntopia",
			"phone": "5124206969",
			"address": "42 Wallaby Way",
			"city": "Sydney",
			"state": "NSW",
			"zipcode": "2000",
			"has_expanded_menu": "1",
			"lat": "-33.865143",
			"lng": "151.209900",
			"is_car_only": "0"
		}
	}`

	server, client := setupMockClient(dummyMerchantResponse)
	defer server.Close()

	s.Client = client

	expectedMerchant := Merchant{
		ID:              "1234",
		FranchiseID:     "1",
		MarketID:        "1",
		Name:            "Farts McGregor's Corntopia",
		Phone:           "5124206969",
		Address:         "42 Wallaby Way",
		City:            "Sydney",
		State:           "NSW",
		Zipcode:         "2000",
		HasExpandedMenu: "1",
		Lat:             "-33.865143",
		Lng:             "151.209900",
		IsCarOnly:       "0",
	}

	actualMerchant, err := s.GetMerchant("1234")
	if err != nil {
		t.Errorf("GetMerchant failed with the following error: %v", err)
	}

	if !reflect.DeepEqual(expectedMerchant, actualMerchant) {
		t.Errorf("Retrieved Merchant differs from expected result.\n")
		t.Errorf("Expected:\n%v+ \n", expectedMerchant)
		t.Errorf("Received:\n%v+ \n", actualMerchant)
	}
}

func TestGetMerchants(t *testing.T) {
	s, err := New(dummyToken)
	if err != nil {
		t.Errorf("Constructor failed with the following error: %v", err)
		t.FailNow()
	}
	s.Secure = false

	dummyMerchantResponse := `
	{
		"merchants": [{
			"id": "1234",
			"franchise_id": "1",
			"market_id": "1",
			"name": "Farts McGregor's Corntopia",
			"phone": "5124206969",
			"address": "42 Wallaby Way",
			"city": "Sydney",
			"state": "NSW",
			"zipcode": "2000",
			"has_expanded_menu": "1",
			"lat": "-33.865143",
			"lng": "151.209900",
			"is_car_only": "0"
		}]
	}`

	server, client := setupMockClient(dummyMerchantResponse)
	defer server.Close()

	s.Client = client

	expectedMerchants := []Merchant{
		Merchant{
			ID:              "1234",
			FranchiseID:     "1",
			MarketID:        "1",
			Name:            "Farts McGregor's Corntopia",
			Phone:           "5124206969",
			Address:         "42 Wallaby Way",
			City:            "Sydney",
			State:           "NSW",
			Zipcode:         "2000",
			HasExpandedMenu: "1",
			Lat:             "-33.865143",
			Lng:             "151.209900",
			IsCarOnly:       "0",
		},
	}

	actualMerchants, err := s.GetMerchants(30.234855, -97.7322537)
	if err != nil {
		t.Errorf("GetMerchants failed with the following error: %v", err)
	}

	if len(actualMerchants) != 1 {
		t.Errorf("More merchants than expected were parsed.")
	}

	if !reflect.DeepEqual(expectedMerchants, actualMerchants) {
		t.Errorf("Retrieved Merchant differs from expected result.\n")
		t.Errorf("Expected:\n%v+ \n", expectedMerchants)
		t.Errorf("Received:\n%v+ \n", actualMerchants)
	}
}

func merchantHoursAreEqual(a, b map[int][]MerchantHours, t *testing.T) bool {
	if len(a) != len(b) {
		return false
	}
	for day := range a {
		for index, hours := range a[day] {
			if !hours.Open.Equal(b[day][index].Open) || !hours.Close.Equal(b[day][index].Close) {
				t.Errorf("\nExpected a[%v][%v].Open to equal \n%v.\nInstead equaled \n%v\n", day, index, b[day][index].Open, a[day][index].Open)
				t.Errorf("\nExpected a[%v][%v].Close to equal \n%v.\nInstead equaled \n%v\n", day, index, b[day][index].Close, a[day][index].Close)
				return false
			}
		}
	}
	return true
}

func TestBuildTimesFromHours(t *testing.T) {
	now := time.Now()

	simple := MerchantHoursResponse{
		Days: []string{"0"},
		Open: []HoursOpen{
			HoursOpen{
				Start: "0700",
				End:   "1500",
			},
		},
	}
	expectedSimple := map[int][]MerchantHours{
		0: []MerchantHours{
			MerchantHours{
				Open:  time.Date(now.Year(), now.Month(), 0, 7, 0, 0, 0, time.UTC),
				Close: time.Date(now.Year(), now.Month(), 0, 15, 0, 0, 0, time.UTC),
			},
		},
	}
	actualSimple, err := simple.buildTimesFromHours()
	if err != nil {
		t.Errorf("Simple parsing failed with error:\n%v", err)
	}
	if !merchantHoursAreEqual(actualSimple, expectedSimple, t) {
		t.Errorf("Simple parsing produced unexpected results")
	}

	slightlyMoreComplicated := MerchantHoursResponse{
		Days: []string{"0"},
		Open: []HoursOpen{
			HoursOpen{
				Start: "0700",
				End:   "+0300",
			},
		},
	}
	expectedComplicated := map[int][]MerchantHours{
		0: []MerchantHours{
			MerchantHours{
				Open:  time.Date(now.Year(), now.Month(), 0, 7, 0, 0, 0, time.UTC),
				Close: time.Date(now.Year(), now.Month(), 1, 3, 0, 0, 0, time.UTC),
			},
		},
	}
	actualComplicated, err := slightlyMoreComplicated.buildTimesFromHours()

	if err != nil {
		t.Errorf("More complicated parsing failed with error:\n%v", err)
	}
	if !merchantHoursAreEqual(actualComplicated, expectedComplicated, t) {
		t.Errorf("More complicated parsing produced unexpected results")
	}

	shouldThrowAnError := MerchantHoursResponse{
		Days: []string{"0"},
		Open: []HoursOpen{
			HoursOpen{
				Start: "NEVER",
				End:   "NOPE",
			},
		},
	}
	z, err := shouldThrowAnError.buildTimesFromHours()
	if err == nil {
		t.Errorf("Somehow the dumb input didn't fail, resulting object:\n%v", z)
	}

	twoOpenings := MerchantHoursResponse{
		Days: []string{"0"},
		Open: []HoursOpen{
			HoursOpen{
				Start: "0700",
				End:   "1500",
			}, {
				Start: "1800",
				End:   "2100",
			},
		},
	}
	expectedTwoOpenings := map[int][]MerchantHours{
		0: []MerchantHours{
			MerchantHours{
				Open:  time.Date(now.Year(), now.Month(), 0, 7, 0, 0, 0, time.UTC),
				Close: time.Date(now.Year(), now.Month(), 0, 15, 0, 0, 0, time.UTC),
			},
			MerchantHours{
				Open:  time.Date(now.Year(), now.Month(), 0, 18, 0, 0, 0, time.UTC),
				Close: time.Date(now.Year(), now.Month(), 0, 21, 0, 0, 0, time.UTC),
			},
		},
	}
	actualTwoOpenings, err := twoOpenings.buildTimesFromHours()
	if err != nil {
		t.Errorf("Two openings parsing failed with error:\n%v", err)
	}
	if !merchantHoursAreEqual(expectedTwoOpenings, actualTwoOpenings, t) {
		t.Errorf("Two openings parsing produced unexpected results")
	}
}
