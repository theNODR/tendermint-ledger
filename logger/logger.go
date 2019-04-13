package logger

const (
	ErrorLedgerWebSocketOpen = "ledger: web socket open"
	ErrorLedgerWebSocket = "ledger: web socket"
	ErrorLedgerInvalidMessage = "ledger: invalid message"
	ErrorLedgerInvalidCommandMessage = "ledger: invalid command message"
	ErrorLedgerHandler = "ledger: handler"
	ErrorLedgerChainSyncLoad = "ledger: chain sync: load"
	ErrorLedgerChainSyncParse = "ledger: chain sync: parse"
	ErrorLedgerParseTranRow = "ledger: parse transaction row"
	ErrorLedgerFinTranSerialize = "ledger: fin tran state data serialize"
	ErrorLedgerFinTranWrite = "ledger: fin tran write"
	ErrorLedgerSyncConnectFailed = "ledger: sync: connection"
	ErrorLedgerSyncTrans = "ledger: sync: parse tran"
)

const (
	EventLedgerWebSocketPing = "web socket ping"
	EventLedgerWebSocketPong = "web socket pong"
	EventLedgerWebSocketOpen = "web socket open"
	EventLedgerWebSocketClose = "web socket close"
	EventLedgerWebSocketMessage = "web socket message"
	EventLedgerAutoCloseBegin = "auto close: begin"
	EventLedgerAutoCloseEnd = "auto close: end"
	EventLedgerFinTranEmpty = "fin tran empty"
	EventLedgerFinTranSuccess = "fin tran success"
	EventLedgerSyncBegin = "sync: begin"
	EventLedgerSyncEnd = "sync: end"
)
