package store

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"sort"
	"sync"

	"github.com/pkg/errors"

	"common"
	"svcledger/blockchain"
	"svcledger/blockchain/trans"
	"svcledger/helpers"
	"svcledger/logger"
)

const commonAddress					string = ""
const defaultTransactionsMaxCount	uint64 = 10000

type TransferChannelStatus = uint

const (
	NotExistTransferChannelStatus	TransferChannelStatus = 1 << iota
	OpenTransferChannelStatus		TransferChannelStatus = 1 << iota
	CloseTransferChannelStatus		TransferChannelStatus = 1 << iota
	ErrorTransferChannelStatus		TransferChannelStatus = 0
)

type ownChannelStatusType = uint
const (
	notOwnChannelStatus ownChannelStatusType = 1 << iota
	ownChannelStatus ownChannelStatusType = 1 << iota
	invalidOwnChannelStatus ownChannelStatusType = 0
)

type FirstTransactions struct {
	Data		LedgerTransactions
	RowCount	uint
}

type NextTransactions struct {
	Data	LedgerTransactions
}

type Ledger struct {
	address			string
	financialState	FinancialState
	keyPair			helpers.KeyPair
	mutex			sync.RWMutex
	pk				string
	transactions	LedgerTransactions
	waiters			Waiters
	writer			blockchain.Writer

	amounts			*Amounts
	lifeTimes		*LifeTimes
}

func NewLedger(
	address string,
	keyPair helpers.KeyPair,
	trans blockchain.TranItems,
	amounts *Amounts,
	lifeTimes *LifeTimes,
	waiters Waiters,
	writer blockchain.Writer,
) (*Ledger, error) {
	pk, err := keyPair.PublicKey()
	if err != nil {
		return nil, err
	}

	transactions := make(LedgerTransactions, 0, defaultTransactionsMaxCount)
	finState := NewFinancialState(address, pk)

	l := &Ledger{
		address: address,
		financialState: finState,
		keyPair: keyPair,
		pk: pk,
		transactions: transactions,
		waiters: waiters,
		writer: writer,

		amounts: amounts,
		lifeTimes: lifeTimes,
	}

	l.AddTrans(trans)

	return l, nil
}

func (l *Ledger) AddTrans(tranItems blockchain.TranItems) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	addTransCount := 0

	for _, item := range tranItems {
		var tran blockchain.ReadTran
		err := json.Unmarshal([]byte(item.Tran), &tran)
		if err != nil {
			continue
		}

		blockchainTran := &BlockchainTran{
			Id: tran.Id,
			ClientTimestamp: tran.TimeStamp,
			PublicKey: tran.PublicKey,
			ServerTimestamp: item.TimeStamp,
		}

		switch tran.Type{
		case blockchain.InitTransactionType:
			var d trans.InitTran
			err = json.Unmarshal(tran.Data, &d)
			if err != nil {
				common.Log.PrintFull(
					common.Printf("data conversion error: %v", err.Error()),
				)
				continue
			}
			addTransCount++
			l.init(blockchainTran, &d)
		case blockchain.AutoCloseTransactionType:
			var d trans.AutoCloseTran
			err = json.Unmarshal(tran.Data, &d)
			if err != nil {
				common.Log.PrintFull(
					common.Printf("data conversion error: %v", err.Error()),
				)
				continue
			}
			addTransCount++
			l.autoClose(blockchainTran, &d)
		case blockchain.FinancialTransactionType:
			var d trans.FinancialTran
			err = json.Unmarshal(tran.Data, &d)
			if err != nil {
				common.Log.PrintFull(
					common.Printf("data conversion error: %v", err.Error()),
				)
				continue
			}
			addTransCount++
			l.writeFinancialTransaction(blockchainTran, &d)
		case blockchain.OpenIncomeChannelTransactionType:
			var d trans.OpenIncomeTran
			err = json.Unmarshal(tran.Data, &d)
			if err != nil {
				common.Log.PrintFull(
					common.Printf("data conversion error: %v", err.Error()),
				)
				continue
			}
			addTransCount++
			l.openIncomeChannel(blockchainTran, &d)
		case blockchain.OpenSpendChannelTransactionType:
			var d trans.OpenSpendTran
			err = json.Unmarshal(tran.Data, &d)
			if err != nil {
				common.Log.PrintFull(
					common.Printf("data conversion error: %v", err.Error()),
				)
				continue
			}
			addTransCount++
			l.openSpendChannel(blockchainTran, &d)
		case blockchain.OpenTransferChannelTransactionType:
			var d trans.OpenTransferTran
			err = json.Unmarshal(tran.Data, &d)
			if err != nil {
				common.Log.PrintFull(
					common.Printf("data conversion error: %v", err.Error()),
				)
				continue
			}
			addTransCount++
			l.openTransferChannel(blockchainTran, &d)
		case blockchain.CloseTransferChannelTransactionType:
			var d trans.CloseTransferTran
			err = json.Unmarshal(tran.Data, &d)
			if err != nil {
				common.Log.PrintFull(
					common.Printf("data conversion error: %v", err.Error()),
				)
				continue
			}
			addTransCount++
			l.closeTransferChannel(blockchainTran, &d)
		case blockchain.CloseSpendChannelTransactionType:
			var d trans.CloseSpendTran
			err = json.Unmarshal(tran.Data, &d)
			if err != nil {
				common.Log.PrintFull(
					common.Printf("data conversion error: %v", err.Error()),
				)
				continue
			}
			addTransCount++
			l.closeSpendChannel(blockchainTran, &d)
		case blockchain.CloseIncomeChannelTransactionType:
			var d trans.CloseIncomeTran
			err = json.Unmarshal(tran.Data, &d)
			if err != nil {
				common.Log.PrintFull(
					common.Printf("data conversion error: %v", err.Error()),
				)
				continue
			}
			addTransCount++
			l.closeIncomeChannel(blockchainTran, &d)
		//case UnknownTransactionType:
		default:
			continue
		}
	}

	if addTransCount > 0 {
		common.Log.PrintFull(
			common.Printf("add trans: %v", addTransCount),
		)
	}
}

func (l *Ledger) Init(amount TransactionAmountType) (TransactionStatus, error) {
	tran := blockchain.NewInitTransactionType(
		&trans.InitTran{
			To: l.address,
			Amount: uint64(amount),
		},
		l.pk,
	)

	ch := make(chan *LedgerTransaction)
	defer func() {
		close(ch)
	}()

	err := l.waiters.Add(tran.Id, ch)
	if err != nil {
		return InvalidTransactionStatus, err
	}

	err = l.sendTran(tran)
	if err != nil {
		return InvalidTransactionStatus, err
	}

	ledgerTran := <-ch
	if ledgerTran == nil {
		return InvalidTransactionStatus, errors.New("ledger: transaction not synced")
	}

	return ledgerTran.Status, nil
}

func (l *Ledger) init(bData *BlockchainTran, tData *trans.InitTran) {
	success := false
	defer func() {
		if !success {
			l.waiters.TryInterrupt(bData.Id)
		}
	}()
	body := &TransactionBody{
		TimeStamp: bData.ClientTimestamp,
		Type: InitTransactionType,
		From: commonAddress,
		To:  tData.To,
		Amount: TransactionAmountType(tData.Amount),
	}

	tran, err := l.createTran(body, bData.Id, bData.PublicKey)
	if err != nil {
		common.Log.PrintFull(
			common.Printf("fail to create init tran: %v", err.Error()),
		)
		return
	}

	ledgerTran := &LedgerTransaction{
		SignedTransaction: tran,
		ServerTimestamp: bData.ServerTimestamp,
	}

	if l.hasValidTrans() {
		ledgerTran.Status = NoneTransactionStatus
	} else {
		ledgerTran.Status = ValidTransactionStatus
	}

	l.addTran(ledgerTran)
	common.Log.PrintFull(
		common.Printf("success init tran: id=%v, to=%v, status=%v", ledgerTran.Id, ledgerTran.To, ledgerTran.Status),
	)
	success = true
}

func (l *Ledger) AutoClose() error {
	tran := blockchain.NewAutoCloseTran(
		&trans.AutoCloseTran{
		},
		l.pk,
	)

	return l.sendTran(tran)
}

