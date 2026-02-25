package students

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/ArkaniLoveCoding/Shcool-manajement/middleware"
	"github.com/ArkaniLoveCoding/Shcool-manajement/middleware/logger"
	"github.com/ArkaniLoveCoding/Shcool-manajement/types"
	"github.com/ArkaniLoveCoding/Shcool-manajement/utils"
)

//type handlerequest that declare the student store for a database logic
type HandleRequest struct{
	db types.StudentStore
}

//func that declare the handler for student
func NewHandlerStudent(db types.StudentStore) *HandleRequest {
	return &HandleRequest{db: db}
}

//func to create a new student
func (h *HandleRequest) RegisterAsStudent_Bp(w http.ResponseWriter, r *http.Request) {

	//get the request id from this func
	requestID := middleware.GetRequestID(r)
	if requestID == "" {
		//make the logger data response for info
		logger.Log.Info("Failed to get the request id from this func!", 
			zap.String("client_ip", r.RemoteAddr),
			zap.String("path", r.URL.Path),
	)
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the request id!", false)
		return 
	}

	//validate the role, cannot register as a student if the role of the user is (siswa)
	role, err := middleware.GetRoleMiddleware(w, r)
	if err != nil {
		//logger if the error is detected
		logger.Log.Error("Failed to get the middleware for role", 
			zap.String("request_id", requestID),
			zap.String("client_ip", r.RemoteAddr),
	)
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the middleware role", err.Error())
		return 
	}
	if role == "guru" {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to access this method!", false)
		return 
	}

	//decode the payload of the struct student register
	var payload types.RegisterAsStudent
	if err := utils.DecodeData(r, &payload); err != nil {
		//make the data response for logger if the decode is failed
		logger.Log.Error("Failed to decode data payload", 
			zap.String("request_id", requestID),
			zap.String("client_ip", r.RemoteAddr),
	)
		utils.ResponseError(w, http.StatusBadRequest, "Failed to decode the data!", err.Error())
		return 
	}

	//make the validator of the payload
	var validate *validator.Validate
	validate = validator.New()
	if err := validate.Struct(&payload); err != nil {
		var errors []string
		for _, erorrValidate := range err.(validator.ValidationErrors) {
			errors = append(errors,fmt.Sprintf("error at field: %s, %s", erorrValidate.Field(), erorrValidate.Error()))
			//make the response data if the validate is failed to detect some error in field
			logger.Log.Error("Failed to doing some validate", 
				zap.String("request_id", requestID),
				zap.String("client_ip", r.RemoteAddr),
		)
			utils.ResponseError(w, http.StatusBadRequest, "Failed to make the validator!", err.Error())
			return 
			}	
		}

	//checking if the student is not have a same name
	students, err := h.db.GetStudentByName(payload.Name)
	if err != nil {
		//logger if some error is detected
		logger.Log.Error("Failed to get the student", 
			zap.String("request_id", requestID),
			zap.String("client_ip", r.RemoteAddr),
	)
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the name of student!", err.Error())
		return 
	}
	if students != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Name has been already exist!", false)
		return 
	}

	//define the time updated and created response
	time_updated_format := time.Now().UTC().Format("2006-01-02")
	time_created_format := time.Now().UTC().Format("2006-01-02")

	//make the struct of payload to interact with the struct of the user
	students_payload := &types.Student{
		Id: uuid.New(),
		Name: payload.Name,
		Class: payload.Class,
		Address: payload.Address,
		Major: payload.Major,
		StudentProfile: payload.StudentProfile,
		Created_at: payload.Created_at,
		Updated_at: payload.Updated_at,
	}

	//declare the context to user
	ctx, cancle := context.WithTimeout(r.Context(), time.Second * 10)
	defer cancle()

	//execute the query of the create user
	if err := h.db.CreateNewStudent(ctx, students_payload); err != nil {
		//logger if some error is detected when we want to create it
		logger.Log.Error("Failed to create a new student", 
			zap.String("request_id", requestID),
			zap.String("client_ip", r.RemoteAddr),
	)
		utils.ResponseError(w, http.StatusBadRequest, "Failed to create the students data!", err.Error())
		return 
	}

	//make the response of the students data
	students_response := types.StudentResponse{
		Id: students_payload.Id,
		Name: students_payload.Name,
		Class: students_payload.Class,
		Address: students_payload.Address,
		Major: students_payload.Major,
		StudentProfile: students_payload.StudentProfile,
		Created_at: time_created_format,
		Updated_at: time_updated_format,
	}
	
	//return a final value
	utils.ResponseSuccess(w, http.StatusCreated, "Register as a student has been successfully", students_response)

}

func (h *HandleRequest) GetAll_Bp(w http.ResponseWriter, r *http.Request) {

	//get the request id from this func
	requestID := middleware.GetRequestID(r)
	if requestID == "" {
		//make the logger data response for info
		logger.Log.Info("Failed to get the request id from this func!", 
			zap.String("client_ip", r.RemoteAddr),
			zap.String("path", r.URL.Path),
	)
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the request id!", false)
		return 
	}

	//get the middleware token to validate this method cannot see the other student
	role_user, err := middleware.GetRoleMiddleware(w, r)
	if err != nil {
		//logger the data response if something error in this method
		logger.Log.Error("Failed to get the role middleware", 
			zap.String("request_id", requestID),
			zap.String("client_ip", r.RemoteAddr),
	)
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the role middleware!", err.Error())
	}
	if role_user == "siswa" {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to access this method!", false)
		return 
	}

	//define the query params 
	limit := r.URL.Query().Get("limit")
	sort := r.URL.Query().Get("sort")
	order := r.URL.Query().Get("order")
	cursor := r.URL.Query().Get("cursor")

	//convert limit into an integer
	limit_convert, err := strconv.Atoi(limit)
	if err != nil {
		//logger the data response if something went wrong with this method
		logger.Log.Error("Failed to conver the data!", 
			zap.String("request_id", requestID),
			zap.String("client_ip", r.RemoteAddr),
	)
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the limit convert!", err.Error())
		return 
	}

	//validate the limit
	if limit_convert > 0 && limit_convert < 50 {
		limit_convert = 10
	}

	//decode the value of the cursor
	var valueCursor any
	var IdCursor string
	decode, err := utils.DecodeCursor(cursor)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the cursor decode", err.Error())
		return 
	}
	if decode != nil {
		t, err := time.Parse(time.RFC3339, decode.Value)
		if err != nil {
			utils.ResponseError(w, http.StatusBadRequest, "Failed to get the time!", err.Error())
			return 
		}
		valueCursor = t
		IdCursor = decode.Id
	}

	//execute the query
	ctx, cancle := context.WithTimeout(r.Context(), time.Second * 10)
	defer cancle()
	students, err := h.db.GetAllStudents(ctx, limit_convert, sort, order, valueCursor, IdCursor)
	if err != nil {
		//logger if the response is nill 
		logger.Log.Error("Failed to get all the data student", 
			zap.String("request_id", requestID),
			zap.String("client_ip", r.RemoteAddr),
	)
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the all of the students", err.Error())
		return 
	}

	//for the next cursor
	var nextCursor *string
	if len(students) > 0 {
		last_students := students[len(students) - 1]
		encode, err := utils.EncodeCursor(last_students.Created_at, last_students.Id.String()) 
		if err == nil {
			nextCursor = &encode
		}
	}

	//make the struct for the data students
	response_user := map[string]interface{}{
		"data_students": students,
		"next_cursor": nextCursor,
	}
	
	//return a final result
	utils.ResponseSuccess(w, http.StatusOK, "Get alll students has been successfully", response_user)

}
