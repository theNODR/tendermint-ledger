package store

import (
	"svcledger/helpers"
)

const (
	finBlockThresholdTokens			TransactionAmountType	= 100000
	finBlockThresholdTrafficVolume	uint64					= 1 << 30
	finBlockThresholdTranCount		uint64					= 100
)

type FinancialTran struct {
	SecondParty		string					`json:"second_party"`
	Amount			TransactionAmountType	`json:"amount"`
}

type FinancialData struct {
	Amount	TransactionAmountType	`json:"amount"`
	Volume	uint64					`json:"volume"`
	To		[]*FinancialTran		`json:"to"`
	From	[]*FinancialTran		`json:"from"`
}

type transactionAmountTypeDict map[string]TransactionAmountType

func (t transactionAmountTypeDict) toArray() []*FinancialTran {
	res := make([]*FinancialTran, 0)

	for key, item := range t {
		res = append(res, &FinancialTran{
			Amount: item,
			SecondParty: key,
		})
	}

	return res
}

type FinancialState interface {
	Add(tran *LedgerTransaction)
	GetFinancialData() *FinancialData
}

type financialState struct {
	address			string
	pk				string

	// На какие адреса переводятся деньги с этого адреса
	to				transactionAmountTypeDict
	// С какого адреса переводятся деньги на этот адрес
	from			transactionAmountTypeDict

	amount			TransactionAmountType
	count			uint64
	volume			uint64

	ownIncome		map[string]*LedgerTransaction
	ownSpend		map[string]*LedgerTransaction
	foreignIncome	map[string]*LedgerTransaction
	foreignSpend	map[string]*LedgerTransaction
}

func NewFinancialState(address string, pk string) FinancialState {
	return &financialState{
		address: address,
		pk: pk,

		to: make(transactionAmountTypeDict),
		from: make(transactionAmountTypeDict),

		ownIncome: make(map[string]*LedgerTransaction),
		ownSpend: make(map[string]*LedgerTransaction),
		foreignIncome: make(map[string]*LedgerTransaction),
		foreignSpend: make(map[string]*LedgerTransaction),
	}
}

func (f *financialState) Add(tran *LedgerTransaction) {
	if !tran.IsValid() {
		return
	}

	switch tran.Type {
	case OpenIncomeChannelTransactionType:
		if tran.From == f.address {
			f.ownIncome[tran.To] = tran
		} else {
			f.foreignSpend[tran.To] = tran
		}
		break
	case OpenSpendChannelTransactionType:
		if tran.From == f.address {
			f.ownSpend[tran.To] = tran
		} else {
			f.foreignSpend[tran.To] = tran
		}
		break
	case OpenTransferChannelTransactionType:
		break
	case CloseTransferChannelTransactionType:
		if tran.PublicKey != f.pk {
			break
		}

		f.amount += tran.Amount
		f.count++
		f.volume += tran.AmountedVolume()

		toAddress := tran.To
		if helpers.IsAddressCommon(toAddress) {
			if v, ok := f.to[toAddress]; ok {
				f.to[toAddress] = v + tran.Amount
			} else {
				f.to[toAddress] = tran.Amount
			}
		} else {
			if foreign, ok := f.foreignSpend[tran.From]; ok {
				if v, ok := f.from[foreign.From]; ok {
					f.from[foreign.From] = v + tran.Amount
				} else {
					f.from[foreign.From] = tran.Amount
				}
			}
		}
		break
	case CloseSpendChannelTransactionType:
		if tran.To == f.address {
			delete(f.ownIncome, tran.From)
		} else {
			delete(f.foreignIncome, tran.From)
		}
		break
	case CloseIncomeChannelTransactionType:
		if tran.To == f.address {
			delete(f.ownIncome, tran.From)
		} else {
			delete(f.foreignIncome, tran.From)
		}
		break
	case FinancialBlockTransactionType:
		if tran.From == f.address &&
			tran.PublicKey == f.pk {
			f.flush()
		}
		break
	}
}

func (f *financialState) GetFinancialData() *FinancialData {
	if f.count > finBlockThresholdTranCount ||
		f.amount > finBlockThresholdTokens ||
		f.volume > finBlockThresholdTrafficVolume {
		return &FinancialData{
			Amount: f.amount,
			Volume: f.volume,
			From: f.from.toArray(),
			To: f.to.toArray(),
		}
	} else {
		return nil
	}
}

func (f *financialState) flush() {
	f.to = make(transactionAmountTypeDict)
	f.from = make(transactionAmountTypeDict)

	f.amount = 0
	f.count = 0
	f.volume = 0
}