func (l *Ledger) autoClose(bData *BlockchainTran, tData *trans.AutoCloseTran) {
	common.Log.Event(
		logger.EventLedgerAutoCloseBegin,
		common.Printf("begin autoclose", ),
	)

	ms := newMemoryState(bData.PublicKey, bData.ClientTimestamp)
	ms.Init(l.transactions)

	openTransfers := ms.GetTransferChannelsForClose()

	common.Log.PrintFull(common.Printf("begin autoclose transfers: %v", len(openTransfers), ), )
	for _, item := range openTransfers {
		openTran := item.Tran()
		signedTran, err := l.createCloseTransferChannelTran(
			openTran.Id,
			openTran.ContractId,
			openTran.From,
			openTran.To,
			openTran.PriceAmount,
			openTran.PriceQuantumPower,
			openTran.PlannedQuantumCount,
			0.0,
			0,
			0,
			bData.ClientTimestamp,
			blockchain.CreateId(bData.PublicKey, bData.ClientTimestamp),
			bData.PublicKey,
		)

		if err != nil {
			continue
		}

		item.Close(nil)
		l.addTran(&LedgerTransaction{
			SignedTransaction: signedTran,
			ServerTimestamp: bData.ServerTimestamp,
			Status: ValidTransactionStatus,
		})
	}

	openIncomes := ms.GetIncomeAddressesForClose()
	common.Log.PrintFull(common.Printf("begin autoclose incomes: %v", len(openIncomes), ), )
	for _,item := range openIncomes {
		openTran := item.Tran()
		signedTran, err := l.createCloseIncomeChannelTran(
			openTran.Id,
			openTran.To,
			openTran.From,
			item.Amount(),
			item.Volume(),
			bData.ClientTimestamp,
			blockchain.CreateId(bData.PublicKey, bData.ClientTimestamp),
			bData.PublicKey,
			nil,
		)

		if err != nil {
			continue
		}

		l.addTran(&LedgerTransaction{
			SignedTransaction: signedTran,
			ServerTimestamp: bData.ServerTimestamp,
			Status: ValidTransactionStatus,
		})
	}

	openSpends := ms.GetSpendAddressesForClose()
	common.Log.PrintFull(common.Printf("begin autoclose spends: %v", len(openSpends), ), )
	for _, item := range openSpends {
		openTran := item.Tran()
		signedTran, err := l.createCloseSpendChannelTran(
			openTran.Id,
			openTran.To,
			openTran.From,
			item.Amount(),
			item.Volume(),
			bData.ClientTimestamp,
			blockchain.CreateId(bData.PublicKey, bData.ClientTimestamp),
			bData.PublicKey,
			nil,
		)

		if err != nil {
			continue
		}

		l.addTran(&LedgerTransaction{
			SignedTransaction: signedTran,
			ServerTimestamp: bData.ServerTimestamp,
			Status: ValidTransactionStatus,
		})
	}

	common.Log.Event(
		logger.EventLedgerAutoCloseEnd,
		common.Printf("autoclose end", ),
	)
}

func (l *Ledger) WriteFinancialTransaction() error {
	tran := blockchain.NewFinancialTran(
		&trans.FinancialTran{
			From: l.address,
			To: commonAddress,
		},
		l.pk,
	)

	return l.sendTran(tran)
}

//ToDo: Not worked. Add FinancialStates class to calc all financial states, not only own
func (l *Ledger) writeFinancialTransaction(bData *BlockchainTran, tData *trans.FinancialTran) {
	success := false
	defer func() {
		if !success {
			l.waiters.TryInterrupt(bData.Id)
		}
	}()
	state := l.financialState.GetFinancialData()
	if state == nil {
		common.Log.Event(
			logger.EventLedgerFinTranEmpty,
			common.Printf("financial state data is empty", ),
		)
		return
	}

	data, err := json.Marshal(state)
	if err != nil {
		common.Log.StatError(
			logger.ErrorLedgerFinTranSerialize,
			common.Printf("financial state data serialize error: %v", err),
		)
		return
	}

	dataString := string(data[:])
	body := &TransactionBody{
		TimeStamp: bData.ClientTimestamp,
		Type: FinancialBlockTransactionType,
		From: tData.From,
		To: tData.To,
		Data: dataString,
	}

	signedTran, err := l.createTran(body, bData.Id, bData.PublicKey)
	if err != nil {
		return
	}

	ledgerTran := &LedgerTransaction{
		SignedTransaction: signedTran,
		Status: ValidTransactionStatus,
		ServerTimestamp: bData.ServerTimestamp,
	}

	l.addTran(ledgerTran)
	success = true

	common.Log.Print(common.Printf("%v", dataString, ))
	common.Log.Event(
		logger.EventLedgerFinTranSuccess,
		common.Printf("financial transaction write success", ),
	)
}

func (l *Ledger) GetFirstTransactions(
	pageSize int,
	page uint,
	fields *GetTransactionsFields,
	order OrderType,
) (*FirstTransactions, error) {
	if page == 0 {
		return nil, errors.New("ledger: GetFirstTransactions: page value should be greater than 0")
	}
	if pageSize == 0 {
		return nil, errors.New("ledger: GetFirstTransactions: page size value should be not equal 0")
	}

	copied := l.getCopiedTrans(fields)
	l.sort(copied, order)

	data, totalCount :=  l.getDataFromCopied(
		copied,
		pageSize,
		page,
	)

	return &FirstTransactions{
		Data: *data,
		RowCount: totalCount,
	}, nil
}

func (l *Ledger) GetNextTransactions(
	startId string,
	pageSize int,
	page uint,
	fields *GetTransactionsFields,
	order OrderType,
) (*NextTransactions, error) {
	if page == 0 {
		return nil, errors.New("ledger: GetNextTransactions: page value should be greater than 0")
	}
	if pageSize == 0 {
		return nil, errors.New("ledger: GetNextTransactions: page size value should be greater than 0")
	}

	arr := *l.getCopiedTrans(fields)

	var pos int = len(arr) - 1
	for ; pos >= 0; pos-- {
		if arr[pos].Id == startId {
			break
		}
	}
	if pos < 0 {
		return nil, errors.New("ledger: GetNextTransactions: invalid startId value")
	}

	arr = arr[:pos]

	copied := &arr
	l.sort(copied, order)

	data, _ :=  l.getDataFromCopied(
		copied,
		pageSize,
		page,
	)

	return &NextTransactions{
		Data: *data,
	}, nil
}

func (l *Ledger) OpenIncomeChannel(peerPublicKey string) (*IncomeChannelState, error) {
	tran := blockchain.NewOpenIncomeTran(
		&trans.OpenIncomeTran{
			PeerPublicKey: peerPublicKey,
			From: l.address,
			LifeTime: l.lifeTimes.IncomeChannel,
			PriceQuantumPower: uint8(l.amounts.PriceQuantumPower),
			PriceAmount: uint64(l.amounts.PriceIncome),
		},
		l.pk,
	)

	var err error
	ch := make(chan *LedgerTransaction)
	defer func() {
		close(ch)
	}()

	err = l.waiters.Add(tran.Id, ch)
	if err != nil {
		return nil, err
	}

	err = l.sendTran(tran)
	if err != nil {
		return nil, err
	}

	ledgerTran := <-ch
	if ledgerTran == nil {
		return nil, errors.New("ledger: transaction not synced")
	}

	return NewIncomeChannelState(ledgerTran), nil
}

func (l *Ledger) openIncomeChannel(bData *BlockchainTran, tData *trans.OpenIncomeTran) {
	success := false
	defer func() {
		if !success {
			l.waiters.TryInterrupt(bData.Id)
		}
	}()

	timestamp :=bData.ClientTimestamp
	address := helpers.CreateIncomeAddress(bData.PublicKey, tData.PeerPublicKey, timestamp)

	body := &TransactionBody{
		TimeStamp: timestamp,
		Type: OpenIncomeChannelTransactionType,
		From: tData.From,
		To: address,
		Amount: 0.0,
		PriceAmount: TransactionAmountType(tData.PriceAmount),
		PriceQuantumPower: QuantumPowerType(tData.PriceQuantumPower),
		TimeLock: timestamp + tData.LifeTime,
	}

	signedTran, err := l.createTran(body, bData.Id, bData.PublicKey)
	if err != nil {
		return
	}

	ledgerTran := &LedgerTransaction{
		SignedTransaction: signedTran,
		ServerTimestamp: bData.ServerTimestamp,
		Status: ValidTransactionStatus,
	}

	l.addTran(ledgerTran)
}

func (l *Ledger) GetIncomeAddressState(address string) (*IncomeChannelState, error) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if !helpers.IsAddressIncome(address) {
		return nil, errors.New(fmt.Sprintf("ledger: GetIncomeAddressState: invalid income address: %v", address))
	}

	pos, err := l.getOpenIncomeChannelPositionByAddress(address)

	if err != nil {
		return nil, err
	}

	if pos > -1 {
		return l.getIncomeChannelState(pos)
	}

	return nil, errors.New(fmt.Sprintf("ledger: GetIncomeAddressState: income address no find: %v", address))
}

