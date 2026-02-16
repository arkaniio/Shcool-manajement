package service

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"github.com/ArkaniLoveCoding/School-manajement/types"
)

type mockStore struct {
	GetUserByEmailFn func(email string) (*types.User, error)
	CreateUserFn func(ctx context.Context, user *types.User) error
	UpdateUserFn func(
		id uuid.UUID,
		ctx context.Context,
		firstname string,
		lastname string,
		password string,
		email string,
		country string,
		address string,
	) error
}

func (m *mockStore) GetUserByEmail(email string) (*types.User, error) {
	
	return m.GetUserByEmailFn(email)

}

func (m *mockStore) CreateUser(ctx context.Context, user *types.User) error {
	
	return m.CreateUserFn(ctx, user)

}

func (m *mockStore) UpdateDataUser(
	id uuid.UUID,
	ctx context.Context,
	firstname string,
	lastname string,
	password string,
	email string,
	country string,
	address string,
	user *types.User, 
	) error {
	
		return m.UpdateUserFn(id, ctx, firstname, lastname, password, email, country, address)

}
