package handler

import (
	"encoding/json"

	"github.com/pkg/errors"

	"svcledger/helpers"
	"svcledger/netHelpers"
	"svcledger/store"
)

type getNextTransactions struct {
	*netHelpers.BaseRequest

	ledger		*store.Ledger
	queries		*store.Queries
	trackerKeyPair		helpers.KeyPair

	cursorId	string
	page		uint
}

type getNextTransactionsRequestData struct {
	PublicKey	string	`json:"pk"`

	Page		uint	`json:"page"`
	CursorId	string	`json:"cursorId"`
}

type getNextTransactionsResponseData struct {
	Items	store.JsonTransactions	`json:"items"`
}

func newGetNextTransactionsRequest(
	request *netHelpers.BaseRequest,
	keyPair helpers.KeyPair,
	ledger *store.Ledger,
	queries *store.Queries,
) (netHelpers.Requester, error) {
	data, err := helpers.GetSignedPayloadData(request.Payload)

	if err != nil {
		return nil, err
	}

	var obj getNextTransactionsRequestData
	err = json.Unmarshal([]byte(data.Message), &obj)
	if err != nil {
		return nil, err
	}
	if obj.PublicKey != data.PublicKey {
		return nil, errors.New("invalid message content")
	}

	return &getNextTransactions{
		BaseRequest: request,

		ledger: ledger,
		queries: queries,
		trackerKeyPair: keyPair,

		cursorId: obj.CursorId,
		page: obj.Page,
	}, nil
}

func (req *getNextTransactions) Handle() (interface{}, error) {
	rawParams, err := req.queries.Get(req.cursorId)
	if err != nil {
		return nil, err
	}
	params := rawParams.(*store.QueryItem)

	trans, err := req.ledger.GetNextTransactions(
		params.StartId,
		params.PageSize,
		req.page,
		params.Fields,
		params.Order,
	)
	if err != nil {
		return nil, err
	}

	data := &getNextTransactionsResponseData{
		Items: store.ToJsonTransactions(trans.Data),
	}
	respData, err := helpers.NewResponseDataInterface(data, req.trackerKeyPair)
	if err != nil {
		return nil, err
	}

	return respData, nil
}
