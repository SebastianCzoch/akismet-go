package akismet

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	apiEndpoints = map[string]apiEndpoint{
		"withOutKey": apiEndpoint{"without-key", "GET", false},
		"withKey":    apiEndpoint{"with-key", "GET", true},
	}
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
