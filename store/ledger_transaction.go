package store

type TransactionStatus = int8
const (
	NoneTransactionStatus		TransactionStatus = -1
	InvalidTransactionStatus	TransactionStatus = 0
	ValidTransactionStatus		TransactionStatus = 1
)

type LedgerTransaction struct {
	*SignedTransaction
	ServerTimestamp	int64
	Status			TransactionStatus
}

func (t *LedgerTransaction) IsValid() bool {
	return t.Status == ValidTransactionStatus
}

type LedgerTransactions = []*LedgerTransaction
