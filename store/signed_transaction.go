package store

import (
	"encoding/base64"

	"github.com/pkg/errors"

	"common"
	"svcledger/helpers"
	"svcledger/blockchain"
)

type SignedTransaction struct {
	*Transaction
	Sign	string
}

func NewSignedTransaction(
	body *TransactionBody,
	keyPair helpers.KeyPair,
	originalId string,
	originalPk string,
) (*SignedTransaction, error) {
	newId := blockchain.CreateId(originalPk, body.TimeStamp)
	if newId != originalId {
		return nil, errors.New("invalid sended id")
	}

	tran := &Transaction{
		TransactionBody: body,
		Id: newId,
		PublicKey: originalPk,
	}

	signBytes, err := keyPair.SignObject(tran)

	if err != nil {
		return nil, err
	}

	return &SignedTransaction{
		Transaction: tran,
		Sign: string(signBytes[:]),
	}, nil
}

func (t *SignedTransaction) Serialize() []string {
	return common.ToString(
		//server_date
		t.TimeStamp,
		//timelock
		t.TimeLock,
		//tx_id
		t.Id,
		// parent_tx_id
		t.ParentId,
		// version
		t.Version,
		// tx_type
		transactionTypeNames[t.Type],
		// from
		t.From,
		// to
		t.To,
		// amount
		t.Amount,
		// volume
		t.Volume,
		// quantum_count
		t.QuantumCount,
		// contract_id
		t.ContractId,
		// price_amount
		t.PriceAmount,
		// price_quantum_power
		t.PriceQuantumPower,
		// planned_quantum_count
		t.PlannedQuantumCount,
		// data
		t.Data,
		// status
		//t.Status,
		// pk
		base64.StdEncoding.EncodeToString([]byte(t.PublicKey)),
		// sign
		base64.StdEncoding.EncodeToString([]byte(t.Sign)),
	)
}

func (t *SignedTransaction) validate() error {
	validator, err := helpers.NewECSDAValidator(t.PublicKey)
	if err != nil {
		return err
	}

	return validator.ValidateObject(t.Transaction, t.Sign)
}

func (t *SignedTransaction) serializeToJson() *JsonTransaction {
	return &JsonTransaction{
		Amount: uint64(t.Amount),
		ContractId: t.ContractId,
		From: t.From,
		Id: t.Id,
		ParentId: t.ParentId,
		PlannedQuantumCount: t.PlannedQuantumCount,
		PriceAmount: uint64(t.PriceAmount),
		PriceQuantumPower: uint8(t.PriceQuantumPower),
		PublicKey: t.PublicKey,
		QuantumCount: t.QuantumCount,
		//Status: int16(t.Status),
		Sign: t.Sign,
		TimeLock: t.TimeLock,
		TimeStamp: t.TimeStamp,
		To: t.To,
		Type: int8(t.Type),
		Version: t.Version,
		Volume: t.Volume,
	}
}
