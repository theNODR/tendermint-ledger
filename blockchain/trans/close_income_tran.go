package trans

type CloseIncomeTran struct {
	PeerPublicKey	string	`json:"k"`
	From			string	`json:"f"`
	To				string	`json:"t"`
}
