package students

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/ArkaniLoveCoding/Shcool-manajement/types"
)

//type for a store student
type StudentStore struct {
	db *sqlx.DB
}

//func that we use when we want to use the store from this db
func NewStudentStore(db *sqlx.DB) *StudentStore {
	return &StudentStore{db: db}
}

//func that create a new student
func (s *StudentStore) CreateNewStudent(ctx context.Context, student *types.Student) error {

	//make the options of transaction of create
	options := &sql.TxOptions{
		ReadOnly: false,
		Isolation: sql.LevelSerializable,
	}

	//make the ctx for context in transaction 
	ctx, cancle := context.WithTimeout(context.Background(), time.Second * 10)
	defer cancle()

	//make the new begintxx for transaction
	tx, err := s.db.BeginTxx(ctx, options)
	if err != nil {
		return errors.New("Failed to settings the db transactions")
	}
	defer tx.Rollback()

	//make the base query for create a new student
	query := `
		INSERT INTO students 
		(id, name, class, address, major, student_profile, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING *;
	`

	//make the method query
	if err := tx.QueryRowContext(
		ctx,
		query, 
		student.Id,
		student.Name,
		student.Class,
		student.Address,
		student.Major,
		student.StudentProfile,
		student.Created_at,
		student.Updated_at,
	).Scan(
		&student.Id,
		&student.Name,
		&student.Class,
		&student.Address,
		&student.Major,
		&student.StudentProfile,
		&student.Created_at,
		&student.Updated_at,
		); err != nil {
			return errors.New("Failed to create a new user!" + err.Error())
		}

	//commit the transaction
	if err := tx.Commit(); err != nil {
		return errors.New("Failed to commit the query of transaction!" + err.Error())
	}

	//final return
	return nil

}

//func to get the name in student db
func (s *StudentStore) GetStudentByName(name string) (*types.Student, error) {

	//make base query
	query := `
		SELECT id, name, class, major, student_profile, created_at, updated_at 
		FROM students WHERE name = $1;
	`

	//make the method of query
	var students types.Student
	if err := s.db.Get(&students, query, name); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return &students, nil

}