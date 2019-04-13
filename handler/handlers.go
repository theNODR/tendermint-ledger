package handler

import (
	"github.com/pkg/errors"

	"svcledger/helpers"
	"svcledger/netHelpers"
	"svcledger/store"
	"svcledger/websocketData"
)

type WSHandlerFunc = func(
	keyPair helpers.KeyPair,
	ledger *store.Ledger,
) netHelpers.CreateRequester

type Handlers struct {
	items map[websocketData.CmdType]WSHandlerFunc
}

func NewHandlers() (*Handlers, error) {
	h := &Handlers{
		items: make(map[websocketData.CmdType]WSHandlerFunc),
	}

	var err error

	err = h.register(websocketData.CloseIncomeChannelCmd, newCloseIncomeChannelRequest)
	if err != nil {
		return nil, err
	}

	err = h.register(websocketData.CloseSpendChannelCmd, newCloseSpendChannelRequest)
	if err != nil {
		return nil, err
	}

	err = h.register(websocketData.CloseTransferChannelCmd, newCloseTransferChannelRequest)
	if err != nil {
		return nil, err
	}

	err = h.register(websocketData.GetIncomeAddressStateCmd, newGetIncomeAddressStateRequest)
	if err != nil {
		return nil, err
	}

	err = h.register(websocketData.GetSpendAddressStateCmd, newGetSpendAddressStateRequest)
	if err != nil {
		return nil, err
	}

	err = h.register(websocketData.GetTransferChannelStateCmd, newTransferChannelStateRequest)
	if err != nil {
		return nil, err
	}

	err = h.register(websocketData.HandshakeCmd, newHandshakeRequest)
	if err != nil {
		return nil, err
	}

	err = h.register(websocketData.OpenIncomeChannelCmd, newOpenIncomeChannelRequest)
	if err != nil {
		return nil, err
	}

	err = h.register(websocketData.OpenSpendChannelCmd, newOpenSpendChannelRequest)
	if err != nil {
		return nil, err
	}

	err = h.register(websocketData.OpenTransferChannelCmd, newOpenTransferChannelRequest)
	if err != nil {
		return nil, err
	}

	err = h.register(websocketData.GetCommonAddressStateCmd, newGetCommonAddressState)
	if err != nil {
		return nil, err
	}

	return h, nil
}

func (h *Handlers) createHandler(newRequest netHelpers.NewRequest) WSHandlerFunc {
	return func(
		keyPair helpers.KeyPair,
		ledger *store.Ledger,
	) netHelpers.CreateRequester {
		return func(baseRequest *netHelpers.BaseRequest) (netHelpers.Requester, error) {
			return newRequest(
				baseRequest,
				keyPair,
				ledger)
		}
	}
}

func (h *Handlers) register(handlerType websocketData.CmdType, request netHelpers.NewRequest) error {
	if _, ok := h.items[handlerType]; ok {
		return errors.New("handler already exist")
	}

	h.items[handlerType] = h.createHandler(request)

	return nil
}

func (h *Handlers) Execute(handlerType websocketData.CmdType) (WSHandlerFunc, error) {
	item, ok := h.items[handlerType]

	if !ok {
		return nil, errors.New("handler not exist")
	}

	return item, nil
}
