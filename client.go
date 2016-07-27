package favor

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// Client is our basic struct for making Favor API requests
type Client struct {
	Token  string
	Client http.Client
	Secure bool
}

// New is a constructor function returning a new instance of a Favor Client.
func New(token string) (*Client, error) {
	// As far as I can tell, in my limited research, the Favor token needs to be 32 digits long
	if len(token) < 32 {
		return nil, fmt.Errorf("The token provided is the incorrect length, a normal favorToken is 32 characters long.")
	}
	c := &Client{Token: token, Secure: true}
	return c, nil
}

// BuildURL constructs our Favor API URL when provided with an endpoint to hit and
// any necessary query params. Example usages would be:
// favor.BuildURL("hello", map[string]string{})
//      => "https://api.askfavor.com/api/v5/hello"
//
// favor.BuildURL("hello", map[string]string{"lol": "yup"})
//      => "https://api.askfavor.com/api/v5/hello?lol=yup"
func (c Client) BuildURL(endpoint string, params map[string]string) string {
	var baseURL string
	if c.Secure { // Lord forgive me, for I have sinned.
		baseURL = "https://api.askfavor.com/api/v5/"
	} else {
		baseURL = "http://api.askfavor.com/api/v5/"
	}
	v := url.Values{}
	for p, val := range params {
		v.Add(p, val)
	}
	var paramsString string
	if len(params) > 0 {
		paramsString = fmt.Sprintf("?%v", v.Encode())
	}

	returnURL := fmt.Sprintf("%v%v%v", baseURL, endpoint, paramsString)

	return returnURL
}

// Unexported function used to actually send requests off.
func (c Client) makeAPIRequest(method string, url string) ([]byte, error) {
	req, err := http.NewRequest(strings.ToUpper(method), url, nil)
	if err != nil {
		err = fmt.Errorf("API request failed to build and returned this error:\n %v", err)
		return nil, err
	}

	req.Header.Add("favorToken", c.Token)

	res, err := c.Client.Do(req)
	if err != nil {
		err = fmt.Errorf("API request failed to complete and returned this error:\n %v", err)
		return nil, err
	}
	defer res.Body.Close()

	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		err = fmt.Errorf("API response failed to close and returned this error:\n %v", err)
		return nil, err
	}

	return responseBody, nil
}

// Unexported function used to actually send requests off.
func (c Client) makeAPIRequestWithBody(method string, url string, body url.Values) ([]byte, error) {
	b := strings.NewReader(body.Encode())
	req, err := http.NewRequest(strings.ToUpper(method), url, b)
	if err != nil {
		err = fmt.Errorf("API request failed to build and returned this error:\n %v", err)
		return nil, err
	}

	req.Header.Add("favorToken", c.Token)

	res, err := c.Client.Do(req)
	if err != nil {
		err = fmt.Errorf("API request failed to complete and returned this error:\n %v", err)
		return nil, err
	}
	defer res.Body.Close()

	_ = req.ParseForm()
	// x, _ := ioutil.ReadAll(req.Body)
	fmt.Printf("req.Method   : %v\n", req.Method)
	fmt.Printf("req.URL      : %v\n", req.URL)
	fmt.Printf("req.Header   : %v\n", req.Header)
	fmt.Printf("req.Host     : %v\n", req.Host)
	fmt.Printf("req.Form     : %v\n", req.Form)
	fmt.Printf("req.PostForm : %v\n", req.PostForm)

	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		err = fmt.Errorf("API response failed to close and returned this error:\n %v", err)
		return nil, err
	}
	fmt.Printf("responseBody : %v\n", string(responseBody))

	return responseBody, nil
}
