package websocketHelper

import (
	"encoding/json"

	"app/websocket"
	"svcledger/websocketData"
)

type ErrorResponse struct {
	payload	[]byte
}

func NewErrorResponse(corrToken string, data string, errorText string) *ErrorResponse {
	payload, _ := json.Marshal(&websocketData.ErrorResponseData{
		CorrToken: corrToken,
		Data: data,
		Error: errorText,
	})

	return &ErrorResponse{
		payload: payload,
	}
}

func (e* ErrorResponse) SendError(ctx websocket.HandlerContext, text string) {
	ctx.Send(&websocket.Message{
		OpCode: websocket.OpText,
		Payload: e.payload,
	})
}
