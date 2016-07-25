package favor

import (
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

func TestURLConstruction(t *testing.T) {
	s, _ := New(dummyToken)
	expectedURL := "https://api.askfavor.com/api/v5/farts"
	actualURL := s.BuildURL("farts", map[string]string{})

	if actualURL != expectedURL {
		t.Errorf("Expected URL to be: %v\nActually generated: %v\n", expectedURL, actualURL)
	}
}
