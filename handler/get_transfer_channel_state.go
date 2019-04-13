package handler

import (
	"encoding/json"
	"github.com/pkg/errors"

	"svcledger/helpers"
	"svcledger/netHelpers"
	"svcledger/store"
	ntypes "svcnodr/types"
)

type getTransferChannelState struct {
	*netHelpers.BaseRequest
	channelId      string
	ledger         *store.Ledger
	trackerKeyPair helpers.KeyPair
}

type getTransferChannelStateRequestData struct {
	ChannelId string `json:"channelId"`
	PublicKey string `json:"pk"`
}

type getTransferChannelStateResponseData struct {
	ChannelId string                       `json:"channelId"`
	State     ntypes.TransferChannelStatus `json:"state"`
	PublicKey string                       `json:"pk"`
}

func newTransferChannelStateRequest(
	baseRequest *netHelpers.BaseRequest,
	keyPair helpers.KeyPair,
	ledger *store.Ledger,
) (netHelpers.Requester, error) {
	data, err := helpers.GetSignedPayloadData(baseRequest.Payload)

	if err != nil {
		return nil, err
	}

	var obj getTransferChannelStateRequestData
	err = json.Unmarshal([]byte(data.Message), &obj)

	if err != nil {
		return nil, err
	}
	if obj.PublicKey != data.PublicKey {
		return nil, errors.New("invalid message content")
	}

	return &getTransferChannelState{
		BaseRequest:    baseRequest,
		channelId:      obj.ChannelId,
		ledger:         ledger,
		trackerKeyPair: keyPair,
	}, nil
}

func (req *getTransferChannelState) Handle() (interface{}, error) {
	state, err := req.ledger.GetTransferChannelState(req.channelId)
	if err != nil {
		return nil, err
	}

	trackerPublicKey, err := req.trackerKeyPair.PublicKey()
	if err != nil {
		return nil, err
	}

	data := &getTransferChannelStateResponseData{
		ChannelId: req.channelId,
		State:     state,
		PublicKey: trackerPublicKey,
	}
	respData, err := helpers.NewResponseDataInterface(data, req.trackerKeyPair)
	if err != nil {
		return nil, err
	}

	return respData, nil
}
