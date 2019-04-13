package handler

import (
	"encoding/json"
	"svcnodr/types"

	"github.com/pkg/errors"

	"svcledger/helpers"
	"svcledger/netHelpers"
	"svcledger/store"
)

type closeIncomeChannelRequest struct {
	*netHelpers.BaseRequest
	ledger         *store.Ledger
	peerAddress    string
	peerPublicKey  string
	trackerKeyPair helpers.KeyPair
}

type closeIncomeAddressRequestData struct {
	Address   string `json:"address"`
	PublicKey string `json:"pk"`
}

type closeIncomeChannelResponseData struct {
	PublicKey string             `json:"pk"`
	State     *types.ChannelFact `json:"state"`
}

func newCloseIncomeChannelRequest(
	baseRequest *netHelpers.BaseRequest,
	keyPair helpers.KeyPair,
	ledger *store.Ledger,
) (netHelpers.Requester, error) {
	data, err := helpers.GetSignedPayloadData(baseRequest.Payload)

	if err != nil {
		return nil, err
	}

	var obj closeIncomeAddressRequestData
	err = json.Unmarshal([]byte(data.Message), &obj)

	if err != nil {
		return nil, err
	}
	if obj.PublicKey != data.PublicKey {
		return nil, errors.New("invalid message content")
	}

	return &closeIncomeChannelRequest{
		BaseRequest:    baseRequest,
		ledger:         ledger,
		peerAddress:    obj.Address,
		peerPublicKey:  obj.PublicKey,
		trackerKeyPair: keyPair,
	}, nil
}

func (req *closeIncomeChannelRequest) Handle() (interface{}, error) {
	_, state, err := req.ledger.CloseIncomeChannel(req.peerAddress, req.peerPublicKey)
	if err != nil {
		return nil, err
	}

	trackerPublicKey, err := req.trackerKeyPair.PublicKey()
	if err != nil {
		return nil, err
	}

	data := &closeIncomeChannelResponseData{
		PublicKey: trackerPublicKey,
		State:     state,
	}
	respData, err := helpers.NewResponseDataInterface(data, req.trackerKeyPair)
	if err != nil {
		return nil, err
	}

	return respData, nil
}
