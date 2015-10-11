// Package akismet is Go utils for working with Akismet spam detection service
package akismet

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// Basic informations about Akismet API
const (
	APIAddress  = "rest.akismet.com"
	APIProtocol = "https"
	APIVersion  = "1.1"
	DateFormat  = time.RFC3339
)

// Client is Akismet client struct
type Client struct {
	apiKey     string
	site       string
	httpClient *http.Client
}

// Options is a struct which contains all of possible arguments for Akismet
type Options struct {
	UserIP      string
	UserAgent   string
	Referrer    string
	Permalink   string
	Author      string
	AuthorEmail string
	AuthorURL   string
	Content     string
	Created     string
	Modified    string
	Lang        string
	Charset     string
	UserRole    string
	IsTest      string
}

type apiEndpoint struct {
	path           string
	method         string
	apiKeyRequired bool
}

var apiEndpoints = map[string]apiEndpoint{
	"verifyKey":    apiEndpoint{"verify-key", "GET", false},
	"commentCheck": apiEndpoint{"comment-check", "POST", true},
	"submitSpam":   apiEndpoint{"submit-spam", "POST", true},
	"submitHam":    apiEndpoint{"submit-ham", "POST", true},
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

// IsSpam is a method which check if passed Options struct is spam or not
func (c *Client) IsSpam(o *Options) (bool, error) {
	v, err := o.parse()
	if err != nil {
		return false, err
	}

	v.Add("blog", c.site)
	endpointURL, err := c.getEndpointURL("commentCheck")
	if err != nil {
		return false, err
	}

	address, err := url.Parse(endpointURL)
	if err != nil {
		return false, err
	}

	address.RawQuery = v.Encode()
	res, err := c.httpClient.Get(address.String())
	if err != nil {
		return false, err
	}

	if res.StatusCode != 200 {
		return false, errors.New("something went wrong, HTTP status code is not equals 200")
	}

	r, _ := getResponseBodyAsString(res)
	switch r {
	case "true":
		return true, nil
	case "invalid":
		return false, errors.New(res.Header.Get("X-Akismet-Debug-Help"))
	}

	return false, nil
}

// SubmitSpam is method which send to Akismet API request about found spam
func (c *Client) SubmitSpam(o *Options) error {
	v, err := o.parse()
	if err != nil {
		return err
	}

	v.Add("blog", c.site)
	endpointURL, err := c.getEndpointURL("submitSpam")
	if err != nil {
		return err
	}

	address, err := url.Parse(endpointURL)
	if err != nil {
		return err
	}

	address.RawQuery = v.Encode()
	res, err := c.httpClient.Get(address.String())
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return errors.New("something went wrong, HTTP status code is not equals 200")
	}

	r, _ := getResponseBodyAsString(res)
	if r != "Thanks for making the web a better place." {
		return errors.New("something went wrong")
	}

	return nil
}

// SubmitHam is method which send to Akismet API request about found ham
func (c *Client) SubmitHam(o *Options) error {
	v, err := o.parse()
	if err != nil {
		return err
	}

	v.Add("blog", c.site)
	endpointURL, err := c.getEndpointURL("submitHam")
	if err != nil {
		return err
	}

	address, err := url.Parse(endpointURL)
	if err != nil {
		return err
	}

	address.RawQuery = v.Encode()
	res, err := c.httpClient.Get(address.String())
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return errors.New("something went wrong, HTTP status code is not equals 200")
	}

	r, _ := getResponseBodyAsString(res)
	if r != "Thanks for making the web a better place." {
		return errors.New("something went wrong")
	}

	return nil
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

func (o *Options) parse() (*url.Values, error) {
	if o.UserIP == "" {
		return nil, errors.New("filed UserIP can not be empty, it is required")
	}

	if o.UserAgent == "" {
		return nil, errors.New("filed UserAgent can not be empty, it is required")
	}

	v := url.Values{}
	v.Add("user_ip", o.UserIP)
	v.Add("user_agent", o.UserAgent)

	if o.Referrer != "" {
		v.Add("referrer", o.Referrer)
	}

	if o.Permalink != "" {
		v.Add("permalink", o.Permalink)
	}

	if o.Author != "" {
		v.Add("comment_author", o.Author)
	}

	if o.AuthorEmail != "" {
		v.Add("comment_author_email", o.AuthorEmail)
	}

	if o.AuthorURL != "" {
		v.Add("comment_author_url", o.AuthorURL)
	}

	if o.Content != "" {
		v.Add("comment_content", o.Content)
	}

	if o.Created != "" {
		created, err := time.Parse(DateFormat, o.Created)
		if err != nil {
			return nil, err
		}
		v.Add("comment_date_gmt", fmt.Sprint(created.Unix()))
	}

	if o.Modified != "" {
		modified, err := time.Parse(DateFormat, o.Modified)
		if err != nil {
			return nil, err
		}
		v.Add("comment_post_modified_gmt", fmt.Sprint(modified.Unix()))
	}

	if o.Lang != "" {
		v.Add("blog_lang", o.Lang)
	}

	if o.Charset != "" {
		v.Add("blog_charset", o.Charset)
	}

	if o.UserRole != "" {
		v.Add("user_role", o.UserRole)
	}

	if o.IsTest != "" {
		v.Add("is_test", o.IsTest)
	}

	return &v, nil
}