func (l *Ledger) OpenSpendChannel(peerPublicKey string) (*SpendChannelState, error) {
	tran := blockchain.NewOpenSpendTran(
		&trans.OpenSpendTran{
			PeerPublicKey: peerPublicKey,
			From: l.address,
			MaxAmount: uint64(l.amounts.Spend),
			PriceQuantumPower: uint8(l.amounts.PriceQuantumPower),
			PriceAmount: uint64(l.amounts.PriceSpend),
			LifeTime: l.lifeTimes.SpendChannel,
		},
		l.pk,
	)

	var err error

	ch := make(chan *LedgerTransaction)
	defer func() {
		close(ch)
	}()

	err = l.waiters.Add(tran.Id, ch)
	if err != nil {
		return nil, err
	}

	err = l.sendTran(tran)
	if err != nil {
		return nil, err
	}

	ledgerTran := <-ch
	if ledgerTran == nil {
		return nil, errors.New("ledger: transaction not synced")
	}

	return NewSpendChannelState(ledgerTran), nil
}

func (l *Ledger) openSpendChannel(bData *BlockchainTran, tData *trans.OpenSpendTran) {
	success := false
	defer func() {
		if !success {
			l.waiters.TryInterrupt(bData.Id)
		}
	}()
	amount, err := l.getAddressAmount(tData.From)
	if err != nil {
		return
	}
	if amount <= 0 {
		return
	}
	if amount > TransactionAmountType(tData.MaxAmount) {
		amount = TransactionAmountType(tData.MaxAmount)
	}

	timestamp := bData.ClientTimestamp
	address := helpers.CreateSpendAddress(bData.PublicKey, tData.PeerPublicKey, timestamp)

	body := &TransactionBody{
		TimeStamp: timestamp,
		Type: OpenSpendChannelTransactionType,
		From: tData.From,
		To: address,
		Amount: amount,
		PriceAmount: TransactionAmountType(tData.PriceAmount),
		PriceQuantumPower: QuantumPowerType(tData.PriceQuantumPower),
		TimeLock: timestamp + tData.LifeTime,
	}

	signedTran, err := l.createTran(body, bData.Id, bData.PublicKey)
	if err != nil {
		return
	}

	ledgerTran := &LedgerTransaction{
		SignedTransaction: signedTran,
		ServerTimestamp: bData.ServerTimestamp,
		Status: ValidTransactionStatus,
	}
	l.addTran(ledgerTran)
}

func (l *Ledger) GetSpendAddressState(address string) (*SpendChannelState, error) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if !helpers.IsAddressSpend(address) {
		return nil, errors.New(fmt.Sprintf("ledger: GetSpendAddressState: invalid spend address: %v", address))
	}

	pos, err := l.getOpenSpendChannelPositionByAddress(address)

	if err != nil {
		return nil, err
	}

	if pos > -1 {
		return l.getSpendChannelState(pos)
	}

	return nil, errors.New(fmt.Sprintf("ledger: GetSpendAddressState: spend address no find: %v", address))
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
) (*TransferChannelState, error) {
	tran := blockchain.NewOpenTransferTran(
		&trans.OpenTransferTran{
			ChannelId: channelId,
			FromAddress: fromAddress,
			FromPk: fromPk,
			ToAddress: toAddress,
			ToPk: toPk,
			PriceAmount: uint64(priceAmount),
			PriceQuantumPower: uint8(priceQuantumPower),
			PlannedQuantumCount: plannedQuantumCount,
			LifeTime: l.lifeTimes.TransferChannel,
		},
		l.pk,
	)

	var err error

	ch := make(chan *LedgerTransaction)
	defer func() {
		close(ch)
	}()

	err = l.waiters.Add(tran.Id, ch)
	if err != nil {
		return nil, err
	}

	err = l.sendTran(tran)
	if err != nil {
		return nil, err
	}

	ledgerTran := <-ch
	if ledgerTran == nil {
		return nil, errors.New("ledger: transaction not synced")
	}

	var data *TransferChannelState = nil

	if ledgerTran.Data != nil {
		data = ledgerTran.Data.(*TransferChannelState)
	}

	return data, nil
}

func (l *Ledger) openTransferChannel(bData *BlockchainTran, tData *trans.OpenTransferTran) {
	if helpers.IsAddressCommon(tData.ToAddress) {
		l.openTransferChannelWithCommonAddress(bData, tData)
	} else {
		l.openTransferChannelWithIncomeAddress(bData, tData)
	}
}

func (l *Ledger) GetTransferChannelState(
	channelId string,
) (TransferChannelStatus, error) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	for i := len(l.transactions) - 1; i >= 0; i-- {
		tran := l.transactions[i]
		if !tran.IsValid() {
			continue
		}
		if tran.ContractId == channelId {
			switch tran.Type {
			case OpenTransferChannelTransactionType:
				return OpenTransferChannelStatus, nil
			case CloseTransferChannelTransactionType:
				return CloseTransferChannelStatus, nil
			default:
				return ErrorTransferChannelStatus, errors.New("ledger: GetTransferChannelState: inconsistent tran log or invalid channel id")
			}
		}
	}

	return NotExistTransferChannelStatus, nil
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
			ChannelId: channelId,
			FromAddress: fromAddress,
			FromPk: fromPk,
			ToAddress: toAddress,
			ToPk: toPk,
			PriceAmount: uint64(priceAmount),
			PriceQuantumPower: uint8(priceQuantumPower),
			PlannedQuantumCount: plannedQuantumCount,
			Amount: uint64(amount),
			QuantumCount: quantumCount,
			Volume: volume,
		},
		l.pk,
	)


	var err error

	ch := make(chan *LedgerTransaction)
	defer func() {
		close(ch)
	}()

	err = l.waiters.Add(tran.Id, ch)
	if err != nil {
		return InvalidTransactionStatus, err
	}

	err = l.sendTran(tran)
	if err != nil {
		return InvalidTransactionStatus, err
	}

	ledgerTran := <-ch
	if ledgerTran == nil {
		return InvalidTransactionStatus, errors.New("ledger: transaction not synced")
	}

	return ledgerTran.Status, nil
}

func (l *Ledger) closeTransferChannel(
	bData *BlockchainTran, tData *trans.CloseTransferTran,
) {
	if helpers.IsAddressCommon(tData.ToAddress) {
		l.closeTransferChannelWithCommonAddress(bData, tData)
	} else {
		l.closeTransferChannelWithIncomeAddress(bData, tData)
	}
}

func (l *Ledger) CloseIncomeChannel(address string, pk string) (TransactionStatus, *ChannelFact, error) {
	errResult := func(err error) (TransactionStatus, *ChannelFact, error) {
		common.Log.Print(common.Printf(
			"CloseIncomeChannel: %v",
			err.Error(),
		))
		return InvalidTransactionStatus, nil, err
	}

	tran := blockchain.NewCloseIncomeTran(
		&trans.CloseIncomeTran{
			From: address,
			To: l.address,
			PeerPublicKey: pk,
		},
		l.pk,
	)

	ch := make(chan *LedgerTransaction)
	defer func() {
		close(ch)
	}()

	var err error
	err = l.waiters.Add(tran.Id, ch)
	if err != nil {
		return errResult(err)
	}

	err = l.sendTran(tran)
	if err != nil {
		return errResult(err)
	}

	ledgerTran := <-ch
	if ledgerTran == nil {
		return errResult(errors.New("transaction not synced"))
	}

	var data *ChannelFact = nil
	if ledgerTran.Data != nil {
		d := ledgerTran.Data.(*IncomeChannelState)
		data = d.State.Fact.ToChannel()
	}

	return ledgerTran.Status, data, nil
}

func (l *Ledger) closeIncomeChannel(bData *BlockchainTran, tData *trans.CloseIncomeTran) {
	var err error
	success := false
	defer func () {
		if !success {
			common.Log.Print(common.Printf(
				"closeIncomeChannel: %v",
				err.Error(),
			))
			l.waiters.TryInterrupt(bData.Id)
		}
	}()

	if !helpers.IsAddressIncome(tData.From) {
		err = errors.New("invalid income address")
		return
	}

	pos, err := l.getOpenIncomeChannelPositionByAddress(tData.From)
	if err != nil {
		return
	}
	if pos < 0 {
		err = errors.New("income address not found")
		return
	}

	tran := l.transactions[pos]
	if !l.validateIncomeAddress(tData.From, tData.PeerPublicKey, tran) {
		err = errors.New("income address has different operator")
		return
	}
	if tran.TimeLock < bData.ClientTimestamp {
		err = errors.New("income channel has already expired")
		return
	}

	hasOpenedChannels, err := l.hasOpenedTransferChannelsInIncomeChannel(pos)
	if err != nil {
		return
	}
	if hasOpenedChannels {
		err =  errors.New("income address has opened channels")
		return
	}

	state, err := l.getIncomeChannelState(pos)
	if err != nil {
		return
	}

	signedTran, err := l.createCloseIncomeChannelTran(
		tran.Id,
		tData.From,
		tran.From,
		state.State.Fact.Amount,
		state.State.Fact.Volume,
		bData.ClientTimestamp,
		bData.Id,
		bData.PublicKey,
		state,
	)
	if err != nil {
		return
	}

	ledgerTran := &LedgerTransaction{
		SignedTransaction: signedTran,
		ServerTimestamp: bData.ServerTimestamp,
		Status: ValidTransactionStatus,
	}

	l.addTran(ledgerTran)
	success = true
}

