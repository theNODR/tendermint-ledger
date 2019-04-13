package store

import (
	"encoding/json"

	"github.com/pkg/errors"

	"svcledger/helpers"
)

type CloseContractMessage struct {
	*ContractMessage
	Amount			string	`json:"tokens"`
	QuantumCount	uint64	`json:"quantumCount"`
	Volume			uint64	`json:"volume"`
}

func NewCloseContractMessage(data *helpers.BiSignedRequestData) (*CloseContractMessage, error) {
	var message CloseContractMessage

	err := json.Unmarshal([]byte(data.Message), &message)

	if err != nil {
		return nil, err
	}

	if !data.Validate(message.BiSignedMessage) {
		return nil, errors.New("invalid close tran message content: public key")
	}

	return &message, nil
}

func (m *CloseContractMessage) IsEqual(another *CloseContractMessage) bool {
	if another == nil {
		return false
	}

	return m.ContractMessage.IsEqual(another.ContractMessage) &&
		m.Amount == another.Amount &&
		m.QuantumCount == another.QuantumCount &&
		m.Volume == another.Volume
}
