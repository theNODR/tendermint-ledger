package store

type TransactionStatus = int8

const (
	NoneTransactionStatus    TransactionStatus = -1
	InvalidTransactionStatus TransactionStatus = 0
	ValidTransactionStatus   TransactionStatus = 1
)
