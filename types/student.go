package types

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type StudentStore interface {
	CreateNewStudent(ctx context.Context, student *Student) error
	GetStudentByName(name string) (*Student, error)
	GetAllStudents(
		ctx context.Context,
		limit int,
		sort string,
		order string,
		cursorValue any,
		cursorID string,
		) ([]Student, error)
}

type Student struct {
	Id 				uuid.UUID 		`db:"id"`
	Name 			string 			`db:"name"`
	Class 			string 			`db:"class"`
	Address 		string			`db:"address"`
	Major 			string 			`db:"major"`
	StudentProfile	string 			`db:"student_profile"`
	Created_at 		time.Time 		`db:"created_at"`
	Updated_at      time.Time 		`db:"updated_at"`
}

type RegisterAsStudent struct {
	Id 				uuid.UUID 		`json:"id"`
	Name 			string 			`json:"name" validate:"required"`
	Class 			string 			`json:"class" validate:"required"`
	Address 		string 			`json:"address" validate:"required"`
	Major 			string 			`json:"major" validate:"required"`
	StudentProfile 	string 			`json:"student_profile"`
	Created_at 		time.Time 		`json:"created_at"`
	Updated_at 		time.Time 		`json:"updated_at"`
}

type UpdateAsStudent struct {
	Id 				uuid.UUID 		`json:"id"`
	Name 			string 			`json:"name"`
	Class 			string 			`json:"class"`
	Address 		string 			`json:"address"`
	Major 			string 			`json:"major"`
	StudentProfile 	string 			`json:"student_profile"`
	Created_at 		time.Time 		`json:"created_at"`
	Updated_at 		time.Time 		`json:"updated_at"`
}

type StudentResponse struct {
	Id 				uuid.UUID 		`json:"id"`
	Name 			string 			`json:"name"`
	Class 			string 			`json:"class"`
	Address 		string 			`json:"address"`
	Major 			string 			`json:"major"`
	StudentProfile 	string 			`json:"student_profile"`
	Created_at 		string			`json:"created_at"`
	Updated_at 		string 			`json:"updated_at"`
}