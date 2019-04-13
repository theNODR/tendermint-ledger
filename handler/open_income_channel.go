package handler

import (
	"github.com/pkg/errors"

	"svcledger/helpers"
	"svcledger/netHelpers"
	"svcledger/store"
)

type openIncomeChannelRequest struct {
	*netHelpers.BaseRequest
	ledger				*store.Ledger
	peerPublicKey		string
	trackerKeyPair		helpers.KeyPair
}

type openIncomeChannelResponseData struct {
	Address		string				`json:"address"`
	Current		*store.ChannelFact	`json:"current"`
	Limit		*store.ChannelPlan	`json:"limit"`
	Price		*store.ChannelPrice	`json:"price"`
	PublicKey	string				`json:"pk"`
	TimeLock	int64				`json:"timelock"`
	LifeTime	int64				`json:"lifetime"`
}

func newOpenIncomeChannelRequest(
	baseRequest *netHelpers.BaseRequest,
	keyPair helpers.KeyPair,
	ledger *store.Ledger,
	_ *store.Queries,
) (netHelpers.Requester, error) {
	data, err := helpers.GetSignedPayloadData(baseRequest.Payload)

	if err != nil {
		return nil, err
	}
	if data.Message != data.PublicKey {
		return nil, errors.New("invalid message")
	}

	return &openIncomeChannelRequest{
		BaseRequest: baseRequest,
		ledger: ledger,
		peerPublicKey: data.PublicKey,
		trackerKeyPair: keyPair,
	}, nil
}

func (req *openIncomeChannelRequest) Handle() (interface{}, error) {
	state, err := req.ledger.OpenIncomeChannel(req.peerPublicKey)
	if err != nil {
		return nil, err
	}

	trackerPublicKey, err := req.trackerKeyPair.PublicKey()
	if err != nil {
		return nil, err
	}

	data := &openIncomeChannelResponseData{
		Address: state.Address,
		Current: state.State.Fact.ToChannel(),
		Limit: state.State.Plan.ToChannel(),
		Price: state.Price.ToChannel(),
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
