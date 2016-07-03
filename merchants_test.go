package solid

import (
	"os"
	"testing"
	"time"
)

func TestGetMerchants(t *testing.T) {
	token := os.Getenv("FAVOR_TOKEN")
	if true { // token == "" {
		t.Skip("skipping test; $FAVOR_TOKEN not set")
	}
	s, err := New(token)

	if err != nil {
		t.Errorf("Constructor failed with the following error: %v", err)
		t.FailNow()
	}

	_, err = s.GetMerchants(30.234855, -97.7322537)
	if err != nil {
		t.Errorf("GetMerchants failed with the following error: %v", err)
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
