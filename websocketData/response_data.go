package websocketData

import "encoding/json"

type ResponseData struct {
	Cmd					CmdType			`json:"cmd"`
	CorrelationToken	string			`json:"corrToken"`
	Data				json.RawMessage	`json:"data"`
	ExData				json.RawMessage	`json:"exData"`
}
