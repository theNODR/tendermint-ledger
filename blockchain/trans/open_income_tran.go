package trans

type OpenIncomeTran struct {
	PeerPublicKey		string	`json:"p"`
	From				string	`json:"f"`
	PriceAmount			uint64	`json:"a"`
	PriceQuantumPower	uint8	`json:"q"`
	LifeTime			int64	`json:"l"`
}
