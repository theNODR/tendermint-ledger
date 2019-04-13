package websocketHelper

import (
	"svcledger/netHelpers"
	"svcledger/websocketData"
)

func NewBaseRequest(message *websocketData.IncomingMessage) *netHelpers.BaseRequest {
	return &netHelpers.BaseRequest{
		Payload: message.Payload,
	}
}
