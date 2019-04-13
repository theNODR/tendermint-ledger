package store

import (
	"svcledger/helpers"
)

type ContractMessage struct {
	*helpers.BiSignedMessage
	ChannelId			string	`json:"channelId"`
	FromAddress			string	`json:"fromAddress"`
	ToAddress			string	`json:"toAddress"`
	PriceAmount			string	`json:"priceTokens"`
	PriceQuantumPower	uint8	`json:"priceQuantumPower"`
	PlannedQuantumCount	uint64	`json:"plannedQuantumCount"`
}

func (t *ContractMessage) IsEqual(another* ContractMessage) bool {
	if another == nil {
		return false
	}

	return t.IsEqual(another) &&
		t.ChannelId == another.ChannelId &&
		t.PriceAmount == another.PriceAmount &&
		t.PriceQuantumPower == another.PriceQuantumPower &&
		t.PlannedQuantumCount == another.PlannedQuantumCount &&
		t.FromAddress == another.FromAddress &&
		t.ToAddress == another.ToAddress
}
