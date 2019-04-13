package handler

import (
	"github.com/pkg/errors"

	"svcledger/helpers"
	"svcledger/netHelpers"
	"svcledger/store"
)

type handshakeRequest struct {
	*netHelpers.BaseRequest
	trackerKeyPair	helpers.KeyPair
}

func newHandshakeRequest(
	baseRequest *netHelpers.BaseRequest,
	keyPair helpers.KeyPair,
	_ *store.Ledger,
	_ *store.Queries,
) (netHelpers.Requester, error) {
	data, err := helpers.GetSignedPayloadData(baseRequest.Payload)

	if err != nil {
		return nil, err
	}
	if data.Message != data.PublicKey {
		return nil, errors.New("invalid message")
	}

	return &handshakeRequest{
		BaseRequest: baseRequest,
		trackerKeyPair: keyPair,
	}, nil
}

func (req *handshakeRequest) Handle() (interface{}, error) {
	trackerPublicKey, err := req.trackerKeyPair.PublicKey()
	if err != nil {
		return nil, err
	}

	data, err := helpers.NewResponseDataString(trackerPublicKey, req.trackerKeyPair)
	if err != nil {
		return nil, err
	}

	return data, nil
}
