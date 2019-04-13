package helpers

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type SignedMessage struct {
	PublicKey	string	`json:"pk"`
}

func (m *SignedMessage) IsEqual(another *SignedMessage) bool {
	if another == nil {
		return false
	}

	return m.PublicKey == another.PublicKey
}

type SignedRequestData struct {
	Message		string	`json:"message"`
	PublicKey	string	`json:"pk"`
	Signature	string	`json:"sign"`
}

func GetSignedRequestData(ctx *gin.Context) (*SignedRequestData, error) {
	var data SignedRequestData
	if err := ctx.ShouldBindJSON(&data); err != nil {
		return nil, err
	}

	return data.check()
}

func GetSignedPayloadData(payload *json.RawMessage) (*SignedRequestData, error) {
	if payload == nil {
		return nil, errors.New("nil payload")
	}

	var data SignedRequestData
	var bytesPayload = []byte(*payload)

	if len(bytesPayload) == 0 {
		return nil, errors.New("empty payload")
	}

	if err := json.Unmarshal([]byte(bytesPayload), &data); err != nil {
		return nil, err
	}

	return data.check()
}

func (rd *SignedRequestData) check() (*SignedRequestData, error) {
	if err := rd.validate(); err != nil {
		return nil, err
	} else {
		return rd, nil
	}
}

func (rd *SignedRequestData) validate() (error) {
	v, err := NewECSDAValidator(rd.PublicKey)
	if err != nil {
		return err
	}

	return v.ValidateString(rd.Message, rd.Signature)
}

func (rd *SignedRequestData) Validate(message *SignedMessage) bool {
	if message == nil {
		return false
	}

	return rd.PublicKey == message.PublicKey
}