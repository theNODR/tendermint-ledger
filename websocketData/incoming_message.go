package websocketData

import "encoding/json"

type CmdType int

const (
	HandshakeCmd CmdType = 1 << iota
	OpenIncomeChannelCmd CmdType = 1 << iota
	OpenSpendChannelCmd CmdType = 1 << iota
	OpenTransferChannelCmd CmdType = 1 << iota
	CloseTransferChannelCmd CmdType = 1 << iota
	CloseIncomeChannelCmd CmdType = 1 << iota
	CloseSpendChannelCmd CmdType = 1 << iota
	GetIncomeAddressStateCmd CmdType = 1 << iota
	GetSpendAddressStateCmd CmdType = 1 << iota
	GetTransferChannelStateCmd CmdType = 1 << iota
	GetFirstTransactionsCmd CmdType = 1 << iota
	GetNextTransactionsCmd CmdType = 1 << iota
	UpdateTransactionsCursorCmd CmdType = 1 << iota
	GetCommonAddressStateCmd CmdType = 1 << iota
)

type IncomingMessage struct {
	Cmd					CmdType				`json:"cmd"`
	CorrelationToken	string				`json:"corrToken"`
	Payload				*json.RawMessage	`json:"payload"`
}
