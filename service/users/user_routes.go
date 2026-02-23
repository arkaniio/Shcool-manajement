package service

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"github.com/ArkaniLoveCoding/Shcool-manajement/middleware"
	"github.com/ArkaniLoveCoding/Shcool-manajement/middleware/logger"
	"github.com/ArkaniLoveCoding/Shcool-manajement/types"
	"github.com/ArkaniLoveCoding/Shcool-manajement/utils"
)

// this is for router that token is not verified in their function!
type HandleRequest struct {
	db types.UserStore
}

func NewHandlerUser(db types.UserStore) *HandleRequest {
	return &HandleRequest{db: db}
}

// controler that take the services on it
func (h *HandleRequest) Register_Bp(w http.ResponseWriter, r *http.Request) {

	//get request id from this func
	requestID := middleware.GetRequestID(r)
	if requestID == "" {
		//logger the data response if request id is zero value
		logger.Log.Info("Failed to get the request id", 
			zap.String("client_ip", r.RemoteAddr),
			zap.String("path", r.URL.Path),
	)
		utils.ResponseError(w, http.StatusBadRequest, "Failed to detect the request id!", false)
		return 
	}

	//create the payload of the json
	var payload types.Register
	if err := utils.DecodeData(r, &payload); err != nil {
		//logger the data response
		logger.Log.Error("Failed to decode data", 
			zap.String("request_id", requestID),
			zap.String("client_ip", r.RemoteAddr),
	)
		utils.ResponseError(w, http.StatusBadRequest, "Failed to decode the data of the json!", err.Error())
		return 
	}

	//validate the email is the right form email
	is_valid := utils.IsValidEmail(payload.Email)
	if !is_valid {
		//logger the data response if is failed 
		logger.Log.Error("Failed , invalid email type!", 
			zap.String("request_id", requestID),
			zap.String("client_ip", r.RemoteAddr),
	)
		utils.ResponseError(w, http.StatusBadRequest, "Invalid email!", false)
		return
	}

	//validate the error message
	var validate *validator.Validate
	validate = validator.New()
	if err := validate.Struct(&payload); err != nil {
		var errors []string
		for _, validator_payload := range err.(validator.ValidationErrors) {
			errors = append(errors, fmt.Sprintf("Error data %s, %s", validator_payload.Field(), validator_payload.Error()))
			logger.Log.Warn("Validation failed",
				zap.String("request_id", requestID),
				zap.Strings("errors", errors),
			)

			utils.ResponseError(w, http.StatusBadRequest, "Validation error", errors)
			return
		}
	}

	// validate if the email and username has been already exist
	users, err := h.db.GetUserByEmailAndUsername(payload.Email, payload.Username)
	if err != nil {
		//logger the data response if get user by username and by email is vailed!
		logger.Log.Error("Failed to get the username and email", 
			zap.String("request_id", requestID),
			zap.String("client_ip", r.RemoteAddr),
	)
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the email and username", err.Error())
		return 
	}
	if users != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Username and email has been already exist!!", false)
		return
	}

	//hash the password user for a better security 
	hash_password, err := utils.HashPassword(payload.Password)
	if err != nil {
		//logger the data response if hash password is failed
		logger.Log.Error("Failed to hash the password", 
			zap.String("request_id", requestID),
			zap.String("client_ip", r.RemoteAddr),
	)	
		utils.ResponseError(w, http.StatusBadRequest, "Failed to hash the password of the data user!", err.Error())
		return
	}

	//parsing the date time into a string
	time_created := time.Now().UTC().Format("2006-01-02")
	time_updated := time.Now().UTC().Format("2006-01-02")

	//parsing the payload into a struct in the db
	final_payload := &types.User{
		Id: uuid.New(),
		Username: payload.Username,
		Email: payload.Email,
		Password: hash_password,
		Profile_Image: payload.Profile_Image,
		Role: payload.Role,
		Created_at: payload.Created_at,
		Updated_at: payload.Updated_at,
	}

	//create a new user with the registration
	ctx, cancle := context.WithTimeout(r.Context(), time.Second * 10)
	defer cancle()

	//execute the query
	if err := h.db.CreateUser(ctx, final_payload); err != nil {
		//logger the data response if create user is failed
		logger.Log.Error("Failed to create new user", 
			zap.String("request_id", requestID),
			zap.String("client_ip", r.RemoteAddr),
	)	
		utils.ResponseError(w, http.StatusBadRequest, "Failed to create a new user", err.Error())
		return
	}

	//parsing into a user response in types user
	users_response := types.UserResponse{
		Id: final_payload.Id,
		Username: final_payload.Username,
		Email: final_payload.Email,
		Password: final_payload.Password,
		Profile_Image: final_payload.Profile_Image,
		Role: final_payload.Role,
		Created_at: time_created,
		Updated_at: time_updated,
	}

	//return a response success
	utils.ResponseSuccess(w, http.StatusCreated, "Created a new user has been successfully!", users_response)

}

