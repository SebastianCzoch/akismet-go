package akismet

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// Basic informations about Akismet API
const (
	APIAddress  = "rest.akismet.com"
	APIProtocol = "https"
	APIVersion  = "1.1"
)

// Client is Akismet client struct
type Client struct {
	apiKey     string
	site       string
	httpClient *http.Client
}

type apiEndpoint struct {
	path           string
	method         string
	apiKeyRequired bool
}

var apiEndpoints = map[string]apiEndpoint{
	"verifyKey": apiEndpoint{"verify-key", "GET", false},
}

// NewClient is function which create new Akismet client
func NewClient(apiKey, site string) *Client {
	return &Client{
		apiKey:     apiKey,
		site:       site,
		httpClient: &http.Client{},
	}
}

// VeryfiClient is method which check key & site parameters are valid
func (c *Client) VeryfiClient() error {
	endpointURL, err := c.getEndpointURL("verifyKey")
	if err != nil {
		return err
	}

	address, err := url.Parse(endpointURL)
	if err != nil {
		return err
	}
	v := address.Query()
	v.Add("key", c.apiKey)
	v.Add("blog", c.site)
	address.RawQuery = v.Encode()

	res, err := c.httpClient.Get(address.String())
	if err != nil {
		return err
	}

	r, _ := getResponseBodyAsString(res)
	if r == "valid" {
		return nil
	}

	return errors.New("invalid key or blog")
}

func (c *Client) getEndpointURL(name string) (string, error) {
	endpoint, err := getEndpoint(name)
	if err != nil {
		return "", err
	}

	address := url.URL{
		Scheme: APIProtocol,
		Path:   fmt.Sprintf("%s/%s", APIVersion, endpoint.path),
		Host:   APIAddress,
	}

	if endpoint.apiKeyRequired {
		address.Host = fmt.Sprintf("%s.%s", c.apiKey, APIAddress)
	}

	return address.String(), nil

}

func getEndpoint(name string) (*apiEndpoint, error) {
	endpoint, ok := apiEndpoints[name]
	if !ok {
		return nil, fmt.Errorf("endpoint %s not found", name)
	}

	return &endpoint, nil
}

func getResponseBodyAsString(response *http.Response) (string, error) {
	res, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return string(res), nil
}