func (l *Ledger) CloseSpendChannel(address string, pk string) (TransactionStatus, *ChannelFact, error) {
	errResult := func(err error) (TransactionStatus, *ChannelFact, error) {
		common.Log.Print(common.Printf(
			"CloseSpendChannel: %v",
			err.Error(),
		))
		return InvalidTransactionStatus, nil, err
	}
	tran := blockchain.NewCloseSpendTran(
		&trans.CloseSpendTran{
			PeerPublicKey: pk,
			From: address,
			To: l.address,
		},
		l.pk,
	)

	ch := make(chan *LedgerTransaction)
	defer func() {
		close(ch)
	}()

	err := l.waiters.Add(tran.Id, ch)
	if err != nil {
		return errResult(err)
	}

	err = l.sendTran(tran)
	if err != nil {
		return errResult(err)
	}

	ledgerTran := <-ch
	if ledgerTran == nil {
		return errResult(errors.New("transaction not synced"))
	}

	var data *ChannelFact
	if ledgerTran.Data != nil {
		d := ledgerTran.Data.(*SpendChannelState)
		data = d.State.Fact.ToChannel()
	}

	return ledgerTran.Status, data, nil
}

func (l *Ledger) closeSpendChannel(bData *BlockchainTran, tData *trans.CloseSpendTran) {
	var err error
	success := false
	defer func () {
		if !success {
			common.Log.Print(common.Printf(
				"closeSpendChannel: %v",
				err.Error(),
			))
			l.waiters.TryInterrupt(bData.Id)
		}
	}()

	if !helpers.IsAddressSpend(tData.From) {
		err = errors.New("invalid spend address")
		return
	}

	pos, err := l.getOpenSpendChannelPositionByAddress(tData.From)
	if err != nil {
		return
	}
	if pos < 0 {
		err = errors.New("spend address not found")
		return
	}

	tran := l.transactions[pos]
	if !l.validateSpendAddress(tData.From, tData.PeerPublicKey, tran) {
		err = errors.New("spend address has different operator")
		return
	}
	if tran.TimeLock < bData.ClientTimestamp {
		err = errors.New("spend channel has already expired")
		return
	}

	hasOpenedChannels, err := l.hasOpenedTransferChannelsInSpendChannel(pos)
	if err != nil {
		return
	}
	if hasOpenedChannels {
		err = errors.New("spend address has openned transfer channels")
		return
	}

	state, err := l.getSpendChannelState(pos)
	if err != nil {
		return
	}
	signedTran, err := l.createCloseSpendChannelTran(
		tran.Id,
		tData.From,
		tran.From,
		state.State.Fact.Amount,
		state.State.Fact.Volume,
		bData.ClientTimestamp,
		bData.Id,
		bData.PublicKey,
		state,
	)
	if err != nil {
		return
	}

	ledgerTran := &LedgerTransaction{
		SignedTransaction: signedTran,
		ServerTimestamp: bData.ServerTimestamp,
		Status: ValidTransactionStatus,
	}

	l.addTran(ledgerTran)
	success = true
}

func (l *Ledger) GetCommonAddressBalance(address string) (*CommonAddressBalanceState, error) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if !helpers.IsAddressCommon(address) {
		return nil, errors.New(fmt.Sprintf("ledger: GetCommonAddressBalance: invalid common address: %v", address))
	}

	balance := NewCommonAddressBalance()

	for i := 0; i < len(l.transactions); i++ {
		tran := l.transactions[i]

		if !tran.IsValid() {
			continue
		}
		if tran.From != address && tran.To != address {
			continue
		}

		switch tran.Type {
		case InitTransactionType:
			if tran.To == address {
				balance.Amount += tran.Amount
			}
			break
		case OpenSpendChannelTransactionType:
			if tran.From == address {
				balance.Amount -= tran.Amount
			}
			break
		case OpenIncomeChannelTransactionType:
			if tran.From == address {
				balance.Amount -= tran.Amount
			}
			break
		case CloseSpendChannelTransactionType:
			if tran.To == address {
				balance.Amount += tran.Amount
			}
			break
		case CloseIncomeChannelTransactionType:
			if tran.To == address {
				balance.Amount += tran.Amount
			}
			break
		case OpenTransferChannelTransactionType:
			if tran.To == address {
				planEarn := tran.PriceAmount * TransactionAmountType(tran.PlannedQuantumCount)
				if tran.PublicKey == l.pk {
					balance.OwnEarn.Plan += planEarn
				} else {
					balance.ForeignEarn.Plan += planEarn
				}
			} else {
				return nil, errors.New("ledger: GetCommonAddressBalance: inconsistent tran log: invalid OpenTransferChannelTransactionType")
			}
			break
		case CloseTransferChannelTransactionType:
			if tran.To == address {
				planEarn := tran.PriceAmount * TransactionAmountType(tran.PlannedQuantumCount)
				factEarn := tran.Amount
				if tran.PublicKey == l.pk {
					balance.OwnEarn.Plan -= planEarn
					balance.OwnEarn.Fact += factEarn
				} else {
					balance.ForeignEarn.Plan -= planEarn
					balance.ForeignEarn.Fact += factEarn
				}
			} else {
				return nil, errors.New("ledger: GetCommonAddressBalance: inconsistent tran log: invalid CloseTransferChannelTransactionType")
			}
			break
		case FinancialBlockTransactionType:
			if tran.From == address {
				balance.Profit += balance.OwnEarn.Fact
				balance.OwnEarn.Fact = 0
			}
			break
		default:
			return nil, errors.New("ledger: GetCommonAddressBalance: inconsistent tran log")
		}
	}

	return balance.ToState(), nil
}

func (l *Ledger) addTran(tran *LedgerTransaction) {
	if tran.Status == ValidTransactionStatus {
		l.transactions = append(l.transactions, tran)
		l.financialState.Add(tran)
	}
	l.waiters.TryClose(tran)
}

func (l *Ledger) sendTran(tran *blockchain.Tran) error {
	return l.writer.Write(tran)
}

func (l *Ledger) openTransferChannelWithCommonAddress(
	bData *BlockchainTran, tData *trans.OpenTransferTran,
) {
	var err error
	success := false
	defer func() {
		if !success {
			common.Log.Print(common.Printf(
				"openTransferChannelWithCommonAddress: %v",
				err.Error(),
			))
			l.waiters.TryInterrupt(bData.Id)
		}
	}()
	status, tran, err := l.preOpenTransferChannelWithCommonAddress(tData)

	switch status {
	case InvalidTransactionStatus:
		return
	case NoneTransactionStatus:
		l.waiters.TryCloseAnother(bData.Id, tran)
		success = true
		return
	}

	timestamp := bData.ClientTimestamp
	timeLock := timestamp + tData.LifeTime

	spendPos, _, err := l.getOwnOpenSpendChannelPositionByAddress(tData.FromAddress, tData.FromPk)
	if err != nil {
		return
	}
	if spendPos == -1 {
		err = errors.New("spend address don't find")
		return
	}

	spendState, err := l.getSpendChannelState(spendPos)
	if err != nil {
		return
	}
	if spendState.TimeLock < timestamp {
		err = errors.New("spend channel has already expired")
		return
	}

	price := NewPrice(TransactionAmountType(tData.PriceAmount), QuantumPowerType(tData.PriceQuantumPower))
	if res, err := price.Cmp(spendState.MaxPrice()); err != nil || res == 1 {
		err = errors.New("invalid spend traffic price")
		return
	}

	plannedAmount := TransactionAmountType(tData.PriceAmount) * TransactionAmountType(tData.PlannedQuantumCount)
	if spendState.State.Plan.Amount < plannedAmount {
		err = errors.New("spend address hasn't got enough amount to download traffic")
		return
	}

	if timeLock > spendState.TimeLock {
		timeLock = spendState.TimeLock
	}

	signedTran, err := l.createOpenTransferChannelTran(
		tData.ChannelId,
		tData.FromAddress,
		tData.ToAddress,
		TransactionAmountType(tData.PriceAmount),
		QuantumPowerType(tData.PriceQuantumPower),
		tData.PlannedQuantumCount,
		timestamp,
		timeLock,
		bData.Id,
		bData.PublicKey,
		NewTransferChannelState(timeLock, timeLock - timestamp),
	)
	if err != nil {
		return
	}

	ledgerTran := &LedgerTransaction{
		SignedTransaction: signedTran,
		Status: ValidTransactionStatus,
		ServerTimestamp: bData.ServerTimestamp,
	}

	l.addTran(ledgerTran)
	success = true
}

