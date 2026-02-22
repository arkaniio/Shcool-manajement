package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"github.com/ArkaniLoveCoding/Shcool-manajement/utils"
)

//middleware bearer token
func TokenIdMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		//get the authorization from header token
		header := r.Header.Get("Authorization")
		if header == "" {
			utils.ResponseError(w, http.StatusBadRequest, "The header of the token is nil!", false)
			return
		}

		//string prefix to get a header token
		token := strings.TrimPrefix(header, "Bearer ")
		if token == "" {
			utils.ResponseError(w, http.StatusBadRequest, "The token is nil!", false)
			return
		}

		//validate the token to make sure the token was added
		token_validate, err := utils.ValidateToken(token)
		if err != nil {
			utils.ResponseError(w, http.StatusBadRequest, "Failed to validate the token!", err.Error())
			return
		}
		if token_validate == nil {
			utils.ResponseError(w, http.StatusBadRequest, "Failed to get the data validate!", false)
			return
		}

		//get the if from tokenvalidate
		user_id, err := uuid.Parse(token_validate.Id)
		if err != nil {
			utils.ResponseError(w, http.StatusBadRequest, "Failed to convert into an uuid!", err.Error())
			return
		}

		//save the data user id to context
		user_id_ctx := context.WithValue(r.Context(), "user_id", user_id)
		if user_id_ctx == nil {
			utils.ResponseError(w, http.StatusBadRequest, "Failed to settings the context value in request client!", false)
			return
		}
		r = r.WithContext(user_id_ctx)

		//save the data of the role user from token
		role_user_ctx := context.WithValue(r.Context(), "role_user", token_validate.Role)
		if role_user_ctx == nil {
			utils.ResponseError(w, http.StatusBadRequest, "Failed to settings the context value in request client", false)
			return 
		}
		r = r.WithContext(role_user_ctx)

		//next http handler
		next.ServeHTTP(w, r)

	})
}

//func that containt the user id value from token jwt
func GetIdMiddleware(w http.ResponseWriter, r *http.Request) (uuid.UUID, error) {

	//get token from id context
	user_id := r.Context().Value("user_id")
	if user_id == "" && user_id == nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the user id from token jwt!", false)
		return uuid.Nil, nil
	}
	
	//convert into a uuid value
	uuid_user, ok := user_id.(uuid.UUID)
	if !ok {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to convert from string into an uuid!", ok)
		return uuid.Nil, nil
	}
	if uuid_user == uuid.Nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to load the uuid user!", false)
		return uuid.Nil, nil
	}

	//return the final value in uuid type of user id
	return uuid_user, nil

}