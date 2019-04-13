package helpers

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type BiSignedMessage struct {
	FromPublicKey	string	`json:"fromPK"`
	ToPublicKey		string	`json:"toPK"`
}

func (m *BiSignedMessage) IsEqual(another *BiSignedMessage) bool {
	if another == nil {
		return false
	}

	return m.FromPublicKey == another.FromPublicKey &&
		m.ToPublicKey == another.ToPublicKey
}

type BiSignedRequestData struct {
	Message			string	`json:"message"`
	FromPublicKey	string	`json:"fromPK"`
	FromSignature	string	`json:"fromSign"`
	ToPublicKey		string	`json:"toPK"`
	ToSignature		string	`json:"toSign"`
}

func GetBiSignedRequestData(ctx *gin.Context) (*BiSignedRequestData, error) {
	var data BiSignedRequestData
	if err := ctx.ShouldBindJSON(&data); err != nil {
		return nil, err
	}

	return data.check()
}

func GetBiSignedPayloadData(payload *json.RawMessage) (*BiSignedRequestData, error) {
	if payload == nil {
		return nil, errors.New("nil payload")
	}

	var data BiSignedRequestData
	var bytesPayload = []byte(*payload)

	if len(bytesPayload) == 0 {
		return nil, errors.New("empty payload")
	}

	if err := json.Unmarshal([]byte(*payload), &data); err != nil {
		return nil, err
	}

	return data.check()
}

func (rd *BiSignedRequestData) check() (*BiSignedRequestData, error) {
	if err := rd.validate(); err != nil {
		return nil, err
	} else {
		return rd, nil
	}
}

func (rd *BiSignedRequestData) validate() (error) {
	err := rd.validateBy(rd.FromPublicKey, rd.FromSignature)

	if err != nil {
		return err
	}

	err = rd.validateBy(rd.ToPublicKey, rd.ToSignature)

	if err != nil {
		return err
	}

	return nil
}

func (rd *BiSignedRequestData) validateBy(pk string, signature string) (error) {
	validator, err := NewECSDAValidator(pk)
	if err != nil {
		return err
	}

	return validator.ValidateString(rd.Message, signature)
}

func (rd *BiSignedRequestData) Validate(message *BiSignedMessage) bool {
	if message == nil {
		return false
	}

	return rd.FromPublicKey == message.FromPublicKey &&
		rd.ToPublicKey == message.ToPublicKey
}
