package utils

import (
	"encoding/json"
	"errors"
	"net/http"
)


type JsonParams struct {
	Message 	string 		`json:"message"`
	Success 	bool		`json:"success"`
	Data 		interface{}	`json:"data"`
}

func ResponseJson (w http.ResponseWriter, code int, data interface{}) error {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	encode := json.NewEncoder(w)

	if err := encode.Encode(data); err != nil {
		return errors.New(err.Error())
	}

	return nil

}

func responseSuccess (w http.ResponseWriter, code int, message string, data interface{}) error {

	response_json := JsonParams{
		Message: message,
		Success: true,
		Data: data,
	}

	if err := ResponseJson(w, code, response_json); err != nil {
		return nil
	}

	return nil

}

func responseError (w http.ResponseWriter, code int, message string, data interface{}) error {

	response_json := JsonParams{
		Message: message,
		Success: false,
		Data: data,
	}

	if err := ResponseJson(w, code, response_json); err != nil {
		return nil
	}

	return nil

}