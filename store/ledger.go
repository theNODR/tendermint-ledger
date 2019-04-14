package store

import (
	"svcledger/blockchain/tendermint"
	"svcnodr/handlers"
	ntypes "svcnodr/types"
	"sync"
	"time"

	"common"
	"svcledger/blockchain"
	"svcledger/blockchain/trans"
	"svcledger/helpers"
)

const commonAddress string = ""
const defaultTransactionsMaxCount uint64 = 10000

type ownChannelStatusType = uint

const (
	notOwnChannelStatus     ownChannelStatusType = 1 << iota
	ownChannelStatus        ownChannelStatusType = 1 << iota
	invalidOwnChannelStatus ownChannelStatusType = 0
)

type Ledger struct {
	address string
	keyPair helpers.KeyPair
	mutex   sync.RWMutex
	pk      string
	client  *tendermint.TendermintClient

	amounts   *Amounts
	lifeTimes *LifeTimes
}

func NewLedger(
	address string,
	keyPair helpers.KeyPair,
	amounts *Amounts,
	lifeTimes *LifeTimes,
	client *tendermint.TendermintClient,
) (*Ledger, error) {
	pk, err := keyPair.PublicKey()
	if err != nil {
		return nil, err
	}

	l := &Ledger{
		address: address,
		keyPair: keyPair,
		pk:      pk,
		client:  client,

		amounts:   amounts,
		lifeTimes: lifeTimes,
	}

	return l, nil
}

func (l *Ledger) Init(amount TransactionAmountType) (TransactionStatus, error) {
	tran := blockchain.NewInitTransactionType(
		&trans.InitTran{
			To:     l.address,
			Amount: uint64(amount),
		},
		l.pk,
	)

	err := l.client.SendInit(tran)
	for err != nil {
		time.Sleep(5 * time.Second)
		common.Log.Error("ErrorInitLedgerBalance", common.Printf("Can`t init ledger balance. %v", err))
		err = l.client.SendInit(tran)
	}

	if err != nil {
		return InvalidTransactionStatus, err
	}

	return ValidTransactionStatus, nil
}

func (l *Ledger) AutoClose() error {
	tran := blockchain.NewAutoCloseTran(
		&trans.AutoCloseTran{
			Timestamp:     common.GetNowUnixMs(),
			LedgerAddress: l.address,
		},
		l.pk,
	)

	return l.client.SendAutoclose(tran)
}

func (l *Ledger) OpenIncomeChannel(peerPublicKey string) (*ntypes.IncomeChannelStateResponse, error) {
	tran := blockchain.NewOpenIncomeTran(
		&trans.OpenIncomeTran{
			PeerPublicKey:     peerPublicKey,
			From:              l.address,
			LifeTime:          l.lifeTimes.IncomeChannel,
			PriceQuantumPower: uint8(l.amounts.PriceQuantumPower),
			PriceAmount:       uint64(l.amounts.PriceIncome),
		},
		l.pk,
	)

	state, err := l.client.SendOpenIncomeChannel(tran)
	if err != nil {
		return nil, err
	}

	return state, nil
}

func (l *Ledger) GetIncomeAddressState(address string) (*ntypes.IncomeChannelStateResponse, error) {
	state, err := l.client.GetIncomeState(address)
	if err != nil {
		return nil, err
	}

	return state, nil
}

func (l *Ledger) OpenSpendChannel(peerPublicKey string) (*ntypes.SpendChannelStateResponse, error) {
	tran := blockchain.NewOpenSpendTran(
		&trans.OpenSpendTran{
			PeerPublicKey:     peerPublicKey,
			From:              l.address,
			MaxAmount:         uint64(l.amounts.Spend),
			PriceQuantumPower: uint8(l.amounts.PriceQuantumPower),
			PriceAmount:       uint64(l.amounts.PriceSpend),
			LifeTime:          l.lifeTimes.SpendChannel,
		},
		l.pk,
	)

	state, err := l.client.SendOpenSpendChannel(tran)
	if err != nil {
		return nil, err
	}
	return state, nil
}

func (l *Ledger) GetSpendAddressState(address string) (*ntypes.SpendChannelStateResponse, error) {
	state, err := l.client.GetSpendChannelState(address)
	if err != nil {
		return nil, err
	}

	return state, err
}

