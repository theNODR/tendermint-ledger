package store

import (
	"github.com/pkg/errors"
)

type Price interface {
	Amount() TransactionAmountType
	QuantumPower() QuantumPowerType
	Cmp(Price) (int8, error)
	ToNewQuantumPower(QuantumPowerType) (Price, error)
	ToChannel() *ChannelPrice
}

type price struct {
	amount			TransactionAmountType
	quantumPower	QuantumPowerType
}

type ChannelPrice struct {
	Amount			string				`json:"amount"`
	QuantumPower	QuantumPowerType	`json:"quantumPower"`
}

func NewPrice(
	amount TransactionAmountType,
	quantumPower QuantumPowerType,
) Price {
	return &price{amount, quantumPower }
}

func (p *price) Amount() TransactionAmountType {
	return p.amount
}

func (p *price) QuantumPower() QuantumPowerType {
	return p.quantumPower
}

// -1 if this less then another
// 0 if this equal another
// 1 if this greater than another
func (p *price) Cmp(another Price) (int8, error) {
	if another == nil {
		return 0, errors.New("another is nil")
	}

	var first Price
	var second Price
	var err error = nil

	if another.QuantumPower() > p.quantumPower {
		first, err = p.ToNewQuantumPower(another.QuantumPower())
		second = another
	} else if another.QuantumPower() < p.quantumPower {
		first = p
		second, err = another.ToNewQuantumPower(p.quantumPower)
	} else {
		first = p
		second = another
	}

	if err != nil {
		return 0, err
	}

	if first.Amount() > second.Amount() {
		return 1, nil
	} else if first.Amount() < second.Amount() {
		return -1, nil
	} else {
		return 0, nil
	}
}

func (p *price) ToNewQuantumPower(newQuantumPower QuantumPowerType) (Price, error) {
	if p.quantumPower > newQuantumPower {
		return nil, errors.New("new quantum power should be greater or equal than current quantum power")
	}

	diffSize := uint64(1 << (newQuantumPower - p.quantumPower))
	return NewPrice(
		p.amount * TransactionAmountType(diffSize),
		newQuantumPower,
	), nil
}

func (p *price) ToChannel() *ChannelPrice {
	return &ChannelPrice{
		Amount: p.amount.ToString(),
		QuantumPower: p.quantumPower,
	}
}