//controller services for handling the router login
func (h *HandleRequest) Login_Bp(w http.ResponseWriter, r *http.Request) {

	//make the request id from this func
	requestID := middleware.GetRequestID(r)
	if requestID == "" {
		logger.Log.Info("Failed to get request!", 
			zap.String("client_ip", r.RemoteAddr),
			zap.String("path", r.URL.Path),
	)
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the request id!", false)
		return 
	}

	//decode the payload of the struct into a json structure
	var payload types.Login
	if err := utils.DecodeData(r, &payload); err != nil {
		//logger the data response if the decode data is failed
		logger.Log.Error("Failed to decode the payload data", 
			zap.String("request_id", requestID),
			zap.String("client_ip", r.RemoteAddr),
	)	
		utils.ResponseError(w, http.StatusBadRequest, "Failed to decode the payload of the login struct!", err.Error())
		return
	}

	//validate the email in the payload (it must be a good structure of the email user)
	is_valid := utils.IsValidEmail(payload.Email)
	if !is_valid {
		//logger the data response if the email is invalid
		logger.Log.Error("Failed , invalid email type!", 
			zap.String("request_id", requestID),
			zap.String("client_ip", r.RemoteAddr),
	)
		utils.ResponseError(w, http.StatusBadRequest, "Invalid email!", false)
		return
	}

	//make the validator of the payload login structure
	var validate *validator.Validate
	validate = validator.New()
	if err := validate.Struct(&payload); err != nil {
		var errors []string
		for _, payload_validator := range err.(validator.ValidationErrors) {
			errors = append(errors, fmt.Sprintf("Error data %s, %s", payload_validator.Field(), payload_validator.Error()))
			logger.Log.Warn("Validation failed",
				zap.String("request_id", requestID),
				zap.Strings("errors", errors),
			)

			utils.ResponseError(w, http.StatusBadRequest, "Validation error", errors)
			return
		}
	}

	//check the email and username (exist or not found)
	users, err := h.db.GetUserByEmailAndUsername(payload.Email, payload.Username) 
	if err != nil {
		//logger the data response if failed to get email and username by users
		logger.Log.Error("Failed to get user by username and email", 
			zap.String("request_id", requestID),
			zap.String("client_ip", r.RemoteAddr),
	)
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get email and username from db!", err.Error())
		return 
	}
	if users == nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get email and username, nill result", false)
		return 
	}
	
	//compare the password
	if err := utils.ComparePassword(users.Password, payload.Password); err != nil {
		//logger the data response if the compare password is failed 
		logger.Log.Error("Failed to compare the password", 
			zap.String("request_id", requestID),
			zap.String("client_ip", r.RemoteAddr),
	)
		utils.ResponseError(w, http.StatusBadRequest, "Failed to compare the password!", err.Error())
		return
	}

	//make the token and refresh token using the payload data
	token, refresh_token, err := utils.GenerateJwt(users.Id, users.Username, users.Email, users.Role)
	if err != nil {
		//logger the data response if the generate jwt is failed
		logger.Log.Error("Failed to generate the jwt !", 
			zap.String("request_id", requestID),
			zap.String("client_ip", r.RemoteAddr),
	)
		utils.ResponseError(w, http.StatusBadRequest, "Failed to create a new token from that function", err.Error())
		return 
	}
	if token == "" && refresh_token == "" {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to detect the token and refresh token!", false)
		return
	}

	//make the response of the login bp
	response_users := make(map[string]interface{})
	response_users["data"] = map[string]interface{}{
		"email": users.Email,
		"username": users.Username,
		"role": users.Role,
		"token": token,
		"refresh_token": refresh_token,
	}

	//return the response is success
	utils.ResponseSuccess(w, http.StatusOK, "Login has been successfully!", response_users)
}

