package handler

import (
	"encoding/json"

	"github.com/pkg/errors"

	"svcledger/helpers"
	"svcledger/netHelpers"
	"svcledger/store"
)

type updateTransactionsCursorRequest struct {
	*netHelpers.BaseRequest
	cursorId		string
	queries			*store.Queries
	trackerKeyPair	helpers.KeyPair
}

type updateTransactionsCursorRequestData struct {
	PublicKey	string	`json:"pk"`
	CursorId	string	`json:"cursorId"`
}

type updateTransactionsCursorResponseData struct {
	PublicKey	string	`json:"pk"`
}

func newUpdateTransactionsCursorRequest(
	baseRequest *netHelpers.BaseRequest,
	keyPair helpers.KeyPair,
	_ *store.Ledger,
	queries *store.Queries,
) (netHelpers.Requester, error) {
	data, err := helpers.GetSignedPayloadData(baseRequest.Payload)
	if err != nil {
		return nil, err
	}

	var obj updateTransactionsCursorRequestData
	err = json.Unmarshal([]byte(data.Message), &obj)
	if err != nil {
		return nil, err
	}
	if obj.PublicKey != data.PublicKey {
		return nil, errors.New("invalid message content")
	}

	return &updateTransactionsCursorRequest{
		BaseRequest: baseRequest,
		cursorId: obj.CursorId,
		queries: queries,
		trackerKeyPair: keyPair,
	}, nil
}

func (req *updateTransactionsCursorRequest) Handle() (interface{}, error) {
	_, err := req.queries.Get(req.cursorId)

	if err != nil {
		return nil, err
	}

	trackerPublicKey, err := req.trackerKeyPair.PublicKey()
	if err != nil {
		return nil, err
	}

	data := &updateTransactionsCursorResponseData{
		PublicKey: trackerPublicKey,
	}
	respData, err := helpers.NewResponseDataInterface(data, req.trackerKeyPair)
	if err != nil {
		return nil, err
	}

	return respData, nil
}
