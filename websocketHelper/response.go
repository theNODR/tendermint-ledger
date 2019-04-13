package websocketHelper

import (
	"encoding/json"

	"app/websocket"
	"svcledger/netHelpers"
	"svcledger/websocketData"
)

type Response struct {
	cmd					websocketData.CmdType
	correlationToken	string
	ctx					websocket.HandlerContext
	data 				interface{}
	exData				interface{}
}

func NewResponse(
	ctx websocket.HandlerContext,
	cmd websocketData.CmdType,
	correlationToken string,
	data interface{},
	exData interface{},
) netHelpers.Responser {
	return &Response{
		cmd: cmd,
		correlationToken: correlationToken,
		ctx: ctx,
		data: data,
		exData: exData,
	}
}

func (resp* Response) Send() error {
	data, err := json.Marshal(resp.data)

	if err != nil {
		return err
	}

	exData, err := json.Marshal(resp.exData)

	if err != nil {
		return err
	}

	response, err := json.Marshal(&websocketData.ResponseData{
		Cmd: resp.cmd,
		CorrelationToken: resp.correlationToken,
		Data: data,
		ExData: exData,
	})

	if err != nil {
		return err
	}

	resp.ctx.Send(&websocket.Message{
		OpCode: websocket.OpText,
		Payload: response,
	})

	return nil
}
