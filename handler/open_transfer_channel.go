package handler

import (
	"strconv"

	"svcledger/helpers"
	"svcledger/netHelpers"
	"svcledger/store"
)

type openTransferChannelRequest struct {
	*netHelpers.BaseRequest
	ledger				*store.Ledger
	message				*store.OpenContractMessage
	trackerKeyPair		helpers.KeyPair
}

type openTransferChannelResponseData struct {
	PublicKey	string	`json:"pk"`
	TimeLock	int64	`json:"timelock"`
	LifeTime	int64	`json:"lifetime"`
}

func newOpenTransferChannelRequest(
	baseRequest *netHelpers.BaseRequest,
	keyPair helpers.KeyPair,
	ledger *store.Ledger,
	_ *store.Queries,
) (netHelpers.Requester, error) {
	data, err := helpers.GetBiSignedPayloadData(baseRequest.Payload)

	if err != nil {
		return nil, err
	}

	message, err := store.NewOpenContractMessage(data)

	if err != nil {
		return nil, err
	}

	return &openTransferChannelRequest{
		BaseRequest: baseRequest,
		ledger: ledger,
		message: message,
		trackerKeyPair: keyPair,
	}, nil
}

func (req *openTransferChannelRequest) Handle() (interface{}, error) {
	priceAmount, err := strconv.ParseUint(req.message.PriceAmount, 10, 64)
	if err != nil {
		return nil, err
	}

	state, err := req.ledger.OpenTransferChannel(
		req.message.ChannelId,
		req.message.FromAddress,
		req.message.FromPublicKey,
		req.message.ToAddress,
		req.message.ToPublicKey,
		store.TransactionAmountType(priceAmount),
		store.QuantumPowerType(req.message.PriceQuantumPower),
		req.message.PlannedQuantumCount,
	)

	if err != nil {
		return nil, err
	}

	trackerPublicKey, err := req.trackerKeyPair.PublicKey()
	if err != nil {
		return nil, err
	}

	data := &openTransferChannelResponseData{
		PublicKey: trackerPublicKey,
		TimeLock: state.TimeLock,
		LifeTime: state.LifeTime,
	}
	respData, err := helpers.NewResponseDataInterface(data, req.trackerKeyPair)
	if err != nil {
		return nil, err
	}

	return respData, nil
}
