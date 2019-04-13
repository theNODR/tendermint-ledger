package blockchain

import (
	"encoding/json"

	"common"
	"svcledger/blockchain/trans"
)

type TransactionType uint8

const (
	InitTransactionType					TransactionType = 1
	AutoCloseTransactionType			TransactionType = 2
	FinancialTransactionType			TransactionType = 3
	CloseIncomeChannelTransactionType	TransactionType = 4
	CloseTransferChannelTransactionType	TransactionType = 5
	CloseSpendChannelTransactionType	TransactionType = 6
	OpenIncomeChannelTransactionType	TransactionType = 7
	OpenSpendChannelTransactionType		TransactionType = 8
	OpenTransferChannelTransactionType	TransactionType = 9
)

type Tran struct {
	Id			string			`json:"i"`
	Data		interface{}		`json:"d"`
	PublicKey	string			`json:"p"`
	Type		TransactionType	`json:"y"`
	TimeStamp	int64			`json:"t"`
}

type ReadTran struct {
	Id			string			`json:"i"`
	Data		json.RawMessage `json:"d"`
	PublicKey	string			`json:"p"`
	Type		TransactionType	`json:"y"`
	TimeStamp	int64			`json:"t"`
}

func newTran(data interface{}, pk string, t TransactionType) *Tran {
	timestamp := common.GetNowUnixMs()
	return &Tran{
		Id: CreateId(pk, timestamp),
		Data: data,
		PublicKey: pk,
		Type: t,
		TimeStamp: timestamp,
	}
}

func (t Tran) ToString() string {
	bytes, err := json.Marshal(&t)
	if err != nil {
		return ""
	}
	return string(bytes[:])
}

func NewInitTransactionType(data *trans.InitTran, pk string) *Tran {
	return newTran(data, pk, InitTransactionType)
}

func NewAutoCloseTran(data *trans.AutoCloseTran, pk string) *Tran {
	return newTran(data, pk, AutoCloseTransactionType)
}

func NewFinancialTran(data *trans.FinancialTran, pk string) *Tran {
	return newTran(data, pk, FinancialTransactionType)
}

func NewOpenIncomeTran(data *trans.OpenIncomeTran, pk string) *Tran {
	return newTran(data, pk, OpenIncomeChannelTransactionType)
}

func NewOpenSpendTran(data *trans.OpenSpendTran, pk string) *Tran {
	return newTran(data, pk, OpenSpendChannelTransactionType)
}

func NewOpenTransferTran(data *trans.OpenTransferTran, pk string) *Tran {
	return newTran(data, pk, OpenTransferChannelTransactionType)
}

func NewCloseTransferTran(data *trans.CloseTransferTran, pk string) *Tran {
	return newTran(data, pk, CloseTransferChannelTransactionType)
}

func NewCloseSpendTran(data *trans.CloseSpendTran, pk string) *Tran {
	return newTran(data, pk, CloseSpendChannelTransactionType)
}

func NewCloseIncomeTran(data *trans.CloseIncomeTran, pk string) *Tran {
	return newTran(data, pk, CloseIncomeChannelTransactionType)
}
