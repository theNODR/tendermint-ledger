package store

type QueryItem struct {
	PageSize	int
	StartId		string

	Order		OrderType
	Fields		*GetTransactionsFields
}