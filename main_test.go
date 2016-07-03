package solid

import (
	"os"
	"testing"
)

var dummyToken string

func init() {
	dummyToken = "thisisarandomstringfortestinglol"
}

func TestBadTokenInput(t *testing.T) {
	s, _ := New("fartsandbutts")
	if s != nil {
		t.Errorf("Constructor allowed new Favor to be built with faulty token")
	}
}

func TestGetMerchants(t *testing.T) {
	token := os.Getenv("FAVOR_TOKEN")
	if token == "" {
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

func TestURLConstruction(t *testing.T) {
	s, _ := New(dummyToken)
	expectedURL := "https://api.askfavor.com/api/v5/farts"
	actualURL := s.BuildURL("farts", map[string]string{})

	if actualURL != expectedURL {
		t.Errorf("Expected URL to be: %v\nActually generated: %v\n", expectedURL, actualURL)
	}
}