//router to get the profile of the user
func (h *HandleRequest) Profile_Bp(w http.ResponseWriter, r *http.Request) {

	//get the request id from this func
	requestID := middleware.GetRequestID(r)
	if requestID == "" {
		//logger the data response if the request id value is zero
		logger.Log.Info("Failed to get the request id!", 
			zap.String("client_id", r.RemoteAddr),
			zap.String("path", r.URL.Path),
	)
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the request id!", false)
		return 
	}

	//get user id from token
	user_id, err := middleware.GetIdMiddleware(w, r)
	if err != nil {
		//logger the data response if failed to get the user_id!
		logger.Log.Error("Failed to get the user id from token header!", 
			zap.String("request_id", requestID),
			zap.String("client_ip", r.RemoteAddr),
	)
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the user id!", err.Error())
		return 
	}
	if user_id == uuid.Nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the user id!", false)
		return 
	}

	//settings the query from store
	users, err := h.db.GetUserById(user_id)
	if err != nil {
		//logger the data response if get user by id is failed
		logger.Log.Error("Failed to get user by id!", 
			zap.String("request_id", requestID),
			zap.String("client_ip", r.RemoteAddr),
	)
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the user by id!", err.Error())
		return 
	}
	if users == nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to detect the data from db because the nil result!", false)
		return 
	}

	//response the json
	response_user := make(map[string]interface{})
	response_user = map[string]interface{}{
		"username": users.Username,
		"email": users.Email,
		"role": users.Role,
	}

	//return final result
	utils.ResponseSuccess(w, http.StatusOK, "Get profile has been successfully!", response_user)

}

