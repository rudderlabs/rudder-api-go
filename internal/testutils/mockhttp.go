package testutils

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Call struct {
	Method         string
	URL            string
	RequestBody    string
	ResponseStatus int
	ResponseBody   string
	ResponseError  error
}

type mockHTTPClient struct {
	t         *testing.T
	callIndex int
	calls     []Call
}

func NewMockHTTPClient(t *testing.T, calls ...Call) *mockHTTPClient {
	c := &mockHTTPClient{
		t:     t,
		calls: calls,
	}

	return c
}

func (c *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	if c.callIndex >= len(c.calls) {
		err := fmt.Errorf("received unexpected request: %v", req)
		c.t.Error(err)
		return nil, err
	}

	call := c.calls[c.callIndex]
	c.callIndex += 1

	assert.Equal(c.t, call.Method, req.Method)
	if call.RequestBody != "" {
		body, err := ioutil.ReadAll(req.Body)
		require.NoError(c.t, err)
		assert.JSONEq(c.t, call.RequestBody, string(body))
	}

	return &http.Response{
		StatusCode: call.ResponseStatus,
		Body:       io.NopCloser(strings.NewReader(call.ResponseBody)),
	}, call.ResponseError
}

func (c *mockHTTPClient) AssertNumberOfCalls() {
	if c.callIndex < len(c.calls) {
		c.t.Errorf("missing calls: expected %d, received %d", len(c.calls), c.callIndex)
	}
}
