package trans

type OpenTransferTran struct {
	ChannelId			string	`json:"c"`
	FromAddress			string	`json:"f"`
	FromPk				string	`json:"k"`
	ToAddress			string	`json:"t"`
	ToPk				string	`json:"l"`
	PriceAmount			uint64	`json:"a"`
	PriceQuantumPower	uint8	`json:"p"`
	PlannedQuantumCount uint64	`json:"q"`
	LifeTime			int64	`json:"i"`
}
