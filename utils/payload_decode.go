package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)


func DecodeData (r *http.Request, v any) error  {

	if r.Header.Get("Content-Type") != "application/json" {
		return fmt.Errorf("content-type must be application/json")
	}

	r.Body = http.MaxBytesReader(nil, r.Body, 1<<20) // 1MB limit

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(v); err != nil {
		return err
	}

	return nil
}