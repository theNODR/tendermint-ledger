package helpers

import (
	"encoding/json"
)

type ResponseData struct {
	Message		string	`json:"message"`
	Signature	[]byte	`json:"sign"`
}

func NewResponseDataBytes(bytes []byte, pair KeyPair) (*ResponseData, error) {
	str := string(bytes[:])

	return NewResponseDataString(str, pair)
}

func NewResponseDataString(str string, pair KeyPair) (*ResponseData, error) {
	signature, err := pair.SignString(str)

	if err != nil {
		return nil, err
	}

	return &ResponseData{
		Message: str,
		Signature: signature,
	}, nil
}

func NewResponseDataInterface(data interface{}, pair KeyPair) (*ResponseData, error) {
	bytes, err := json.Marshal(data)

	if err != nil {
		return nil, err
	}

	return NewResponseDataBytes(bytes, pair)
}
