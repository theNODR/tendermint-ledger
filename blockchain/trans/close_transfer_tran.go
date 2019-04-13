package trans

type CloseTransferTran struct {
	ChannelId			string	`json:"c"`
	FromAddress 		string	`json:"f"`
	FromPk				string	`json:"k"`
	ToAddress			string	`json:"t"`
	ToPk				string	`json:"l"`
	PriceAmount			uint64	`json:"z"`
	PriceQuantumPower	uint8	`json:"w"`
	PlannedQuantumCount uint64	`json:"o"`
	Amount				uint64	`json:"a"`
	QuantumCount		uint64	`json:"q"`
	Volume				uint64	`json:"v"`
}
