package client_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/rudderlabs/rudder-api-go/client"
	"github.com/rudderlabs/rudder-api-go/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClientRETLSourcesList(t *testing.T) {
	ctx := context.Background()

	calls := []testutils.Call{
		{
			Validate: func(req *http.Request) bool {
				return testutils.ValidateRequest(t, req, "GET", "https://api.rudderstack.com/v2/retl-sources", "")
			},
			ResponseStatus: 200,
			ResponseBody: `{
				"sources": [{
					"id": "id-1",
					"type": "type-1",
					"name": "name-1",
					"config": {"key":"val-1"}
				},  {
					"id": "id-2",
					"type": "type-2",
					"name": "name-2",
					"config": {"key":"val-2"}
				}],
				"paging": {
					"total": 3,
					"next": "/retl-sources?page=2"
				}
			}`,
		},
		{
			Validate: func(req *http.Request) bool {
				return testutils.ValidateRequest(t, req, "GET", "https://api.rudderstack.com/v2/retl-sources?page=2", "")
			},
			ResponseStatus: 200,
			ResponseBody: `{
				"sources": [{
					"id": "id-3",
					"type": "type-3",
					"name": "name-3",
					"config": {"key":"val-3"}
				}],
				"paging": {
					"total": 3
				}
			}`,
		},
	}

	httpClient := testutils.NewMockHTTPClient(t, calls...)

	c, err := client.New("some-access-token", client.WithHTTPClient(httpClient))
	require.NoError(t, err)

	page, err := c.RETLSources.List(ctx)
	require.NoError(t, err)
	assert.NotNil(t, page)
	assert.Len(t, page.RETLSources, 2)
	assert.Equal(t, client.RETLSource{ID: "id-1", Type: "type-1", Name: "name-1", Config: []byte(`{"key":"val-1"}`)}, page.RETLSources[0])
	assert.Equal(t, client.RETLSource{ID: "id-2", Type: "type-2", Name: "name-2", Config: []byte(`{"key":"val-2"}`)}, page.RETLSources[1])
	assert.Equal(t, 3, page.Paging.Total)
	assert.Equal(t, "/retl-sources?page=2", page.Paging.Next)

	page, err = c.RETLSources.Next(ctx, page.Paging)
	require.NoError(t, err)
	assert.NotNil(t, page)
	assert.Len(t, page.RETLSources, 1)
	assert.Equal(t, client.RETLSource{ID: "id-3", Type: "type-3", Name: "name-3", Config: []byte(`{"key":"val-3"}`)}, page.RETLSources[0])
	assert.Equal(t, 3, page.Paging.Total)
	assert.Equal(t, "", page.Paging.Next)

	page, err = c.RETLSources.Next(ctx, page.Paging)
	require.NoError(t, err)
	assert.Nil(t, page)

	httpClient.AssertNumberOfCalls()
}

func TestClientRETLSourcesGet(t *testing.T) {
	ctx := context.Background()

	calls := []testutils.Call{
		{
			Validate: func(req *http.Request) bool {
				return testutils.ValidateRequest(t, req, "GET", "https://api.rudderstack.com/v2/retl-sources/some-id", "")
			},
			ResponseStatus: 200,
			ResponseBody: `{
				"source": {
					"id": "some-id",
					"name": "some-name",
					"type": "some-type",
					"config": {
						"key1": "val1"
					},
					"createdAt": "2020-01-01T01:01:01Z",
					"updatedAt": "2020-01-02T01:01:01Z"	
				}
			}`,
		},
	}

	httpClient := testutils.NewMockHTTPClient(t, calls...)

	c, err := client.New("some-access-token", client.WithHTTPClient(httpClient))
	require.NoError(t, err)

	source, err := c.RETLSources.Get(ctx, "some-id")
	require.NoError(t, err)
	assert.NotNil(t, source)
	assert.Equal(t, "some-id", source.ID)
	assert.Equal(t, "some-name", source.Name)
	assert.Equal(t, "some-type", source.Type)
	assert.Equal(t, time.Date(2020, 1, 1, 1, 1, 1, 0, time.UTC), *source.CreatedAt)
	assert.Equal(t, time.Date(2020, 1, 2, 1, 1, 1, 0, time.UTC), *source.UpdatedAt)

	httpClient.AssertNumberOfCalls()
}

