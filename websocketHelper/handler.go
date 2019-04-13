package websocketHelper

import (
	"fmt"
	"encoding/json"
	"time"

	"app/websocket"
	"common"
	"svcledger/logger"
	"svcledger/netHelpers"
	"svcledger/websocketData"
)

type handlerTimers struct {
	BeforeCreate	int64	`json:"beforeCreate"`
	Created			int64	`json:"created"`
	Handled			int64	`json:"handled"`
}

func Handler(createRequest netHelpers.CreateRequester, ctx websocket.HandlerContext, message *websocketData.IncomingMessage) {
	exData := &handlerTimers{
		BeforeCreate: time.Now().UnixNano(),
	}
	req, err := createRequest(
		NewBaseRequest(message),
	)

	exData.Created = time.Now().UnixNano()
	if err != nil {
		sendError(ctx, message, exData, err)
		return
	}
	data, err := req.Handle()

	if err != nil {
		sendError(ctx, message, exData, err)
		return
	}
	exData.Handled = time.Now().UnixNano()
	resp := NewResponse(
		ctx,
		message.Cmd,
		message.CorrelationToken,
		data,
		exData,
	)
	err = resp.Send()

	if err != nil {
		sendError(ctx, message, exData, err)
		return
	}
}

func sendError(ctx websocket.HandlerContext, message *websocketData.IncomingMessage, exData *handlerTimers, err error) {
	payload := ""

	if message.Payload != nil {
		payload = string([]byte(*message.Payload)[:])
	}

	exDataBytes, _ := json.Marshal(exData)
	exDataStr := ""
	if exDataBytes != nil {
		exDataStr = string(exDataBytes[:])
	}
	errorResponse := NewErrorResponse(message.CorrelationToken, exDataStr, err.Error())
	errorLogText := fmt.Sprintf(
		"[WS ledger] error: %v %v %v %v",
		ctx.ConnectionId(),
		err.Error(),
		message.CorrelationToken,
		payload,
	)
	common.Log.StatError(
		logger.ErrorLedgerHandler,
		common.Printf(errorLogText, ),
	)
	errorResponse.SendError(ctx, errorLogText)
}
