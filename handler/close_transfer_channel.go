package handler

import (
	"strconv"

	"svcledger/helpers"
	"svcledger/netHelpers"
	"svcledger/store"
)

type closeTransferChannelRequest struct {
	*netHelpers.BaseRequest
	ledger				*store.Ledger
	message				*store.CloseContractMessage
	trackerKeyPair		helpers.KeyPair
}

type closeTransferChannelResponseData struct {
	PublicKey	string	`json:"pk"`
}

func newCloseTransferChannelRequest(
	baseRequest *netHelpers.BaseRequest,
	keyPair helpers.KeyPair,
	ledger *store.Ledger,
	_ *store.Queries,
) (netHelpers.Requester, error) {
	data, err := helpers.GetBiSignedPayloadData(baseRequest.Payload)

	if err != nil {
		return nil, err
	}

	message, err := store.NewCloseContractMessage(data)

	if err != nil {
		return nil, err
	}

	return &closeTransferChannelRequest{
		BaseRequest: baseRequest,
		message: message,
		trackerKeyPair: keyPair,
		ledger: ledger,
	}, nil
}

func (req *closeTransferChannelRequest) Handle() (interface{}, error) {
	priceAmount, err := strconv.ParseUint(req.message.PriceAmount, 10, 64)
	if err != nil {
		return nil, err
	}

	amount, err := strconv.ParseUint(req.message.Amount, 10, 64)
	if err != nil {
		return nil, err
	}

	_, err = req.ledger.CloseTransferChannel(
		req.message.ChannelId,
		req.message.FromAddress,
		req.message.FromPublicKey,
		req.message.ToAddress,
		req.message.ToPublicKey,
		store.TransactionAmountType(priceAmount),
		store.QuantumPowerType(req.message.PriceQuantumPower),
		req.message.PlannedQuantumCount,
		store.TransactionAmountType(amount),
		req.message.QuantumCount,
		req.message.Volume,
	)
	if err != nil {
		return nil, err
	}

	trackerPublicKey, err := req.trackerKeyPair.PublicKey()
	if err != nil {
		return nil, err
	}

	data := &closeTransferChannelResponseData{
		PublicKey: trackerPublicKey,
	}
	respData, err := helpers.NewResponseDataInterface(data, req.trackerKeyPair)
	if err != nil {
		return nil, err
	}

	return respData, nil
}
