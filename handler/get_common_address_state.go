package handler

import (
	"encoding/json"

	"github.com/pkg/errors"
	"svcledger/helpers"
	"svcledger/netHelpers"
	"svcledger/store"
	ntypes "svcnodr/types"
)

type getCommonAddressState struct {
	*netHelpers.BaseRequest

	address        string
	ledger         *store.Ledger
	trackerKeyPair helpers.KeyPair
}

type getCommonAddressStateRequestData struct {
	Address   string `json:"address"`
	PublicKey string `json:"pk"`
}

type getCommonAddressStateResponseData struct {
	Address   string                              `json:"address"`
	State     *ntypes.TrackerBalanceStateResponse `json:"state"`
	PublicKey string                              `json:"pk"`
}

func newGetCommonAddressState(
	baseRequest *netHelpers.BaseRequest,
	keyPair helpers.KeyPair,
	ledger *store.Ledger,
) (netHelpers.Requester, error) {
	data, err := helpers.GetSignedPayloadData(baseRequest.Payload)
	if err != nil {
		return nil, err
	}

	var obj getCommonAddressStateRequestData
	err = json.Unmarshal([]byte(data.Message), &obj)
	if err != nil {
		return nil, err
	}
	if data.PublicKey != obj.PublicKey {
		return nil, errors.New("invalid message content")
	}
	if !helpers.IsAddressCommon(obj.Address) {
		return nil, errors.New("invalid common address")
	}

	return &getCommonAddressState{
		BaseRequest:    baseRequest,
		address:        obj.Address,
		ledger:         ledger,
		trackerKeyPair: keyPair,
	}, nil
}

func (req *getCommonAddressState) Handle() (interface{}, error) {
	state, err := req.ledger.GetCommonAddressBalance(req.address)
	if err != nil {
		return nil, err
	}

	trackerPublicKey, err := req.trackerKeyPair.PublicKey()
	if err != nil {
		return nil, err
	}

	data := &getCommonAddressStateResponseData{
		Address:   req.address,
		State:     state,
		PublicKey: trackerPublicKey,
	}
	respData, err := helpers.NewResponseDataInterface(data, req.trackerKeyPair)
	if err != nil {
		return nil, err
	}

	return respData, nil
}
