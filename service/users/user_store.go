package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/ArkaniLoveCoding/Shcool-manajement/types"
	"github.com/ArkaniLoveCoding/Shcool-manajement/utils"
)

type Store struct {
	store *sqlx.DB
}

func NewStore(store *sqlx.DB) *Store {
	return &Store{store: store}
}

//func get user by email and username 
func (s *Store) GetUserByEmailAndUsername(email string, username string) (*types.User, error) {

	//declare the user
	var user types.User

	//base query for select method
	query := `SELECT 
	id, username, email, password, profile_image, role, created_at, updated_at
	FROM users WHERE email = $1 AND username = $2;`

	//second base queries
	err := s.store.Get(&user, query, email, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	//return final result
	return &user, nil
}

//func create a new user
func (s *Store) CreateUser(ctx context.Context, user *types.User) error  {

	//use transactions options
	tx_options := &sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly: false,
	}

	//declare the context for setup base query
	ctx, cancle := context.WithTimeout(context.Background(), time.Second * 10)
	defer cancle()

	//setup the transaction
	tx, err := s.store.BeginTxx(ctx, tx_options)
	if err != nil {
		return errors.New("Failed to doing transactions!")
	}
	defer tx.Rollback()
	
	//base query
	query := `
		INSERT INTO users (id, username, email, password, profile_image, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING *;
	`

	//second base queries
	if err := tx.QueryRowContext(
		ctx,
		query,
		user.Id,
		user.Username,
		user.Email,
		user.Password,
		user.Profile_Image,
		user.Role,
		user.Created_at,
		user.Updated_at,
	).Scan(
		&user.Id,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Profile_Image,
		&user.Role,
		&user.Created_at,
		&user.Updated_at,
		); err != nil {
		return nil
	}

	//commit if the transaction has been successfully
	if err := tx.Commit(); err != nil {
		return errors.New("Failed to commit the transaction!")
	}

	return nil
}

//func update the users identity
func (s *Store) UpdateDataUser(
	id uuid.UUID, 
	ctx context.Context, 
	payload types.Update,
	) error {

	//settings the options for a transaction
	tx_options := &sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly: false,
	}

	//declare the context for a query
	ctx, cancle := context.WithTimeout(context.Background(), time.Second * 10)
	defer cancle()

	//setup the transaction
	tx, err := s.store.BeginTxx(ctx, tx_options)
	if err != nil {
		return errors.New("Failed to doing transactions!")
	}
	defer tx.Rollback()

	//setup the args and args id
	var	settings []string
	argsId := 1
	var args []interface{}

	//if the users wants to update their username
	if payload.Username != nil {
		settings = append(settings, fmt.Sprintf("username=$%d", argsId))
		args = append(args, *payload.Username)
		argsId++
	}

	//if the users wants to update their email
	if payload.Email != nil {
		settings = append(settings, fmt.Sprintf("email=$%d", argsId))
		args = append(args, *payload.Email)
		argsId++
	}

	//if the users wants to update their password 
	if payload.Password != nil {
		hash_password, err := utils.HashPassword(*payload.Password)
		if err != nil {
			return nil
		}
		settings = append(settings, fmt.Sprintf("password=$%d", argsId))
		args = append(args, hash_password)
		argsId++
	}

	//if the users wants to update their profile_image
	if payload.Profile_Image != nil {
		settings = append(settings, fmt.Sprintf("profile_image=$%d", argsId))
		args = append(args, *payload.Profile_Image)
		argsId++
	}

	//validate if the no one field changes
	if len(args) == 0 {
		return errors.New("No one data changes")
	}

	//update the updated at
	settings = append(settings, fmt.Sprintf("updated_at=$%d", argsId))
	args = append(args, time.Now().UTC())
	argsId++

	//define a fullquery for this db
	fullquery := fmt.Sprintf("UPDATE users SET %s WHERE id = $%d", strings.Join(settings, ", "), argsId)
	args = append(args, id)

	//scan into a rows user
	rows, err := tx.ExecContext(ctx, fullquery, args...)
	if err != nil {
		return errors.New(err.Error())
	}
	
	//checking the rows of the db 
	result, err := rows.RowsAffected()
	if err != nil {
		return errors.New("No one changes in db, error: " + err.Error())
	}
	if result == 0 {
		return errors.New("Failed to scan!")
	}

	//commit the transaction if the transaction has been successfully
	if err := tx.Commit(); err != nil {
		return errors.New("Failed to commit the transactions!" + err.Error())
	}

	return nil

}

//func get user by id
func (s *Store) GetUserById(id uuid.UUID) (*types.User, error) {

	//declare the user
	var users types.User

	//setup the base query
	query := `
		SELECT id, username, email, password, profile_image, role, created_at, updated_at
		FROM users WHERE id = $1
	`
	if query == "" {
		return nil, errors.New("The query is nil")
	}

	//second base query
	if err := s.store.Get(&users, query, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("Failed to get the data because is nil!")
		}
		return nil, errors.New("Failed to get the data user because the nil")
	}

	//return fiinal result
	return &users, nil

}