//func that update the users data or profile_image
func (h *HandleRequest) Update_Bp(w http.ResponseWriter, r *http.Request) {

	//get the request id from this func
	requestID := middleware.GetRequestID(r)
	if requestID == "" {
		//logger the data response if the request id value is zero!
		logger.Log.Info("Failed to get the request id!", 
			zap.String("client_ip", r.RemoteAddr),
			zap.String("path", r.URL.Path),
	)
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the request id!", false)
		return 
	}

	//declare the id of the parameters
	vars := mux.Vars(r)
	id := vars["id"]
	user_id, err := uuid.Parse(id)
	if err != nil {
		//logger the data response if the uuid parse is failed
		logger.Log.Error("Failed to convert data into an uuid!", 
			zap.String("request_id", requestID),
			zap.String("client_ip", r.RemoteAddr),
	)
		utils.ResponseError(w, http.StatusBadRequest, "Failed to convert data string into a uuid type!", err.Error())
		return 
	}
	if user_id == uuid.Nil {
		utils.ResponseError(w, http.StatusBadRequest, "Invalid uuid type!", false)
		return 
	}

	//declare the form validaton for the size of the file image
	r.Body = http.MaxBytesReader(w, r.Body, 2 << 20)

	//parse multipart form to setting the request is the form file
	if err := r.ParseMultipartForm(2 << 20); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to parse the multipart form data for a request!", err.Error())
		return 
	}

	//declare the form value for each field in db
	var payload types.Update
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")

	//declare the form file for profile image
	file_image, header, err := r.FormFile("profile_image")
	if err != nil {
		//logger the data response 
		logger.Log.Error("Failed to get the profile image", 
			zap.String("request_id", requestID),
			zap.String("client_ip", r.RemoteAddr),
	)
		if err != http.ErrMissingFile {
			utils.ResponseError(w, http.StatusBadRequest, "Failed to detect profile image!", err.Error())
			return 
		}
	}
	if err == nil {
		//validate the type of the file profile_image in db and read the file
		buff := make([]byte, 512)
		read_buff, err := file_image.Read(buff)
		if err != nil {
			//logger the data response if the read file image is failed
			logger.Log.Error("Failed to read image file", 
				zap.String("request_id", requestID),
				zap.String("client_ip", r.RemoteAddr),
		)
			utils.ResponseError(w, http.StatusBadRequest, "Failed to read the image file!", err.Error())
			return 
		}
		if read_buff == 0 {
			utils.ResponseError(w, http.StatusBadRequest, "Invalid lenght of the file data!", false)
			return
		}
		type_content := http.DetectContentType(buff)
		if type_content != "image/jpg" && type_content != "image/png" && type_content != "image/jpeg" {
			//logger the data response if the type of image is failed
			logger.Log.Error("Failed because the image file is invalid!", 
				zap.String("request_id", requestID),
				zap.String("client_ip", r.RemoteAddr),
		)
			utils.ResponseError(w, http.StatusBadRequest, "Failed content file type!", false)
			return
		}
		file_image.Seek(0, 0)

		//make the name of the file profile image and then join it into a some folder in this project
		filename := uuid.New().String() + filepath.Ext(header.Filename)
		upload_dir := "uploads_user"
		if err := os.MkdirAll(upload_dir, os.ModePerm); err != nil {
			//logger the data response if the make folder is failed
			logger.Log.Error("Failed!", 
				zap.String("request_id", requestID),
				zap.String("client_ip", r.RemoteAddr),
		)
			utils.ResponseError(w, http.StatusBadRequest, "Failed to settings the mkdir all", err.Error())
			return 
		}
		path_final := filepath.Join(upload_dir, filename)
		folder, err := os.Create(path_final)
		if err != nil {
			//logger the data response if create image file is failed
			logger.Log.Error("Failed because the image file is invalid!", 
				zap.String("request_id", requestID),
				zap.String("client_ip", r.RemoteAddr),
		)
			utils.ResponseError(w, http.StatusBadRequest, "Failed to create the folder!", err.Error())
			return 
		}
		if path_final == "" {
			utils.ResponseError(w, http.StatusBadRequest, "Failed to create the folder of uploads!", false)
			return 
		}
		defer folder.Close()

		//copy the result into a io reader
		dst, err := io.Copy(folder, file_image)
		if err != nil {
			//logger the data response if the copy io is failed
			logger.Log.Error("Failed because the lenght of the file is zero!", 
				zap.String("request_id", requestID),
				zap.String("client_ip", r.RemoteAddr),
		)
			utils.ResponseError(w, http.StatusBadRequest, "Failed to copy the data from fileimage!", err.Error())
			return 
		} 
		if dst == 0 {
			utils.ResponseError(w, http.StatusBadRequest, "Failed to copy the data because the length of the data is 0", false)
			return 
		}
		users_profile_image, err := h.db.GetUserById(user_id)
		if err != nil {
			//logger the data response if users is failed
			logger.Log.Error("Failed because id is invalid!", 
				zap.String("request_id", requestID),
				zap.String("client_ip", r.RemoteAddr),
		)
			utils.ResponseError(w, http.StatusBadRequest, "Failed to get the users from db!", err.Error())
			return 
		}

		//if the users have a profile image in their data profile, and then they want to update the profile image again
		//it will be delete the old filename and old path name, that can be replace with new filename and pathname
		//we use this because, if the profile image picture user that we save it into one folder and
		//try to think that if a million users is update their proifile image together, its make server down and worst method
		if users_profile_image.Profile_Image != "" {
			path_old := users_profile_image.Profile_Image
			if _, err := os.Stat(path_old); !os.IsNotExist(err) {
				if err := os.Remove(path_old); err != nil {
					//logger the data response if the file is not exist
				logger.Log.Error("Failed ro remove data!", 
					zap.String("request_id", requestID),
					zap.String("client_ip", r.RemoteAddr),
				)
					utils.ResponseError(w, http.StatusBadRequest, "Failed to remove the new data of uploads!", err.Error())
					return 
				}
			} 
		}
		payload.Profile_Image = &path_final
	}

	// validate if the username, email or password and profile image is nil, it will be return a nil not a string
	// so we can validate it use the conditional validation to make sure that if the email, username or password
	// is not a nill value, it will be passed into a payload data
	if username != "" {
		payload.Username = &username
	}
	if email != "" {
		payload.Email = &email
	}
	if password != "" {
		payload.Password = &password
	}
	
	//settings the context and setup the query
	ctx, cancle := context.WithTimeout(r.Context(), time.Second * 10)
	defer cancle()
	if err := h.db.UpdateDataUser(user_id, ctx, payload); err != nil {
		//logger the data response if the update data is failed
			logger.Log.Error("Failed because the id or data is invalid!!", 
				zap.String("request_id", requestID),
				zap.String("client_ip", r.RemoteAddr),
		)
		utils.ResponseError(w, http.StatusBadRequest, "Failed to update the data user!", err.Error())
		return 
	}

	//get the user data in db
	users, err := h.db.GetUserById(user_id)
	if err != nil {
		//logger the data response if the id is invalid
			logger.Log.Error("Failed because the id is invalid!", 
				zap.String("request_id", requestID),
				zap.String("client_ip", r.RemoteAddr),
		)
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get data user from db!", err.Error())
		return 
	}

	//make the user response
	user_update_response := types.UserResponse{
		Id: users.Id,
		Username: users.Username,
		Email: users.Email,
		Password: users.Password,
		Profile_Image: users.Profile_Image,
		Role: users.Role,
		Created_at: time.Now().UTC().Format("2006-01-02"),
		Updated_at: time.Now().UTC().Format("2006-01-02"),
	}

	//return final results
	utils.ResponseSuccess(w, http.StatusOK, "Update data user has been successfully!", user_update_response)

}

