package client

import (
	"context"
	"encoding/json"
	"time"
)

type RETLSource struct {
	ID          string          `json:"id,omitempty"`
	Name        string          `json:"name"`
	Type        string          `json:"type"`
	SourceType  string          `json:"sourceType"`
	AccountID   string          `json:"accountId"`
	WorkspaceID string          `json:"workspaceId"`
	IsEnabled   bool            `json:"enabled"`
	Config      json.RawMessage `json:"config"`
	CreatedAt   *time.Time      `json:"createdAt,omitempty"`
	UpdatedAt   *time.Time      `json:"updatedAt,omitempty"`
}

type retlSources struct {
	*service
}

type RETLSourcesPage struct {
	APIPage
	RETLSources []RETLSource `json:"sources"`
}

func (s *retlSources) Next(ctx context.Context, paging Paging) (*RETLSourcesPage, error) {
	page := &RETLSourcesPage{}
	ok, err := s.service.next(ctx, paging, page)
	if !ok {
		page = nil
	}
	return page, err
}

func (s *retlSources) List(ctx context.Context) (*RETLSourcesPage, error) {
	page := &RETLSourcesPage{}
	if err := s.list(ctx, page); err != nil {
		return nil, err
	}

	return page, nil
}

func (s *retlSources) Get(ctx context.Context, id string) (*RETLSource, error) {
	response := struct{ Source *RETLSource }{}
	if err := s.get(ctx, id, &response); err != nil {
		return nil, err
	}

	return response.Source, nil
}

func (s *retlSources) Create(ctx context.Context, retlSource *RETLSource) (*RETLSource, error) {
	// copy input and remove fields that should not be in request body without modifying input
	src := *retlSource
	src.ID = ""

	response := struct{ Source *RETLSource }{}
	if err := s.create(ctx, &src, &response); err != nil {
		return nil, err
	}

	return response.Source, nil
}

func (s *retlSources) Update(ctx context.Context, retlSource *RETLSource) (*RETLSource, error) {
	// copy input and remove ID from request body without modifying input
	src := *retlSource
	src.ID = ""

	response := struct{ Source *RETLSource }{}
	if err := s.update(ctx, retlSource.ID, &src, &response); err != nil {
		return nil, err
	}

	return response.Source, nil
}

func (s *retlSources) Delete(ctx context.Context, id string) error {
	return s.service.delete(ctx, id)
}
