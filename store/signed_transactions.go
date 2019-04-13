package store

type SignedTransactions = []*SignedTransaction
type JsonTransactions = []*JsonTransaction

func ToJsonTransactions(t []*LedgerTransaction) JsonTransactions {
	size := len(t)
	result := make(JsonTransactions, size, size)

	for i, _ := range t {
		result[i] = t[i].serializeToJson()
	}

	return result
}
