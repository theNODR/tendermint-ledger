package svcledger

import (
	"encoding/json"
	"fmt"

	"app/websocket"
	"common"
	"svcledger/handler"
	"svcledger/helpers"
	"svcledger/logger"
	"svcledger/store"
	"svcledger/websocketData"
	"svcledger/websocketHelper"
)

func handlerFunc(
	keyPair helpers.KeyPair,
	ledger *store.Ledger,
) (websocket.HandlerFunc, error) {
	handlers, err := handler.NewHandlers()

	if err != nil {
		return nil, err
	}

	return func(ctx websocket.HandlerContext) {
		switch ctx.EventType() {
		case websocket.SocketEventPing:
			common.Log.Event(
				logger.EventLedgerWebSocketPing,
				common.Printf("[WS ledger] ping: %v", ctx.ConnectionId()),
			)
			break

		case websocket.SocketEventPong:
			common.Log.Event(
				logger.EventLedgerWebSocketPong,
				common.Printf("[WS ledger] pong: %v", ctx.ConnectionId()),
			)
			break

		case websocket.SocketEventOpen:
			common.Log.Event(
				logger.EventLedgerWebSocketOpen,
				common.Printf("[WS ledger] open: %v", ctx.ConnectionId()),
			)
			break

		case websocket.SocketEventClose:
			common.Log.Event(
				logger.EventLedgerWebSocketClose,
				common.Printf("[WS ledger] close: %v", ctx.ConnectionId()),
			)
			break

		case websocket.SocketEventError:
			if !ctx.ExistConnection() {
				common.Log.StatError(
					logger.ErrorLedgerWebSocketOpen,
					common.Printf("[WS ledger] error: on open connection: %v", ctx.Error()),
				)
			} else {
				common.Log.StatError(
					logger.ErrorLedgerWebSocket,
					common.Printf("[WS ledger] error: %v: %v", ctx.ConnectionId(), ctx.Error()),
				)
			}
			break

		case websocket.SocketEventMessage:
			common.Log.Event(
				logger.EventLedgerWebSocketMessage,
				common.Printf("[WS ledger] message: %v", ctx.ConnectionId()),
			)
			incomingMessage := &websocketData.IncomingMessage{}
			err := json.Unmarshal(ctx.Message().Payload, incomingMessage)

			if err != nil {
				sendInvalidPayloadError(ctx, err)
				return
			}

			handlerFunc, err := handlers.Execute(incomingMessage.Cmd)

			if err != nil {
				sendInvalidCmdError(ctx, incomingMessage)
				return
			}

			wsHandlerFunc := handlerFunc(
				keyPair,
				ledger,
			)

			websocketHelper.Handler(
				wsHandlerFunc,
				ctx,
				incomingMessage,
			)

			break

		default:
			return
		}
	}, nil
}

func sendInvalidPayloadError(ctx websocket.HandlerContext, err error) {
	errorText := fmt.Sprintf("invalid message: %v", err.Error())
	errorResponse := websocketHelper.NewErrorResponse("", string(ctx.Message().Payload[:]), errorText)

	errorLogText := fmt.Sprintf(
		"[WS ledger] error: invalid message: %v %v %v",
		ctx.ConnectionId(),
		ctx.MessagePayloadString(),
		err.Error(),
	)
	common.Log.StatError(
		logger.ErrorLedgerInvalidMessage,
		common.Printf(errorLogText),
	)
	errorResponse.SendError(ctx, errorLogText)
}

func sendInvalidCmdError(ctx websocket.HandlerContext, message *websocketData.IncomingMessage) {
	errorText := fmt.Sprintf("invalid message command: %v", message.Cmd)
	errorResponse := websocketHelper.NewErrorResponse(message.CorrelationToken, message.CorrelationToken, errorText)
	errorLogText := fmt.Sprintf(
		"[WS ledger] error: invalid message command: %v %v %v",
		ctx.ConnectionId(),
		message.CorrelationToken,
		message.Cmd,
	)
	common.Log.StatError(
		logger.ErrorLedgerInvalidCommandMessage,
		common.Printf(errorLogText),
	)
	errorResponse.SendError(ctx, errorLogText)
}
