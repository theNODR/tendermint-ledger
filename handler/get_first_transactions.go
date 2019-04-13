package handler

import (
	"encoding/json"

	"github.com/pkg/errors"

	"svcledger/helpers"
	"svcledger/netHelpers"
	"svcledger/store"
)

const (
	defaultPage		uint = 1
	defaultPageSize	int = 20
)

type getFirstTransactions struct {
	*netHelpers.BaseRequest

	ledger			*store.Ledger
	queries			*store.Queries
	trackerKeyPair	helpers.KeyPair

	fields			*store.GetTransactionsFields
	limits			*store.PageParams
	order			store.OrderType
}

func newGetFirstTransactionsRequest(
	baseRequest *netHelpers.BaseRequest,
	keyPair helpers.KeyPair,
	ledger *store.Ledger,
	queries *store.Queries,
) (netHelpers.Requester, error) {
	data, err := helpers.GetSignedPayloadData(baseRequest.Payload)

	if err != nil {
		return nil, err
	}

	var obj store.GetFirstTransactionsRequestData
	err = json.Unmarshal([]byte(data.Message), &obj)
	if err != nil {
		return nil, err
	}
	if obj.PublicKey != data.PublicKey {
		return nil, errors.New("invalid message content")
	}

	limits := obj.Limits

	if limits == nil {
		limits = &store.PageParams{
			Page: defaultPage,
			PageSize: defaultPageSize,
		}
	} else {
		if limits.Page == 0 {
			limits.Page = 1
		}

		if limits.PageSize == 0 {
			limits.PageSize = defaultPageSize
		}

		// Если число элементов на странице отрицательное, считаем что осуществляется возврат всех транзкаций на одной странице
		if limits.PageSize < 0 {
			limits.Page = 1
		}
	}

	order := obj.Order

	if order != store.AscOrder {
		order = store.DescOrder
	}

	return &getFirstTransactions{
		BaseRequest: baseRequest,

		ledger: ledger,
		queries: queries,
		trackerKeyPair: keyPair,

		fields: obj.Fields,
		limits: limits,
		order: order,
	}, nil
}

func (req *getFirstTransactions) Handle() (interface{}, error) {
	trans, err := req.ledger.GetFirstTransactions(
		req.limits.PageSize,
		req.limits.Page,
		req.fields,
		req.order,
	)

	if err != nil {
		return nil, err
	}

	startId := ""

	if trans.RowCount > 0 {
		startId = trans.Data[0].Id
	}

	cursorId, err := req.queries.Add(&store.QueryItem{
		PageSize: req.limits.PageSize,
		StartId: startId,
		Order: req.order,
		Fields: req.fields,
	})

	if err != nil {
		return nil, err
	}

	var pageCount uint
	var pageSize uint

	if req.limits.PageSize < 0 {
		pageCount = 1
		pageSize = trans.RowCount
	} else {
		pageSize = uint(req.limits.PageSize)
		pageCount = trans.RowCount / pageSize
		if trans.RowCount % pageSize != 0 {
			pageCount++
		}
	}


	data := &store.GetFirstTransactionsResponseData{
		Items: store.ToJsonTransactions(trans.Data),
		Meta: &store.MetaTransactionsResponseData{
			CursorId: cursorId,
			PerPage: pageSize,
			RowCount: trans.RowCount,
			PageCount: pageCount,
			CurrentPage: req.limits.Page,
		},
	}
	respData, err := helpers.NewResponseDataInterface(data, req.trackerKeyPair)
	if err != nil {
		return nil, err
	}

	return respData, nil
}
