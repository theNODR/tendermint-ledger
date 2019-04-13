package store

import (
	"encoding/json"

	"github.com/pkg/errors"

	"svcledger/helpers"
)

type OpenContractMessage struct {
	*ContractMessage
}

func NewOpenContractMessage(data *helpers.BiSignedRequestData) (*OpenContractMessage, error) {
	var message OpenContractMessage

	err := json.Unmarshal([]byte(data.Message), &message)

	if err != nil {
		return nil, err
	}

	if !data.Validate(message.BiSignedMessage) {
		return nil, errors.New("invalid open tran message content: public key")
	}

	return &message, nil
}

func (m *OpenContractMessage) IsEqual(another *OpenContractMessage) bool {
	if another == nil {
		return false
	}

	return m.ContractMessage.IsEqual(another.ContractMessage)
}
