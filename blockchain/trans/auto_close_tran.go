package trans

type AutoCloseTran struct {
	Timestamp     int64  `json:"time"`
	LedgerAddress string `json:"address"`
}