func TestClientRETLSourcesCreate(t *testing.T) {
	ctx := context.Background()

	calls := []testutils.Call{
		{
			Validate: func(req *http.Request) bool {
				return testutils.ValidateRequest(t, req, "POST", "https://api.rudderstack.com/v2/retl-sources", `{
					"name": "some-name",
					"type": "some-type",
					"enabled": true,
					"config": { "key1": "val1" }
				}`)
			},
			ResponseStatus: 200,
			ResponseBody: `{
				"source": {
					"id": "some-id",
					"name": "some-name",
					"type": "some-type",
					"config": {
						"key1": "val1"
					},
					"createdAt": "2020-01-01T01:01:01Z",
					"updatedAt": "2020-01-02T01:01:01Z"
				}
			}`,
		},
	}

	httpClient := testutils.NewMockHTTPClient(t, calls...)

	c, err := client.New("some-access-token", client.WithHTTPClient(httpClient))
	require.NoError(t, err)

	source, err := c.RETLSources.Create(ctx, &client.RETLSource{
		Name:      "some-name",
		Type:      "some-type",
		IsEnabled: true,
		Config: json.RawMessage([]byte(`{
			"key1": "val1"
		}`)),
	})
	require.NoError(t, err)
	assert.NotNil(t, source)
	assert.Equal(t, "some-id", source.ID)
	assert.Equal(t, "some-name", source.Name)
	assert.Equal(t, "some-type", source.Type)
	assert.Equal(t, time.Date(2020, 1, 1, 1, 1, 1, 0, time.UTC), *source.CreatedAt)
	assert.Equal(t, time.Date(2020, 1, 2, 1, 1, 1, 0, time.UTC), *source.UpdatedAt)

	httpClient.AssertNumberOfCalls()
}

func TestClientRETLSourcesUpdate(t *testing.T) {
	ctx := context.Background()

	calls := []testutils.Call{
		{
			Validate: func(req *http.Request) bool {
				return testutils.ValidateRequest(t, req, "PUT", "https://api.rudderstack.com/v2/retl-sources/some-id", `{
					"name": "some-name",
					"type": "some-type",
					"enabled": true,
					"config": { "key1": "val1" }
				}`)
			},
			ResponseStatus: 200,
			ResponseBody: `{
				"source": {
					"id": "some-id",
					"name": "some-name",
					"type": "some-type",
					"config": {
						"key1": "val1"
					},
					"createdAt": "2020-01-01T01:01:01Z",
					"updatedAt": "2020-01-02T01:01:01Z"
				}
			}`,
		},
	}

	httpClient := testutils.NewMockHTTPClient(t, calls...)

	c, err := client.New("some-access-token", client.WithHTTPClient(httpClient))
	require.NoError(t, err)

	source, err := c.RETLSources.Update(ctx, &client.RETLSource{
		ID:        "some-id",
		Name:      "some-name",
		Type:      "some-type",
		IsEnabled: true,
		Config: json.RawMessage([]byte(`{
			"key1": "val1"
		}`)),
	})
	require.NoError(t, err)
	assert.NotNil(t, source)
	assert.Equal(t, "some-id", source.ID)
	assert.Equal(t, "some-name", source.Name)
	assert.Equal(t, "some-type", source.Type)
	assert.Equal(t, time.Date(2020, 1, 1, 1, 1, 1, 0, time.UTC), *source.CreatedAt)
	assert.Equal(t, time.Date(2020, 1, 2, 1, 1, 1, 0, time.UTC), *source.UpdatedAt)

	httpClient.AssertNumberOfCalls()
}
