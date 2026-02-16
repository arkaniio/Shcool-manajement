package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/ArkaniLoveCoding/Shcool-manajement/types"
)

type Store struct {
	store *sqlx.DB
}

func NewStore(store *sqlx.DB) *Store {
	return &Store{store: store}
}

func (s *Store) GetUserByEmail(email string) (*types.User, error) {
	var user types.User
	query := `SELECT 
	id, firstname, lastname, password, email, country, address, role, token, refresh_token, created_at, updated_at
	FROM users WHERE email = $1;`
	err := s.store.Get(&user, query, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return &user, nil
}

func (s *Store) GetUserById(id uuid.UUID) (*types.User, error) {

	var user types.User
	err := s.store.Get(&user, "SELECT id, firstname, lastname, password, email, country, address, role, token, refresh_token, created_at, updated_at FROM users WHERE id = $1", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}
	return &user, nil

}

func (s *Store) CreateUser(ctx context.Context, user *types.User) error  {

	tx_options := &sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly: false,
	}
	ctx, cancle := context.WithTimeout(context.Background(), time.Second * 10)
	defer cancle()

	tx, err := s.store.BeginTxx(ctx, tx_options)
	if err != nil {
		return errors.New("Failed to doing transactions!")
	}

	defer tx.Rollback()
	
	query := `
		INSERT INTO users (id, firstname, lastname, 
		password, email, country, address, role, token, refresh_token, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING *;
	`

	if err := tx.QueryRowContext(
		ctx,
		query,
		user.Id,
		user.Firstname,
		user.Lastname,
		user.Password,
		user.Email,
		user.Country,
		user.Address,
		user.Role,
		user.Token,
		user.Rerfresh_token,
		user.Created_at,
		user.Updated_at,
	).Scan(
		&user.Id,
		&user.Firstname,
		&user.Lastname,
		&user.Password,
		&user.Email,
		&user.Country,
		&user.Address,
		&user.Role,
		&user.Token,
		&user.Rerfresh_token,
		&user.Created_at,
		&user.Updated_at,
		); err != nil {
		return nil
	}

	if err := tx.Commit(); err != nil {
		return errors.New("Failed to commit the transaction!")
	}

	return nil
}

func (s *Store) UpdateDataUser(
	id uuid.UUID, 
	ctx context.Context, 
	firstname string,
	lastname string,
	password string,
	email string,
	country string,
	address string,
	users *types.User,
	) error {

	tx_options := &sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly: false,
	}
	ctx, cancle := context.WithTimeout(context.Background(), time.Second * 10)
	defer cancle()

	tx, err := s.store.BeginTxx(ctx, tx_options)
	if err != nil {
		return errors.New("Failed to doing transactions!")
	}

	defer tx.Rollback()

	query := `
		UPDATE users 
		SET firstname = $2,
			lastname = $3,
			password = $4,
			email = $5,
			country = $6,
			address = $7
		WHERE id = $1;
	`
	var u = users

	if err := tx.QueryRowContext(
		ctx,
		query,
		id,
		firstname,
		lastname,
		password,
		email,
		country,
		address,
	).Scan(
		&u.Id,
		&u.Firstname,
		&u.Lastname,
		&u.Password,
		&u.Email,
		&u.Country,
		&u.Address,
	); err != nil {
		return nil
	}

	if err := tx.Commit(); err != nil {
		return errors.New("Failed to commit the transactions!")
	}

	return nil

}


