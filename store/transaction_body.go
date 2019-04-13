package store

import "strconv"

type TransactionType = int16
const (
	InitTransactionType 				TransactionType = 1 << iota
	OpenIncomeChannelTransactionType	TransactionType = 1 << iota
	OpenSpendChannelTransactionType   	TransactionType = 1 << iota
	OpenTransferChannelTransactionType	TransactionType = 1 << iota
	CloseTransferChannelTransactionType	TransactionType = 1 << iota
	CloseSpendChannelTransactionType	TransactionType = 1 << iota
	CloseIncomeChannelTransactionType	TransactionType = 1 << iota
	FinancialBlockTransactionType		TransactionType = 1 << iota
	AutoCloseTransactionType			TransactionType = 1 << iota
	UnknownTransactionType 				TransactionType = 0
)
var transactionTypeNames = map[TransactionType]string{
	InitTransactionType: "TX_INIT",
	OpenIncomeChannelTransactionType: "TX_OPEN_INCOME",
	OpenSpendChannelTransactionType: "TX_OPEN_SPEND",
	OpenTransferChannelTransactionType: "TX_OPEN_TRANSFER",
	CloseTransferChannelTransactionType: "TX_CLOSE_TRANSFER",
	CloseSpendChannelTransactionType: "TX_CLOSE_SPEND",
	CloseIncomeChannelTransactionType: "TX_CLOSE_INCOME",
	FinancialBlockTransactionType: "TX_FINANCIAL_BLOCK",
}
var transactionTypeCodes = map[string]TransactionType {
	"TX_INIT": InitTransactionType,
	"TX_OPEN_INCOME": OpenIncomeChannelTransactionType,
	"TX_OPEN_SPEND": OpenSpendChannelTransactionType,
	"TX_OPEN_TRANSFER": OpenTransferChannelTransactionType,
	"TX_CLOSE_TRANSFER": CloseTransferChannelTransactionType,
	"TX_CLOSE_SPEND": CloseSpendChannelTransactionType,
	"TX_CLOSE_INCOME": CloseIncomeChannelTransactionType,
	"TX_FINANCIAL_BLOCK": FinancialBlockTransactionType,
}

type QuantumPowerType uint8
func (p QuantumPowerType) Volume() uint64 {
	return 1 << p
}

type TransactionAmountType uint64
func (t TransactionAmountType) ToString() string {
	return strconv.FormatUint(uint64(t), 10)
}

type TransactionBody struct {
	TimeStamp			int64
	TimeLock			int64
	ParentId			string
	Version				string
	Type				TransactionType
	From				string
	To					string
	Amount				TransactionAmountType
	Volume				uint64
	QuantumCount		uint64
	ContractId			string
	PriceAmount			TransactionAmountType
	PriceQuantumPower	QuantumPowerType
	PlannedQuantumCount	uint64
	Data				interface{}
}

func (t *TransactionBody) PlannedVolume() uint64 {
	return t.PlannedQuantumCount * t.PriceQuantumPower.Volume()
}

func (t *TransactionBody) AmountedVolume() uint64 {
	return t.QuantumCount * t.PriceQuantumPower.Volume()
}

func (t *TransactionBody) PlannedAmount() TransactionAmountType {
	return t.PriceAmount * TransactionAmountType(t.PlannedVolume())
}
