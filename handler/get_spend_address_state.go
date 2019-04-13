package handler

import (
	"encoding/json"

	"github.com/pkg/errors"
	"svcledger/helpers"
	"svcledger/netHelpers"
	"svcledger/store"
	"svcnodr/types"
)

type getSpendAddressState struct {
	*netHelpers.BaseRequest
	ledger         *store.Ledger
	peerAddress    string
	trackerKeyPair helpers.KeyPair
}

type getSpendAddressStateRequestData struct {
	Address   string `json:"address"`
	PublicKey string `json:"pk"`
}

type getSpendAddressStateResponseData struct {
	Address   string              `json:"address"`
	Current   *types.ChannelFact  `json:"current"`
	Limit     *types.ChannelPlan  `json:"limit"`
	Price     *types.ChannelPrice `json:"price"`
	PublicKey string              `json:"pk"`
	TimeLock  int64               `json:"timelock"`
}

func newGetSpendAddressStateRequest(
	baseRequest *netHelpers.BaseRequest,
	keyPair helpers.KeyPair,
	ledger *store.Ledger,
) (netHelpers.Requester, error) {
	data, err := helpers.GetSignedPayloadData(baseRequest.Payload)

	if err != nil {
		return nil, err
	}

	var obj getSpendAddressStateRequestData
	err = json.Unmarshal([]byte(data.Message), &obj)

	if err != nil {
		return nil, err
	}
	if obj.PublicKey != data.PublicKey {
		return nil, errors.New("invalid message content")
	}

	return &getSpendAddressState{
		BaseRequest:    baseRequest,
		ledger:         ledger,
		peerAddress:    obj.Address,
		trackerKeyPair: keyPair,
	}, nil
}

func (req *getSpendAddressState) Handle() (interface{}, error) {
	state, err := req.ledger.GetSpendAddressState(req.peerAddress)
	if err != nil {
		return nil, err
	}

	trackerPublicKey, err := req.trackerKeyPair.PublicKey()
	if err != nil {
		return nil, err
	}

	data := &getSpendAddressStateResponseData{
		Address:   state.Address,
		Current:   state.Current,
		Limit:     state.Limit,
		Price:     state.Price,
		PublicKey: trackerPublicKey,
		TimeLock:  state.TimeLock,
	}
	respData, err := helpers.NewResponseDataInterface(data, req.trackerKeyPair)
	if err != nil {
		return nil, err
	}

	return respData, nil
}