func (l *Ledger) preOpenTransferChannelWithCommonAddress(
	tData *trans.OpenTransferTran,
) (TransactionStatus, *LedgerTransaction, error) {
	if tData.PlannedQuantumCount == 0 {
		return InvalidTransactionStatus, nil, errors.New("ledger: preOpenTransferChannelWithCommonAddress: transfer zero quantum is invalid")
	}

	if !helpers.IsAddressSpend(tData.FromAddress) {
		return InvalidTransactionStatus, nil, errors.New("ledger: preOpenTransferChannelWithCommonAddress: invalid spend address")
	}
	if !helpers.IsAddressCommon(tData.ToAddress) {
		return InvalidTransactionStatus, nil, errors.New("ledger: preOpenTransferChannelWithCommonAddress: invalid common income address")
	}

	tran, err := l.getOpenTransferChannelTran(
		tData.ChannelId,
		tData.FromAddress,
		tData.ToAddress,
		TransactionAmountType(tData.PriceAmount),
		QuantumPowerType(tData.PriceQuantumPower),
		tData.PlannedQuantumCount,
	)
	if err != nil {
		return InvalidTransactionStatus, nil, err
	}
	if tran != nil {
		return NoneTransactionStatus, tran, nil
	}

	_, _, err = l.getOwnOpenSpendChannelPositionByAddress(tData.FromAddress, tData.FromPk)
	if err != nil {
		return InvalidTransactionStatus, nil, err
	}

	return ValidTransactionStatus, nil, nil
}

func (l *Ledger) openTransferChannelWithIncomeAddress(
	bData *BlockchainTran, tData *trans.OpenTransferTran,
) {
	var err error
	success := false
	defer func() {
		if !success {
			common.Log.Print(common.Printf(
				"openTransferChannelWithIncomeAddress: %v",
				err.Error(),
			))
			l.waiters.TryInterrupt(bData.Id)
		}
	}()

	status, tran, err := l.preOpenTransferChannelWithIncomeAddress(tData)

	timestamp := bData.ClientTimestamp
	timeLock := timestamp + tData.LifeTime

	switch status {
	case InvalidTransactionStatus:
		return
	case NoneTransactionStatus:
		l.waiters.TryCloseAnother(bData.Id, tran)
		success = true
		return
	}

	incomePos, _, err := l.getOwnOpenIncomeChannelPositionByAddress(tData.ToAddress, tData.ToPk)
	if err != nil {
		return
	}
	if incomePos == -1 {
		err = errors.New("income address don't find")
		return
	}

	spendPos, _, err := l.getOwnOpenSpendChannelPositionByAddress(tData.FromAddress, tData.FromPk)
	if err != nil {
		return
	}
	if spendPos == -1 {
		err = errors.New("spend address don't find")
		return
	}

	spendState, err := l.getSpendChannelState(spendPos)
	if err != nil {
		return
	}
	if spendState.TimeLock < timestamp {
		err = errors.New("spend channel has already expired")
		return
	}

	incomeState, err := l.getIncomeChannelState(incomePos)
	if err != nil {
		return
	}
	if incomeState.TimeLock < timestamp {
		err = errors.New("income channel has already expired")
		return
	}

	price := NewPrice(TransactionAmountType(tData.PriceAmount), QuantumPowerType(tData.PriceQuantumPower))
	if res, err := price.Cmp(spendState.MaxPrice()); err != nil || res == 1 {
		err = errors.New("invalid spend traffic price")
		return
	}
	if res, err := price.Cmp(incomeState.MinPrice()); err != nil || res == -1 {
		err = errors.New("invalid income traffic price")
		return
	}

	plannedAmount := TransactionAmountType(tData.PriceAmount) * TransactionAmountType(tData.PlannedQuantumCount)
	if spendState.State.Plan.Amount < plannedAmount {
		err = errors.New("spend address hasn't got enough amount to download traffic")
		return
	}

	if timeLock > spendState.TimeLock {
		timeLock = spendState.TimeLock
	}
	if timeLock > incomeState.TimeLock {
		timeLock = incomeState.TimeLock
	}

	signedTran, err := l.createOpenTransferChannelTran(
		tData.ChannelId,
		tData.FromAddress,
		tData.ToAddress,
		TransactionAmountType(tData.PriceAmount),
		QuantumPowerType(tData.PriceQuantumPower),
		tData.PlannedQuantumCount,
		timestamp,
		timeLock,
		bData.Id,
		bData.PublicKey,
		NewTransferChannelState(timeLock, timeLock - timestamp),
	)
	if err != nil {
		return
	}

	ledgerTran := &LedgerTransaction{
		SignedTransaction: signedTran,
		ServerTimestamp: bData.ServerTimestamp,
		Status: ValidTransactionStatus,
	}
	l.addTran(ledgerTran)
	success = true
}

func (l *Ledger) preOpenTransferChannelWithIncomeAddress(
	tData *trans.OpenTransferTran,
) (TransactionStatus, *LedgerTransaction, error) {
	if tData.PlannedQuantumCount == 0 {
		return InvalidTransactionStatus, nil, errors.New("transfer zero quantum is invalid")
	}
	if !helpers.IsAddressSpend(tData.FromAddress) {
		return InvalidTransactionStatus, nil, errors.New("invalid spend address")
	}
	if !helpers.IsAddressIncome(tData.ToAddress) {
		return InvalidTransactionStatus, nil, errors.New("invalid income address")
	}

	tran, err := l.getOpenTransferChannelTran(
		tData.ChannelId,
		tData.FromAddress,
		tData.ToAddress,
		TransactionAmountType(tData.PriceAmount),
		QuantumPowerType(tData.PriceQuantumPower),
		tData.PlannedQuantumCount,
	)
	if err != nil {
		return InvalidTransactionStatus, nil, err
	}
	if tran != nil {
		return NoneTransactionStatus, tran, nil
	}

	_, _, err = l.getOwnOpenIncomeChannelPositionByAddress(tData.ToAddress, tData.ToPk)
	if err != nil {
		return InvalidTransactionStatus, nil, err
	}

	_, _, err = l.getOwnOpenSpendChannelPositionByAddress(tData.FromAddress, tData.FromPk)
	if err != nil {
		return InvalidTransactionStatus, nil, err
	}

	return ValidTransactionStatus, nil, nil
}

func (l *Ledger) closeTransferChannelWithCommonAddress(
	bData *BlockchainTran, tData *trans.CloseTransferTran,
) {
	var err error
	success := false
	defer func() {
		if !success {
			if err != nil {
				common.Log.Print(common.Printf(
					"closeTransferChannelWithCommonAddress: %v",
					err.Error(),
				))
			}
			l.waiters.TryInterrupt(bData.Id)
		}
	}()
	status, openTran, closeTran, err := l.preCloseTransferChannelWithCommonAddress(tData)
	switch status {
	case InvalidTransactionStatus:
		return
	case NoneTransactionStatus:
		l.waiters.TryCloseAnother(bData.Id, closeTran)
		success = true
		return
	}

	timestamp := bData.ClientTimestamp
	if openTran.TimeLock < timestamp {
		err = errors.New("transfer channel has already expired")
		return
	}

	spendPos, _, err := l.getOwnOpenSpendChannelPositionByAddress(tData.FromAddress, tData.FromPk)
	if err != nil {
		return
	}
	if spendPos == -1 {
		err = errors.New("spend address did't find")
		return
	}

	spendState, err := l.getSpendChannelState(spendPos)
	if err != nil {
		return
	}
	if spendState.TimeLock < timestamp {
		err = errors.New("spend channel has already expired")
		return
	}

	signedTran, err := l.createCloseTransferChannelTran(
		openTran.Id,
		tData.ChannelId,
		tData.FromAddress,
		tData.ToAddress,
		TransactionAmountType(tData.PriceAmount),
		QuantumPowerType(tData.PriceQuantumPower),
		tData.PlannedQuantumCount,
		TransactionAmountType(tData.Amount),
		tData.QuantumCount,
		tData.Volume,
		timestamp,
		bData.Id,
		bData.PublicKey,
	)

	ledgerTran := &LedgerTransaction{
		SignedTransaction: signedTran,
		Status: ValidTransactionStatus,
		ServerTimestamp: bData.ServerTimestamp,
	}

	l.addTran(ledgerTran)
	success = true
}

