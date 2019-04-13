package handler

import (
	"encoding/json"

	"github.com/pkg/errors"
	"svcledger/helpers"
	"svcledger/netHelpers"
	"svcledger/store"
	"svcnodr/types"
)

type getIncomeAddressState struct {
	*netHelpers.BaseRequest
	ledger         *store.Ledger
	peerAddress    string
	trackerKeyPair helpers.KeyPair
}

type getIncomeAddressStateRequestData struct {
	Address   string `json:"address"`
	PublicKey string `json:"pk"`
}

type getIncomeAddressStateResponseData struct {
	Address   string              `json:"address"`
	Current   *types.ChannelFact  `json:"current"`
	Limit     *types.ChannelPlan  `json:"limit"`
	Price     *types.ChannelPrice `json:"price"`
	PublicKey string              `json:"pk"`
	TimeLock  int64               `json:"timelock"`
}

func newGetIncomeAddressStateRequest(
	baseRequest *netHelpers.BaseRequest,
	keyPair helpers.KeyPair,
	ledger *store.Ledger,
) (netHelpers.Requester, error) {
	data, err := helpers.GetSignedPayloadData(baseRequest.Payload)

	if err != nil {
		return nil, err
	}

	var obj getIncomeAddressStateRequestData
	err = json.Unmarshal([]byte(data.Message), &obj)

	if err != nil {
		return nil, err
	}
	if obj.PublicKey != data.PublicKey {
		return nil, errors.New("invalid message content")
	}

	return &getIncomeAddressState{
		BaseRequest:    baseRequest,
		ledger:         ledger,
		peerAddress:    obj.Address,
		trackerKeyPair: keyPair,
	}, nil
}

func (req *getIncomeAddressState) Handle() (interface{}, error) {
	state, err := req.ledger.GetIncomeAddressState(req.peerAddress)
	if err != nil {
		return nil, err
	}

	trackerPublicKey, err := req.trackerKeyPair.PublicKey()
	if err != nil {
		return nil, err
	}

	data := &getIncomeAddressStateResponseData{
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
