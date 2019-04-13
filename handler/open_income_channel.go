package handler

import (
	"fmt"
	"github.com/pkg/errors"
	"svcnodr/types"

	"svcledger/helpers"
	"svcledger/netHelpers"
	"svcledger/store"
)

type openIncomeChannelRequest struct {
	*netHelpers.BaseRequest
	ledger         *store.Ledger
	peerPublicKey  string
	trackerKeyPair helpers.KeyPair
}

type openIncomeChannelResponseData struct {
	Address   string              `json:"address"`
	Current   *types.ChannelFact  `json:"current"`
	Limit     *types.ChannelPlan  `json:"limit"`
	Price     *types.ChannelPrice `json:"price"`
	PublicKey string              `json:"pk"`
	TimeLock  int64               `json:"timelock"`
	LifeTime  int64               `json:"lifetime"`
}

func newOpenIncomeChannelRequest(
	baseRequest *netHelpers.BaseRequest,
	keyPair helpers.KeyPair,
	ledger *store.Ledger,
) (netHelpers.Requester, error) {
	data, err := helpers.GetSignedPayloadData(baseRequest.Payload)

	if err != nil {
		return nil, err
	}
	if data.Message != data.PublicKey {
		return nil, errors.New("invalid message")
	}

	return &openIncomeChannelRequest{
		BaseRequest:    baseRequest,
		ledger:         ledger,
		peerPublicKey:  data.PublicKey,
		trackerKeyPair: keyPair,
	}, nil
}

func (req *openIncomeChannelRequest) Handle() (interface{}, error) {
	state, err := req.ledger.OpenIncomeChannel(req.peerPublicKey)
	if err != nil {
		fmt.Printf("Close income channel err:%v", err)
		return nil, err
	}
	fmt.Printf("Close income channel:%v", state)

	trackerPublicKey, err := req.trackerKeyPair.PublicKey()
	if err != nil {
		return nil, err
	}

	data := &openIncomeChannelResponseData{
		Address:   state.Address,
		Current:   state.Current,
		Limit:     state.Limit,
		Price:     state.Price,
		PublicKey: trackerPublicKey,
		TimeLock:  state.TimeLock,
		LifeTime:  state.LifeTime,
	}
	respData, err := helpers.NewResponseDataInterface(data, req.trackerKeyPair)
	if err != nil {
		return nil, err
	}

	return respData, nil
}
