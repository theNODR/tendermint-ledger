package store

type PageParams struct {
	Page		uint	`json:"page"`
	PageSize	int		`json:"pageSize"`
}

type OrderType = uint

const (
	AscOrder	OrderType = 1 << iota
	DescOrder	OrderType = 1 << iota
	NoneOrder	OrderType = 0
)

type GetTransactionsFieldType struct {
	Value	TransactionType		`json:"value"`
}

type GetTransactionsFieldString struct {
	Value	string				`json:"value"`
}

type GetTransactionsFieldStatus struct {
	Value	TransactionStatus	`json:"value"`
}

type GetTransactionsFields struct {
	Type		*GetTransactionsFieldType	`json:"tx_type"`
	ChannelId	*GetTransactionsFieldString	`json:"channel_id"`
	From		*GetTransactionsFieldString	`json:"from"`
	To			*GetTransactionsFieldString	`json:"to"`
	Status		*GetTransactionsFieldStatus	`json:"status"`
}

type GetFirstTransactionsRequestData struct {
	PublicKey		string					`json:"pk"`
	Limits			*PageParams				`json:"limits"`
	Order			OrderType				`json:"order"`
	Fields			*GetTransactionsFields	`json:"fields"`
}

type MetaTransactionsResponseData struct {
	CursorId	string	`json:"cursorId"`
	RowCount	uint	`json:"totalCount"`
	PageCount	uint	`json:"pageCount"`
	CurrentPage	uint	`json:"currentPage"`
	PerPage		uint	`json:"perPage"`
}

type GetFirstTransactionsResponseData struct {
	Meta	*MetaTransactionsResponseData	`json:"_meta"`
	Items	JsonTransactions				`json:"items"`
}