func (l *Ledger) preCloseTransferChannelWithCommonAddress(
	tData *trans.CloseTransferTran,
) (TransactionStatus, *LedgerTransaction, *LedgerTransaction, error) {
	if tData.QuantumCount > tData.PlannedQuantumCount {
		return InvalidTransactionStatus, nil, nil, errors.New("preCloseTransferChannelWithCommonAddress: invalid quantum count: greater than planneng quantum count")
	}

	quantumVolume := tData.QuantumCount * QuantumPowerType(tData.PriceQuantumPower).Volume()
	if quantumVolume < tData.Volume {
		return InvalidTransactionStatus, nil, nil, errors.New("preCloseTransferChannelWithCommonAddress: invalid quantum count: volume are greater than paid")
	}

	if TransactionAmountType(tData.Amount) != TransactionAmountType(tData.PriceAmount) * TransactionAmountType(tData.QuantumCount) {
		return InvalidTransactionStatus, nil, nil, errors.New("preCloseTransferChannelWithCommonAddress: invalid amount: calculated to another value")
	}

	if !helpers.IsAddressSpend(tData.FromAddress) {
		return InvalidTransactionStatus, nil, nil, errors.New("preCloseTransferChannelWithCommonAddress: invalid spend address")
	}

	if !helpers.IsAddressCommon(tData.ToAddress) {
		return InvalidTransactionStatus, nil, nil, errors.New("preCloseTransferChannelWithCommonAddress: invalid common address")
	}

	openTran, err := l.getOpenTransferChannelTran(
		tData.ChannelId,
		tData.FromAddress,
		tData.ToAddress,
		TransactionAmountType(tData.PriceAmount),
		QuantumPowerType(tData.PriceQuantumPower),
		tData.PlannedQuantumCount,
	)
	if err != nil {
		return InvalidTransactionStatus, nil, nil, err
	}
	if openTran == nil {
		return InvalidTransactionStatus, nil, nil, errors.New("preCloseTransferChannelWithCommonAddress: open transfer channel tran not found")
	}

	closeTran, err := l.getCloseTransferChannelTran(
		openTran.Id,
		tData.ChannelId,
		tData.FromAddress,
		tData.ToAddress,
		TransactionAmountType(tData.PriceAmount),
		QuantumPowerType(tData.PriceQuantumPower),
		tData.PlannedQuantumCount,
		TransactionAmountType(tData.Amount),
		tData.QuantumCount,
		tData.Volume,
	)
	if err != nil {
		return InvalidTransactionStatus, openTran, nil, err
	}
	if closeTran != nil {
		return NoneTransactionStatus, openTran, closeTran, nil
	}

	return ValidTransactionStatus, openTran, nil, nil
}

func (l *Ledger) closeTransferChannelWithIncomeAddress(
	bData *BlockchainTran, tData *trans.CloseTransferTran,
) {
	var err error
	success := false
	defer func() {
		if !success {
			if err != nil {
				common.Log.Print(common.Printf(
					"closeTransferChannelWithIncomeAddress: %v",
					err.Error(),
				))
			}
			l.waiters.TryInterrupt(bData.Id)
		}
	}()
	status, openTran, closeTran, err := l.preCloseTransferChannelWithIncomeAddress(tData)

	switch status {
	case InvalidTransactionStatus:
		return
	case NoneTransactionStatus:
		l.waiters.TryCloseAnother(bData.Id,closeTran)
		success = true
		return
	}

	timestamp := bData.ClientTimestamp

	if openTran.TimeLock < timestamp {
		err = errors.New("transfer channel has already expired")
		return
	}

	incomePos, _, err := l.getOwnOpenIncomeChannelPositionByAddress(tData.ToAddress, tData.ToPk)
	if err != nil {
		return
	}
	if incomePos == -1 {
		err = errors.New("income address didn't find")
		return
	}

	spendPos, _, err := l.getOwnOpenSpendChannelPositionByAddress(tData.FromAddress, tData.FromPk)
	if err != nil {
		return
	}
	if spendPos == -1 {
		err = errors.New("spend address didn't find")
		return
	}

	incomeState, err := l.getIncomeChannelState(incomePos)
	if err != nil {
		return
	}
	if incomeState.TimeLock < timestamp  {
		err = errors.New("income channel has already expired")
		return
	}

	spendState, err := l.getSpendChannelState(spendPos)
	if err != nil {
		return
	}
	if spendState.TimeLock < timestamp {
		err = errors.New("spend channel has already expired")
		return
	}

	signedTran, err := l.createCloseTransferChannelTran(
		openTran.Id,
		tData.ChannelId,
		tData.FromAddress,
		tData.ToAddress,
		TransactionAmountType(tData.PriceAmount),
		QuantumPowerType(tData.PriceQuantumPower),
		tData.PlannedQuantumCount,
		TransactionAmountType(tData.Amount),
		tData.QuantumCount,
		tData.Volume,
		timestamp,
		bData.Id,
		bData.PublicKey,
	)
	if err != nil {
		return
	}

	ledgerTran := &LedgerTransaction{
		SignedTransaction: signedTran,
		ServerTimestamp: bData.ServerTimestamp,
		Status: ValidTransactionStatus,
	}
	l.addTran(ledgerTran)
	success = true
}

func (l *Ledger) preCloseTransferChannelWithIncomeAddress(
	tData *trans.CloseTransferTran,
) (TransactionStatus, *LedgerTransaction, *LedgerTransaction, error) {
	// Проверка: количество переданных квантов не больше запланированного
	if tData.QuantumCount > tData.PlannedQuantumCount{
		return InvalidTransactionStatus, nil, nil, errors.New("preCloseTransferChannelWithIncomeAddress: invalid quantum count: greater than planning quantum count")
	}

	// Объём траффика к оплате
	quantumVolume := tData.QuantumCount * QuantumPowerType(tData.PriceQuantumPower).Volume()

	// Проверка фактичкий объем переданного траффика не превышает объём траффика к оплате
	if quantumVolume < tData.Volume {
		return InvalidTransactionStatus, nil, nil, errors.New("preCloseTransferChannelWithIncomeAddress: invalid quantum count: volume are greater than paid")
	}

	// Переводимое число токенов совпадает с вычисленным числом токенов для перевода
	if TransactionAmountType(tData.Amount) != TransactionAmountType(tData.PriceAmount) * TransactionAmountType(tData.QuantumCount) {
		return InvalidTransactionStatus, nil, nil, errors.New("preCloseTransferChannelWithIncomeAddress: invalid amount: calculated to another value")
	}

	if !helpers.IsAddressSpend(tData.FromAddress) {
		return InvalidTransactionStatus, nil, nil, errors.New("preCloseTransferChannelWithIncomeAddress: invalid spend address")
	}

	if !helpers.IsAddressIncome(tData.ToAddress) {
		return InvalidTransactionStatus, nil, nil, errors.New("preCloseTransferChannelWithIncomeAddress: invalid income address")
	}

	openTran, err := l.getOpenTransferChannelTran(
		tData.ChannelId,
		tData.FromAddress,
		tData.ToAddress,
		TransactionAmountType(tData.PriceAmount),
		QuantumPowerType(tData.PriceQuantumPower),
		tData.PlannedQuantumCount,
	)
	if err != nil {
		return InvalidTransactionStatus, nil, nil, err
	}
	if openTran == nil {
		return InvalidTransactionStatus, nil, nil, errors.New("preCloseTransferChannelWithIncomeAddress: open transfer channel tran not found")
	}

	closeTran, err := l.getCloseTransferChannelTran(
		openTran.Id,
		tData.ChannelId,
		tData.FromAddress,
		tData.ToAddress,
		TransactionAmountType(tData.PriceAmount),
		QuantumPowerType(tData.PriceQuantumPower),
		tData.PlannedQuantumCount,
		TransactionAmountType(tData.Amount),
		tData.QuantumCount,
		tData.Volume,
	)
	if err != nil {
		return InvalidTransactionStatus, openTran, nil, err
	}
	if closeTran != nil {
		return NoneTransactionStatus, openTran, closeTran, nil
	}

	return ValidTransactionStatus, openTran, nil, nil
}

func (l *Ledger) getCopiedTrans(fields *GetTransactionsFields) *LedgerTransactions {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	var copied LedgerTransactions

	if fields != nil {
		copied = make(LedgerTransactions, 0, len(l.transactions))

		for _, item := range l.transactions {
			if fields.Type != nil {
				if item.Type != fields.Type.Value {
					continue
				}
			}

			if fields.To != nil {
				if item.To != fields.To.Value {
					continue
				}
			}

			if fields.From != nil {
				if item.From != fields.From.Value {
					continue
				}
			}

			if fields.Status != nil {
				if item.Status != fields.Status.Value {
					continue
				}
			}

			if fields.ChannelId != nil {
				if item.ContractId != fields.ChannelId.Value {
					continue
				}
			}

			copied = append(copied, item)
		}
	} else {
		copied = make(LedgerTransactions, len(l.transactions))
		copy(copied, l.transactions)
	}

	sort.Slice(copied, func(i int, j int) bool {
		return copied[i].TimeStamp < copied[j].TimeStamp
	})

	return &copied
}

