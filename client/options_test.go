package client_test

import (
	"context"
	"testing"

	"github.com/rudderlabs/rudder-api-go/client"
	"github.com/rudderlabs/rudder-api-go/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClientOptionBaseURL(t *testing.T) {
	c, err := client.New("some-access-token", client.WithBaseURL("https://some-base-url"))
	assert.NoError(t, err)
	assert.Equal(t, "https://some-base-url/path", c.URL("path"))
}

func TestClientOptionBaseURLEmpty(t *testing.T) {
	_, err := client.New("some-access-token", client.WithBaseURL(""))
	assert.Error(t, err, client.ErrInvalidBaseURL)
}

func TestClientOptionHTTPClient(t *testing.T) {
	httpClient := testutils.NewMockHTTPClient(t, testutils.Call{
		Method:         "GET",
		URL:            "https://example.com/path",
		ResponseStatus: 200,
		ResponseBody:   "test",
	})

	c, err := client.New("some-access-token",
		client.WithBaseURL("https://example.com"),
		client.WithHTTPClient(httpClient))
	require.Nil(t, err)

	res, err := c.Do(context.Background(), "GET", "path", nil)
	require.NoError(t, err)
	assert.Equal(t, "test", string(res))
	httpClient.AssertNumberOfCalls()
}

func TestClientOptionHTTPClientNil(t *testing.T) {
	_, err := client.New("some-access-token", client.WithHTTPClient(nil))
	assert.Equal(t, client.ErrInvalidHTTPClient, err)
}
