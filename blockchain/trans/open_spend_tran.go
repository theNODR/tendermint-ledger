package trans

type OpenSpendTran struct {
	PeerPublicKey		string	`json:"p"`
	From				string	`json:"f"`
	MaxAmount			uint64	`json:"m"`
	PriceAmount			uint64	`json:"a"`
	PriceQuantumPower	uint8	`json:"q"`
	LifeTime			int64	`json:"l"`
}