func (l *Ledger) GetCommonAddressBalance(address string) (*ntypes.TrackerBalanceStateResponse, error) {
	return nil, nil
}

func (l *Ledger) OpenTransferChannel(
	channelId string,
	fromAddress string,
	fromPk string,
	toAddress string,
	toPk string,
	priceAmount TransactionAmountType,
	priceQuantumPower QuantumPowerType,
	plannedQuantumCount uint64,
) (*handlers.OpenTransferResponse, error) {
	tran := blockchain.NewOpenTransferTran(
		&trans.OpenTransferTran{
			ChannelId:           channelId,
			FromAddress:         fromAddress,
			FromPk:              fromPk,
			ToAddress:           toAddress,
			ToPk:                toPk,
			PriceAmount:         uint64(priceAmount),
			PriceQuantumPower:   uint8(priceQuantumPower),
			PlannedQuantumCount: plannedQuantumCount,
			LifeTime:            l.lifeTimes.TransferChannel,
		},
		l.pk,
	)

	state, err := l.client.SendOpenTransferChannel(tran)
	if err != nil {
		return nil, err
	}
	return state, nil
}

func (l *Ledger) GetTransferChannelState(channelId string) (ntypes.TransferChannelStatus, error) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	state, err := l.client.GetTransferChannelState(channelId)
	if err != nil {
		return ntypes.ErrorTransferChannelStatus, err
	}

	return state.Status, nil
}

func (l *Ledger) CloseTransferChannel(
	channelId string,
	fromAddress string,
	fromPk string,
	toAddress string,
	toPk string,
	priceAmount TransactionAmountType,
	priceQuantumPower QuantumPowerType,
	plannedQuantumCount uint64,
	amount TransactionAmountType,
	quantumCount uint64,
	volume uint64,
) (TransactionStatus, error) {
	tran := blockchain.NewCloseTransferTran(
		&trans.CloseTransferTran{
			ChannelId:           channelId,
			FromAddress:         fromAddress,
			FromPk:              fromPk,
			ToAddress:           toAddress,
			ToPk:                toPk,
			PriceAmount:         uint64(priceAmount),
			PriceQuantumPower:   uint8(priceQuantumPower),
			PlannedQuantumCount: plannedQuantumCount,
			Amount:              uint64(amount),
			QuantumCount:        quantumCount,
			Volume:              volume,
		},
		l.pk,
	)

	var err error
	err = l.client.SendCloseTransfer(tran)
	if err != nil {
		return InvalidTransactionStatus, err
	}

	return ValidTransactionStatus, nil
}

func (l *Ledger) CloseIncomeChannel(address string, pk string) (TransactionStatus, *ntypes.ChannelFact, error) {
	errResult := func(err error) (TransactionStatus, *ntypes.ChannelFact, error) {
		common.Log.Print(common.Printf(
			"CloseIncomeChannel: %v",
			err.Error(),
		))
		return InvalidTransactionStatus, nil, err
	}

	tran := blockchain.NewCloseIncomeTran(
		&trans.CloseIncomeTran{
			From:          address,
			To:            l.address,
			PeerPublicKey: pk,
		},
		l.pk,
	)

	state, err := l.client.SendCloseIncome(tran)
	if err != nil {
		return errResult(err)
	}

	return ValidTransactionStatus, state.State, nil
}

func (l *Ledger) CloseSpendChannel(address string, pk string) (TransactionStatus, *ntypes.ChannelFact, error) {
	errResult := func(err error) (TransactionStatus, *ntypes.ChannelFact, error) {
		common.Log.Print(common.Printf(
			"CloseSpendChannel: %v",
			err.Error(),
		))
		return InvalidTransactionStatus, nil, err
	}
	tran := blockchain.NewCloseSpendTran(
		&trans.CloseSpendTran{
			PeerPublicKey: pk,
			From:          address,
			To:            l.address,
		},
		l.pk,
	)

	state, err := l.client.SendCloseSpend(tran)
	if err != nil {
		return errResult(err)
	}

	return ValidTransactionStatus, state.State, nil
}