//func to get the profile image picture in postman / url
func (h *HandleRequest) Image_Bp(w http.ResponseWriter, r *http.Request) {

	//get the request id from this func
	requestID := middleware.GetRequestID(r)
	if requestID == "" {
		//logger the data response if the request id value is zero!
		logger.Log.Info("Failed to get the request id!", 
			zap.String("client_ip", r.RemoteAddr),
			zap.String("path", r.URL.Path),
	)
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the request id!", false)
		return 
	}

	//to get the params of the file
	filename_params := mux.Vars(r)
	filename := filename_params["filename"]
	if filename == "" {
		//logger the data response if the params is null
			logger.Log.Error("Failed because the params is invalid!", 
				zap.String("request_id", requestID),
				zap.String("client_ip", r.RemoteAddr),
		)
		utils.ResponseError(w, http.StatusBadRequest, "Failed to create the params of the image!", false)
		return 
	}

	//get the path name of the profile image
	path_name := filepath.Join("uploads_user", filename)
	if _, err := os.Stat(path_name); os.IsNotExist(err) {
		//logger the data response the checking data file is invalid
			logger.Log.Error("Failed because the file is not exist", 
				zap.String("request_id", requestID),
				zap.String("client_ip", r.RemoteAddr),
		)
		utils.ResponseError(w, http.StatusBadRequest, "Failed to detect the file or path name", err.Error())
		return 
	}

	//serve http for file
	http.ServeFile(w, r, path_name)

}