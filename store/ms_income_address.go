package store

type incomeAddressMemoryState struct {
	*baseAddressMemoryState
}

func newIncomeAddressMemoryState(tran *LedgerTransaction, isOwn bool) addressMemoryStater {
	return &incomeAddressMemoryState{
		newBaseAddressMemoryState(
			tran,
			0.0,
			isOwn,
			0,
		),
	}
}

func (s *incomeAddressMemoryState) CloseTransfer(id string, amount TransactionAmountType, volume uint64) {
	delete(s.transfers, id)
	s.amount += amount
	s.volume += volume
}
