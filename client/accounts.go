package client

import (
	"context"
	"encoding/json"
	"time"
)

type Account struct {
	ID          string          `json:"id,omitempty"`
	UserID      string          `json:"userId,omitempty"`
	WorkspaceID string          `json:"workspaceId,omitempty"`
	Name        string          `json:"name"`
	Type        string          `json:"type"`
	Category    string          `json:"category"`
	Options     json.RawMessage `json:"options"`
	CreatedAt   *time.Time      `json:"createdAt,omitempty"`
	UpdatedAt   *time.Time      `json:"updatedAt,omitempty"`
}

type accounts struct {
	*service
}

type AccountsPage struct {
	APIPage
	Accounts []Account `json:"accounts"`
}

func (s *accounts) Next(ctx context.Context, paging Paging) (*AccountsPage, error) {
	page := &AccountsPage{}
	ok, err := s.service.next(ctx, paging, page)
	if !ok {
		page = nil
	}
	return page, err
}

func (s *accounts) List(ctx context.Context) (*AccountsPage, error) {
	page := &AccountsPage{}
	if err := s.list(ctx, page); err != nil {
		return nil, err
	}

	return page, nil
}

func (s *accounts) Get(ctx context.Context, id string) (*Account, error) {
	response := struct{ Account *Account }{}
	if err := s.get(ctx, id, &response); err != nil {
		return nil, err
	}

	return response.Account, nil
}

func (s *accounts) Create(ctx context.Context, account *Account) (*Account, error) {
	// copy input and remove fields that should not be in request body without modifying input
	act := *account
	act.ID = ""

	response := struct{ Account *Account }{}
	if err := s.create(ctx, &act, &response); err != nil {
		return nil, err
	}

	return response.Account, nil
}

func (s *accounts) Update(ctx context.Context, account *Account) (*Account, error) {
	// copy input and remove ID from request body without modifying input
	act := *account
	act.ID = ""

	response := struct{ Account *Account }{}
	if err := s.update(ctx, account.ID, &act, &response); err != nil {
		return nil, err
	}

	return response.Account, nil
}

func (s *accounts) Delete(ctx context.Context, id string) error {
	return s.service.delete(ctx, id)
}
