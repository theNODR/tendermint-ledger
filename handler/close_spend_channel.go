package handler

import (
	"encoding/json"

	"github.com/pkg/errors"

	"svcledger/helpers"
	"svcledger/netHelpers"
	"svcledger/store"
)

type closeSpendChannelRequest struct {
	*netHelpers.BaseRequest
	ledger				*store.Ledger
	peerAddress			string
	peerPublicKey		string
	trackerKeyPair		helpers.KeyPair
}

type closeSpendAddressRequestData struct {
	Address		string	`json:"address"`
	PublicKey	string	`json:"pk"`
}

type closeSpendChannelResponseData struct {
	PublicKey	string				`json:"pk"`
	State		*store.ChannelFact	`json:"state"`
}

func newCloseSpendChannelRequest(
	baseRequest *netHelpers.BaseRequest,
	keyPair helpers.KeyPair,
	ledger *store.Ledger,
	_ *store.Queries,
) (netHelpers.Requester, error) {
	data, err := helpers.GetSignedPayloadData(baseRequest.Payload)

	if err != nil {
		return nil, err
	}

	var obj closeSpendAddressRequestData
	err = json.Unmarshal([]byte(data.Message), &obj)

	if err != nil {
		return nil, err
	}
	if obj.PublicKey != data.PublicKey {
		return nil, errors.New("invalid message content")
	}

	return &closeSpendChannelRequest{
		BaseRequest: baseRequest,
		peerAddress: obj.Address,
		peerPublicKey: obj.PublicKey,
		trackerKeyPair: keyPair,
		ledger: ledger,
	}, nil
}

func (req *closeSpendChannelRequest) Handle() (interface{}, error) {
	_, state, err := req.ledger.CloseSpendChannel(req.peerAddress, req.peerPublicKey)
	if err != nil {
		return nil, err
	}

	trackerPublicKey, err := req.trackerKeyPair.PublicKey()
	if err != nil {
		return err, nil
	}

	data := &closeSpendChannelResponseData{
		PublicKey: trackerPublicKey,
		State: state,
	}
	respData, err := helpers.NewResponseDataInterface(data, req.trackerKeyPair)
	if err != nil {
		return nil, err
	}

	return respData, nil
}
