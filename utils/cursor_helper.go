package utils

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

type Cursor struct {
	Value 		string 	`json:"cursor"`
	Id 			string 	`json:"id"`
}

func EncodeCursor(value any, id string) (string, error) {

	c := Cursor{
		Value: fmt.Sprintf("%v", value),
		Id: id,
	}

	data, err := json.Marshal(c)
	if err != nil {
		return "", nil
	}

	return base64.StdEncoding.EncodeToString(data), nil

}

func DecodeCursor(encoding string) (*Cursor, error) {

	if encoding == "" {
		return nil, fmt.Errorf("Failed to get the encoding data!")
	}

	data, err := base64.StdEncoding.DecodeString(encoding)
	if err != nil {
		return nil, fmt.Errorf("Failed to decode the data!")
	}

	var c Cursor
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, fmt.Errorf("Failed to decode the data")
	}

	return &c, nil

}