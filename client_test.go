package favor

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

var dummyToken string

func init() {
	dummyToken = "thisisarandomstringfortestinglol"
}

func setupMockClient(response string) (*httptest.Server, http.Client) {
	/*
		Originally I wanted to just have the API make requests to the HTTPS endpoints
		that Favor sets up. Unfortunately, I couldn't get httptest.NewTLSServer working
		so I had to make security a boolean value, which is sad, and makes me sad.

		If somebody can/wants to issue a PR fixing this, you will have my eternal gratitude.

		This code is lovingly borrowed from http://keighl.com/post/mocking-http-responses-in-golang/
	*/

	// Test server that always responds with 200 code, and specific payload
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, response)
	}))

	// Make a transport that reroutes all traffic to the example server
	transport := &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			return url.Parse(server.URL)
		},
	}
	// Make a http.Client with the transport
	httpClient := http.Client{Transport: transport}
	return server, httpClient
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
