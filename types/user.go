package types

import (
	"context"
	"time"

	"github.com/google/uuid"
)


type UserStore interface {
	GetUserByEmailAndUsername(email string, username string) (*User, error)
	CreateUser(ctx context.Context, user *User) error 
	UpdateDataUser(
		id uuid.UUID,
		ctx context.Context,
		payload Update,
		) error
	GetUserById(id uuid.UUID) (*User, error)
}

type User struct {
	Id				uuid.UUID	`db:"id"`
	Username 		string  	`db:"username"`
	Email 			string  	`db:"email"`
	Password 		string 		`db:"password"`
	Profile_Image 	string 		`db:"profile_image"`
	Role 			string 		`db:"role"`
	Created_at 		time.Time 	`db:"created_at"`
	Updated_at		time.Time 	`db:"updated_at"`
}

type Register struct {
	Id 				uuid.UUID	`json:"id"`
	Username 		string 		`json:"username" validate:"required,min=2,max=100"`
	Email 			string 		`json:"email" validate:"required,email,min=2,max=100"`
	Password 		string 		`json:"password" validate:"required,min=2,max=100"`
	Profile_Image 	string 		`json:"profile_image"`
	Role 			string 		`json:"role"`
	Created_at 		time.Time 	`json:"created_at"`
	Updated_at 		time.Time 	`json:"updated_at"`
}

type Login struct {
	Username 		string 	`json:"username" validate:"required"`
	Email 			string	`json:"email" validate:"required,email"`
	Password 		string	`json:"password" validate:"required,min=2,max=100"`
}

type UserResponse struct {
	Id				uuid.UUID	`json:"id"`
	Username 		string 		`json:"username"`
	Email 			string 		`json:"email"`
	Password 		string 		`json:"password"`
	Profile_Image   string 		`json:"profile_image"`
	Role 			string 		`json:"role"`
	Created_at 		string  	`json:"created_at"`
	Updated_at 		string 		`json:"updated_at"`
}

type Update struct {
	Username 		*string 	`json:"username"`
	Email 			*string  	`json:"email"`
	Password 		*string 	`json:"password"`
	Profile_Image 	*string 	`json:"profile_image"`
	Updated_at      *string  	`json:"updated_at"`
}
