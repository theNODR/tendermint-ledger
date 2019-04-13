package store

type spendAddressMemoryState struct {
	*baseAddressMemoryState
}

func newSpendAddressMemoryState(tran *LedgerTransaction, isOwn bool) addressMemoryStater {
	return &spendAddressMemoryState{
		newBaseAddressMemoryState(
			tran,
			tran.Amount,
			isOwn,
			0,
		),
	}
}

func (s *spendAddressMemoryState) CloseTransfer(id string, amount TransactionAmountType, volume uint64) {
	delete(s.transfers, id)
	s.volume += volume
	s.amount -= amount
}