func (l *Ledger) sort(trans *LedgerTransactions, order OrderType) {
	if order == NoneOrder {
		return
	}

	arr := *trans
	sort.Slice(arr, func(i int, j int) bool {
		if order == DescOrder {
			return arr[i].TimeStamp > arr[j].TimeStamp
		} else {
			return arr[i].TimeStamp < arr[j].TimeStamp
		}
	})
}

func (l *Ledger) getDataFromCopied(
	copied *LedgerTransactions,
	pageSize int, page uint,
) (*LedgerTransactions, uint) {
	arr := *copied
	totalCount := uint(len(arr))
	var from uint
	var to uint

	if pageSize < 0 {
		from = 0
		to = totalCount
	} else {
		from = (page - 1) * uint(pageSize)
		to = from + uint(pageSize)
	}

	var data LedgerTransactions

	if from >= totalCount {
		data = make(LedgerTransactions, 0, 0)
	} else {
		if to > totalCount {
			to = totalCount
		}

		data = arr[from:to]
	}

	return &data, totalCount
}

func (l *Ledger) createPayload(data interface{}) (*json.RawMessage, error) {
	reqData, err := helpers.NewResponseDataInterface(data, l.keyPair)

	if err != nil {
		return nil, err
	}

	signedRequestData := &helpers.SignedRequestData{
		PublicKey: l.pk,
		Message:   reqData.Message,
		Signature: base64.StdEncoding.EncodeToString(reqData.Signature),
	}

	var res json.RawMessage
	res, err = json.Marshal(signedRequestData)

	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (l *Ledger) createTran(
	body *TransactionBody,
	id string,
	pk string,
) (*SignedTransaction, error) {
	signedTran, err := NewSignedTransaction(body, l.keyPair, id, pk)
	if err != nil {
		return nil, err
	}

	return signedTran, nil
}

func (l *Ledger) hasValidTrans() bool {
	// Проверяем, что для данного кошелька не было начального начисления
	// Упрощение относительно оригинального кода
	for i := 0; i < len(l.transactions); i++ {
		if l.transactions[i].IsValid() &&
			l.transactions[i].Type == InitTransactionType &&
			l.transactions[i].To == l.address {
			return true
		}
	}

	return false
}

func (l *Ledger) createOpenTransferChannelTran(
	channelId string,
	fromAddress string,
	toAddress string,
	priceAmount TransactionAmountType,
	priceQuantumPower QuantumPowerType,
	plannedQuantumCount uint64,
	timestamp int64,
	timeLock int64,
	id string,
	pk string,
	state *TransferChannelState,

) (*SignedTransaction, error) {
	body := &TransactionBody{
		TimeStamp: timestamp,
		TimeLock: timeLock,
		Type: OpenTransferChannelTransactionType,
		From: fromAddress,
		To: toAddress,
		ContractId: channelId,
		PriceAmount: priceAmount,
		PriceQuantumPower: priceQuantumPower,
		PlannedQuantumCount: plannedQuantumCount,
		Data: state,
	}

	return l.createTran(body, id, pk)
}

func (l *Ledger) createCloseTransferChannelTran(
	parentId string,
	channelId string,
	fromAddress string,
	toAddress string,
	priceAmount TransactionAmountType,
	priceQuantumPower QuantumPowerType,
	plannedQuantumCount uint64,
	amount TransactionAmountType,
	quantumCount uint64,
	volume uint64,
	timestamp int64,
	id string,
	pk string,
) (*SignedTransaction, error){
	body := &TransactionBody{
		TimeStamp: timestamp,
		Type: CloseTransferChannelTransactionType,
		ParentId: parentId,
		From: fromAddress,
		To: toAddress,
		ContractId: channelId,
		PriceAmount: priceAmount,
		PriceQuantumPower: priceQuantumPower,
		PlannedQuantumCount: plannedQuantumCount,
		Amount: amount,
		QuantumCount: quantumCount,
		Volume: volume,
	}

	return l.createTran(body, id, pk)
}

func (l *Ledger) createCloseIncomeChannelTran(
	parentId string,
	fromAddress string,
	toAddress string,
	amount TransactionAmountType,
	volume uint64,
	timestamp int64,
	id string,
	publicKey string,
	state *IncomeChannelState,
) (*SignedTransaction, error) {
	body := &TransactionBody{
		TimeStamp: timestamp,
		Type: CloseIncomeChannelTransactionType,
		ParentId: parentId,
		From: fromAddress,
		To: toAddress,
		Amount: amount,
		Volume: volume,
		Data: state,
	}

	return l.createTran(body, id, publicKey)
}

func (l *Ledger) createCloseSpendChannelTran(
	parentId string,
	fromAddress string,
	toAddress string,
	amount TransactionAmountType,
	volume uint64,
	timestamp int64,
	id string,
	pk string,
	state *SpendChannelState,
) (*SignedTransaction, error) {
	body := &TransactionBody{
		TimeStamp: timestamp,
		Type: CloseSpendChannelTransactionType,
		ParentId: parentId,
		From: fromAddress,
		To: toAddress,
		Amount: amount,
		Volume: volume,
		Data: state,
	}

	return l.createTran(body, id, pk)
}

// адрес должен быть общим адресом (common address)
// функция не валидилует входящий адрес
func (l *Ledger) getAddressAmount(address string) (TransactionAmountType, error) {
	var amount TransactionAmountType = 0
	for i := 0; i < len(l.transactions); i++ {
		tran := l.transactions[i]

		if !tran.IsValid() {
			continue
		}
		if tran.To != address && tran.From != address {
			continue
		}

		switch tran.Type {
		case InitTransactionType:
			if tran.To == address {
				amount += tran.Amount
			}
			break
		case OpenSpendChannelTransactionType:
			if tran.From == address {
				amount -= tran.Amount
			}
			break
		case OpenIncomeChannelTransactionType:
			if tran.From == address {
				amount -= tran.Amount
			}
			break
		case CloseSpendChannelTransactionType:
			if tran.To == address {
				amount += tran.Amount
			}
			break
		case CloseIncomeChannelTransactionType:
			if tran.To == address {
				amount += tran.Amount
			}
			break
		case OpenTransferChannelTransactionType:
			if tran.From != address {
				return TransactionAmountType(0), errors.New("ledger: getAddressAmount: inconsistent tran log: invalid OpenTransferChannelTransactionType")
			}
			break
		case CloseTransferChannelTransactionType:
			if tran.From != address {
				return TransactionAmountType(0), errors.New("ledger: getAddressAmount: inconsistent tran log: invalid CloseTransferChannelTransactionType")
			}
			break
		case FinancialBlockTransactionType:
			break
		default:
			tranJsonBytes, _ := json.Marshal(tran)
			return TransactionAmountType(0), errors.New(fmt.Sprintf("ledger: getAddressAmount: inconsistent tran log: %v", string(tranJsonBytes[:])))
		}
	}

	return amount, nil
}

// Позиция транзакции в журнале
// Если не найдено: -1
func (l *Ledger) getOpenIncomeChannelPositionByPublicKey(peerPublicKey string) (int, error) {
	closedAddresses := make(map[string]string)

	for i := len(l.transactions) - 1; i >= 0; i-- {
		tran := l.transactions[i]

		if !tran.IsValid() {
			continue
		}

		switch tran.Type {
		case CloseIncomeChannelTransactionType:
			closedAddresses[tran.From] = tran.ParentId
			break
		case OpenIncomeChannelTransactionType:
			address := helpers.CreateIncomeAddress(tran.PublicKey, peerPublicKey, tran.TimeStamp)
			if address == tran.To {
				parentId, ok := closedAddresses[address]

				if !ok {
					return i, nil
				}
				if parentId != tran.Id {
					return -1, errors.New("ledger: getOpenIncomeChannelPositionByPublicKey: inconsistent tran log")
				}

				return -1, nil
			}
			break
		}
	}

	return -1, nil
}

// Позиция транзакции в журнале
// Если не найдено: -1
func (l *Ledger) getOpenIncomeChannelPositionByAddress(address string) (int, error) {
	for i := len(l.transactions) - 1; i >= 0; i-- {
		tran := l.transactions[i]
		if !tran.IsValid() {
			continue
		}
		switch tran.Type {
		case CloseIncomeChannelTransactionType:
			if tran.From == address {
				return -1, errors.New("ledger: getOpenIncomeChannelPositionByAddress: income channel already close")
			}
			break
		case OpenIncomeChannelTransactionType:
			if tran.To == address {
				return i, nil
			}
		}
	}

	return -1, nil
}

func (l *Ledger) getOwnOpenIncomeChannelPositionByAddress(address string, pk string) (pos int, status ownChannelStatusType, err error) {
	pos, err = l.getOpenIncomeChannelPositionByAddress(address)

	if err != nil {
		return -1, invalidOwnChannelStatus, err
	}

	if pos < 0 {
		return -1, notOwnChannelStatus, nil
	}

	if !l.validateOwnIncomeAddress(address, pk, l.transactions[pos]) {
		return pos, notOwnChannelStatus, nil
	}

	return pos, ownChannelStatus, nil
}

// Позиция транзакции в журнале
// Если не найдено: -1
func (l *Ledger) getOpenSpendChannelPositionByPublicKey(peerPublicKey string) (int, error) {
	closedAddresses := make(map[string]string)

	for i := len(l.transactions) - 1; i >= 0; i-- {
		tran := l.transactions[i]
		if !tran.IsValid() {
			continue
		}
		switch tran.Type {
		case CloseSpendChannelTransactionType:
			closedAddresses[tran.From] = tran.ParentId
			break
		case OpenSpendChannelTransactionType:
			address := helpers.CreateSpendAddress(tran.PublicKey, peerPublicKey, tran.TimeStamp)
			if address == tran.To {
				parentId, ok := closedAddresses[address]

				if !ok {
					return i, nil
				}

				if parentId != tran.Id {
					return -1, errors.New("ledger: getOpenSpendChannelPositionByPublicKey: inconsistent tran log")
				}

				return -1, nil

			}
			break
		}
	}

	return -1, nil
}

// Позиция транзакции в журнале
// Если не найдено: -1
func (l *Ledger) getOpenSpendChannelPositionByAddress(address string) (int, error) {
	for i := len(l.transactions) - 1; i >= 0; i-- {
		tran := l.transactions[i]
		if !tran.IsValid() {
			continue
		}
		switch tran.Type {
		case CloseSpendChannelTransactionType:
			if tran.From == address {
				return -1, errors.New("ledger: getOpenSpendChannelPositionByAddress: spend channel already close")
			}
			break
		case OpenSpendChannelTransactionType:
			if tran.To == address {
				return i, nil
			}
		}
	}

	return -1, nil
}

func (l *Ledger) getOwnOpenSpendChannelPositionByAddress(address string, pk string) (pos int, status ownChannelStatusType, err error) {
	pos, err = l.getOpenSpendChannelPositionByAddress(address)

	if err != nil {
		return -1, invalidOwnChannelStatus, err
	}

	if pos < 0 {
		return -1, notOwnChannelStatus, nil
	}

	if !l.validateOwnSpendAddress(address, pk, l.transactions[pos]) {
		return pos, notOwnChannelStatus, nil
	}

	return pos, ownChannelStatus, nil
}

// Ожидается выполнение под блокировкой
// Ожидается допустимое состояние журнала транзакций
// Считается, что транзакция должна быть в журнале
// Считается, что канал не закрыт
func (l *Ledger) getIncomeChannelState(pos int) (*IncomeChannelState, error) {
	n := len(l.transactions)

	start := l.transactions[pos]
	state := NewIncomeChannelState(start)

	for i := pos + 1; i < n; i++ {
		tran := l.transactions[i]
		if !tran.IsValid() {
			continue
		}
		var err error = nil
		switch tran.Type {
		case OpenTransferChannelTransactionType:
			if tran.To == start.To {
				err = state.AddPlan(tran)
			}
			break
		case CloseTransferChannelTransactionType:
			if tran.To == start.To {
				err = state.AddFact(tran)
			}
			break
		}
		if err != nil {
			return nil, err
		}
	}

	return state, nil
}

func (l *Ledger) getSpendChannelState(pos int) (*SpendChannelState, error) {
	n := len(l.transactions)

	start := l.transactions[pos]
	state := NewSpendChannelState(start)

	for i := pos + 1; i < n; i++ {
		tran := l.transactions[i]
		if !tran.IsValid() {
			continue
		}
		var err error = nil

		switch tran.Type {
		case OpenTransferChannelTransactionType:
			if tran.From == start.To {
				err = state.AddPlan(tran)
			}
			break
		case CloseTransferChannelTransactionType:
			if tran.From == start.To {
				err = state.AddFact(tran)
			}
			break
		}
		if err != nil {
			return nil, err
		}
	}

	return state, nil
}

func (l *Ledger) getOpenTransferChannelTran(
	channelId string,
	fromAddress string,
	toAddress string,
	priceAmount TransactionAmountType,
	priceQuantumPower QuantumPowerType,
	plannedQuantumCount uint64,
) (*LedgerTransaction, error) {
	for i := len(l.transactions) - 1; i >= 0; i-- {
		tran := l.transactions[i]
		if !tran.IsValid() {
			continue
		}
		if tran.Type == OpenTransferChannelTransactionType && tran.ContractId == channelId {
			if tran.From == fromAddress &&
				tran.To == toAddress &&
				tran.PriceAmount == priceAmount &&
				tran.PriceQuantumPower == priceQuantumPower &&
				tran.PlannedQuantumCount == plannedQuantumCount {
				return tran, nil
			}

			return nil, errors.New("ledger: getOpenTransferChannelTran: tran log contain OTC tran with invalid values")
		}
	}

	return nil, nil
}

func (l *Ledger) getCloseTransferChannelTran(
	parentId string,
	channelId string,
	fromAddress string,
	toAddress string,
	priceAmount TransactionAmountType,
	priceQuantumPower QuantumPowerType,
	plannedQuantumCount uint64,
	amount TransactionAmountType,
	quantumCount uint64,
	volume uint64,
) (*LedgerTransaction, error) {
	for i := len(l.transactions) - 1; i >= 0; i-- {
		tran := l.transactions[i]
		if !tran.IsValid() {
			continue
		}
		if tran.Type == CloseTransferChannelTransactionType && tran.ContractId == channelId {
			if tran.ParentId == parentId &&
				tran.From == fromAddress &&
				tran.To == toAddress &&
				tran.PriceAmount == priceAmount &&
				tran.PriceQuantumPower == priceQuantumPower &&
				tran.PlannedQuantumCount == plannedQuantumCount &&
				tran.Amount == amount &&
				tran.QuantumCount == quantumCount &&
				tran.Volume == volume {
				return tran, nil
			}

			return nil, errors.New("ledger: getCloseTransferChannelTran: tran log contains CTC tran with invalid values")
		}

	}
	return nil, nil
}

func (l *Ledger) hasOpenedTransferChannelsInIncomeChannel(pos int) (bool, error) {
	address := l.transactions[pos].To
	transferChannels := make(map[string]bool)

	for i := pos + 1; i < len(l.transactions); i++ {
		tran := l.transactions[i]

		if !tran.IsValid() {
			continue
		}

		switch tran.Type {
		case OpenTransferChannelTransactionType:
			if tran.To == address {
				transferChannels[tran.Id] = true
			}
			break
		case CloseTransferChannelTransactionType:
			if tran.To == address {
				_, ok := transferChannels[tran.ParentId]
				if !ok {
					return false, errors.New("ledger: hasOpenedTransferChannelsInIncomeChannel: inconsistent tran log")
				}
				transferChannels[tran.ParentId] = false
			}
			break
		}
	}

	for _, isOpened := range transferChannels {
		if isOpened {
			return true, nil
		}
	}

	return false, nil
}

func (l *Ledger) hasOpenedTransferChannelsInSpendChannel(pos int) (bool, error) {
	address := l.transactions[pos].To
	transferChannels := make(map[string]bool)

	for i := pos + 1; i < len(l.transactions); i++ {
		tran := l.transactions[i]
		if !tran.IsValid() {
			continue
		}
		switch tran.Type {
		case OpenTransferChannelTransactionType:
			if tran.From == address {
				transferChannels[address] = true
			}
			break
		case CloseTransferChannelTransactionType:
			if tran.From == address {
				_, ok := transferChannels[address]

				if !ok {
					return false, errors.New("ledger: hasOpenedTransferChannelsInSpendChannel: inconsistent tran log")
				}

				transferChannels[address] = false
			}
			break
		}
	}

	for _, isOpened := range transferChannels {
		if isOpened {
			return true, nil
		}
	}

	return false, nil
}

func (l *Ledger) validateIncomeAddress(address string, pk string, openChannelTran *LedgerTransaction) bool {
	return address == helpers.CreateIncomeAddress(openChannelTran.PublicKey, pk, openChannelTran.TimeStamp)
}

func (l *Ledger) validateOwnIncomeAddress(address string, pk string, openChannelTran *LedgerTransaction) bool {
	return openChannelTran.PublicKey == l.pk && l.validateIncomeAddress(address, pk, openChannelTran)
}

func (l *Ledger) validateSpendAddress(address string, pk string, openChannelTran *LedgerTransaction) bool {
	return address == helpers.CreateSpendAddress(openChannelTran.PublicKey, pk, openChannelTran.TimeStamp)
}

func (l *Ledger) validateOwnSpendAddress(address string, pk string, openChannelTran *LedgerTransaction) bool {
	return openChannelTran.PublicKey == l.pk && l.validateSpendAddress(address, pk, openChannelTran)
}
