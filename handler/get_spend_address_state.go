package handler

import (
	"encoding/json"

	"github.com/pkg/errors"

	"svcledger/helpers"
	"svcledger/netHelpers"
	"svcledger/store"
)

type getSpendAddressState struct {
	*netHelpers.BaseRequest
	ledger			*store.Ledger
	peerAddress		string
	trackerKeyPair	helpers.KeyPair
}

type getSpendAddressStateRequestData struct {
	Address		string	`json:"address"`
	PublicKey	string	`json:"pk"`
}

type getSpendAddressStateResponseData struct {
	Address		string				`json:"address"`
	Current		*store.ChannelFact	`json:"current"`
	Limit		*store.ChannelPlan	`json:"limit"`
	Price		*store.ChannelPrice	`json:"price"`
	PublicKey	string				`json:"pk"`
	TimeLock	int64				`json:"timelock"`
}

func newGetSpendAddressStateRequest(
	baseRequest *netHelpers.BaseRequest,
	keyPair helpers.KeyPair,
	ledger *store.Ledger,
	_ *store.Queries,
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
		BaseRequest: baseRequest,
		ledger: ledger,
		peerAddress: obj.Address,
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
		Address: state.Address,
		Current: state.State.Fact.ToChannel(),
		Limit: state.State.Plan.ToChannel(),
		Price: state.Price.ToChannel(),
		PublicKey: trackerPublicKey,
		TimeLock: state.TimeLock,
	}
	respData, err := helpers.NewResponseDataInterface(data, req.trackerKeyPair)
	if err != nil {
		return nil, err
	}

	return respData, nil
}
