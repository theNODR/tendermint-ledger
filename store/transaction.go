package store

type Transaction struct {
	*TransactionBody
	Id			string
	PublicKey	string
}
