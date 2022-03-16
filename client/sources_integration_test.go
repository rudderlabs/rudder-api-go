//go:build integrationtest
// +build integrationtest

package client_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/rudderlabs/rudder-api-go/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newClient() (*client.Client, error) {
	return client.New(os.Getenv("RUDDERSTACK_API_ACCESS_TOKEN"), client.WithBaseURL("http://localhost:5555/v2"))
}

func TestClientSourcesIntegration(t *testing.T) {
	c, err := newClient()
	require.NoError(t, err)

	page, err := c.Sources.List(context.Background())
	require.NoError(t, err)

	assert.Len(t, page.Sources, 1)
	j, _ := json.Marshal(page.Sources[0])
	fmt.Printf("%v", string(j))
	assert.Equal(t, 0, page.Paging.Total)
	assert.Empty(t, page.Paging.Next)
}
