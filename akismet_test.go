package akismet

import (
	"net/url"
	"os"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	apiEndpoints["withOutKey"] = apiEndpoint{"without-key", "GET", false}
	apiEndpoints["withKey"] = apiEndpoint{"with-key", "GET", true}

	os.Exit(m.Run())
}

func TestNewClient(t *testing.T) {
	client := NewClient("test_api_key", "test_site")
	assert.NotNil(t, client)
}

func TestGetEndpointURLFail(t *testing.T) {
	client := NewClient("test_api_key", "test_site")
	address, err := client.getEndpointURL("notExistEndpoint")
	assert.Empty(t, address)
	assert.Error(t, err)
}

func TestGetEndpointURL(t *testing.T) {
	client := NewClient("test_api_key", "test_site")
	address, err := client.getEndpointURL("withOutKey")
	assert.Nil(t, err)
	assert.Equal(t, "https://rest.akismet.com/1.1/without-key", address)
}

func TestGetEndpointURLWithKey(t *testing.T) {
	client := NewClient("test_api_key", "test_site")
	address, err := client.getEndpointURL("withKey")
	assert.Nil(t, err)
	assert.Equal(t, "https://test_api_key.rest.akismet.com/1.1/with-key", address)
}

func TestVeryfiClientNotValid(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", "https://rest.akismet.com/1.1/verify-key?blog=test_site&key=test_api_key", httpmock.NewStringResponder(200, "invalid"))

	client := NewClient("test_api_key", "test_site")
	err := client.VeryfiClient()
	assert.Error(t, err)
}

func TestVeryfiClientInternalError(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", "https://rest.akismet.com/1.1/verify-key?blog=test_site&key=test_api_key", httpmock.NewStringResponder(500, ""))

	client := NewClient("test_api_key", "test_site")
	err := client.VeryfiClient()
	assert.Error(t, err)
}

func TestVeryfiClient(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", "https://rest.akismet.com/1.1/verify-key?blog=test_site&key=test_api_key", httpmock.NewStringResponder(200, `valid`))

	client := NewClient("test_api_key", "test_site")
	err := client.VeryfiClient()
	assert.Nil(t, err)
}

func TestParseOptions(t *testing.T) {
	options := Options{
		UserIP:      "",
		UserAgent:   "",
		Referrer:    "",
		Permalink:   "",
		Author:      "",
		AuthorEmail: "",
		AuthorURL:   "",
		Content:     "",
		Created:     "",
		Modified:    "",
		Lang:        "",
		Charset:     "",
		UserRole:    "",
		IsTest:      "",
	}

	_, err := options.parse()
	assert.Error(t, err)

	options = Options{
		UserIP:      "127.0.0.1",
		UserAgent:   "TestUserAgent",
		Referrer:    "",
		Permalink:   "",
		Author:      "",
		AuthorEmail: "",
		AuthorURL:   "",
		Content:     "",
		Created:     "",
		Modified:    "",
		Lang:        "",
		Charset:     "",
		UserRole:    "",
		IsTest:      "",
	}

	r, err := options.parse()
	assert.Nil(t, err)
	expected := url.Values{}
	expected.Add("user_ip", "127.0.0.1")
	expected.Add("user_agent", "TestUserAgent")
	assert.Equal(t, expected, *r)

	options = Options{
		UserIP:      "127.0.0.1",
		UserAgent:   "TestUserAgent",
		Referrer:    "TestReferer",
		Permalink:   "TestPermaLink",
		Author:      "TestAuthor",
		AuthorEmail: "TestAuthorEmail",
		AuthorURL:   "TestAuthorURL",
		Content:     "TestContent",
		Created:     "2012-11-01T22:08:41+00:00",
		Modified:    "2017-11-01T22:08:41+00:00",
		Lang:        "en_EN",
		Charset:     "UTF-8",
		UserRole:    "regular",
		IsTest:      "no",
	}

	r, err = options.parse()
	assert.Nil(t, err)
	expected = url.Values{}
	expected.Add("user_ip", "127.0.0.1")
	expected.Add("user_agent", "TestUserAgent")
	expected.Add("referrer", "TestReferer")
	expected.Add("permalink", "TestPermaLink")
	expected.Add("comment_author", "TestAuthor")
	expected.Add("comment_author_email", "TestAuthorEmail")
	expected.Add("comment_author_url", "TestAuthorURL")
	expected.Add("comment_content", "TestContent")
	expected.Add("comment_date_gmt", "1351807721")
	expected.Add("comment_post_modified_gmt", "1509574121")
	expected.Add("blog_lang", "en_EN")
	expected.Add("blog_charset", "UTF-8")
	expected.Add("user_role", "regular")
	expected.Add("is_test", "no")
	assert.EqualValues(t, expected, *r)

	options = Options{
		UserIP:      "127.0.0.1",
		UserAgent:   "TestUserAgent",
		Referrer:    "",
		Permalink:   "",
		Author:      "",
		AuthorEmail: "",
		AuthorURL:   "",
		Content:     "",
		Created:     "2015-01-01 00:00:00",
		Modified:    "",
		Lang:        "",
		Charset:     "",
		UserRole:    "",
		IsTest:      "",
	}
	r, err = options.parse()
	assert.Error(t, err)
}

func TestIsSpamMissingRequiredOptions(t *testing.T) {
	client := NewClient("test_api_key", "test_site")
	options := &Options{}
	_, err := client.IsSpam(options)
	assert.EqualError(t, err, "filed UserIP can not be empty, it is required")

	options.UserIP = "TestIP"
	_, err = client.IsSpam(options)
	assert.EqualError(t, err, "filed UserAgent can not be empty, it is required")
}

func TestIsSpamInternal(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", "https://test_api_key.rest.akismet.com/1.1/comment-check?blog=test_site&user_agent=TestUserAgent&user_ip=127.0.0.1", httpmock.NewStringResponder(500, ""))

	client := NewClient("test_api_key", "test_site")
	options := &Options{UserIP: "127.0.0.1", UserAgent: "TestUserAgent"}
	_, err := client.IsSpam(options)
	assert.Error(t, err)
}

func TestIsSpamTrue(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", "https://test_api_key.rest.akismet.com/1.1/comment-check?blog=test_site&user_agent=TestUserAgent&user_ip=127.0.0.1", httpmock.NewStringResponder(200, "true"))

	client := NewClient("test_api_key", "test_site")
	options := &Options{UserIP: "127.0.0.1", UserAgent: "TestUserAgent"}
	res, err := client.IsSpam(options)
	assert.Nil(t, err)
	assert.True(t, res)
}

func TestIsSpamInvalid(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", "https://test_api_key.rest.akismet.com/1.1/comment-check?blog=test_site&user_agent=TestUserAgent&user_ip=127.0.0.1", httpmock.NewStringResponder(200, "invalid"))

	client := NewClient("test_api_key", "test_site")
	options := &Options{UserIP: "127.0.0.1", UserAgent: "TestUserAgent"}
	res, err := client.IsSpam(options)
	assert.Error(t, err)
	assert.False(t, res)
}

func TestIsSpamFalse(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", "https://test_api_key.rest.akismet.com/1.1/comment-check?blog=test_site&user_agent=TestUserAgent&user_ip=127.0.0.1", httpmock.NewStringResponder(200, "false"))

	client := NewClient("test_api_key", "test_site")
	options := &Options{UserIP: "127.0.0.1", UserAgent: "TestUserAgent"}
	res, err := client.IsSpam(options)
	assert.Nil(t, err)
	assert.False(t, res)
}
