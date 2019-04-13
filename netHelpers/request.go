package netHelpers

import (
	"encoding/json"

	"svcledger/helpers"
	"svcledger/store"
)

type Requester interface {
	Handle() (interface{}, error)
}

type BaseRequest struct {
	Payload *json.RawMessage
}

type CreateRequester = func(request *BaseRequest) (Requester, error)

type NewRequest = func(
	request *BaseRequest,
	keyPair helpers.KeyPair,
	ledger *store.Ledger,
) (Requester, error)
