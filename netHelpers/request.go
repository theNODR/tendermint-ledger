package netHelpers

import (
	"encoding/json"

	"svcledger/store"
	"svcledger/helpers"
)

type Requester interface {
	Handle() (interface{}, error)
}

type BaseRequest struct {
	Payload				*json.RawMessage
}

type CreateRequester = func(request *BaseRequest) (Requester, error)

type NewRequest = func(
	request *BaseRequest,
	keyPair helpers.KeyPair,
	ledger *store.Ledger,
	queries *store.Queries,
) (Requester, error)