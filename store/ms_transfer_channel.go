package store

type transferChannelMemoryStater interface {
	Close(tran *LedgerTransaction)
	IsAutoClosable(timestamp int64) bool
	Tran() *LedgerTransaction
}

type transferChannelMemoryState struct {
	tran	*LedgerTransaction

	income	addressMemoryStater
	spend	addressMemoryStater
}

func newTransferChannelMemoryState(
	tran *LedgerTransaction,
	income addressMemoryStater,
	spend addressMemoryStater,
) transferChannelMemoryStater {
	return &transferChannelMemoryState{
		tran: tran,
		income: income,
		spend: spend,
	}
}

func (s *transferChannelMemoryState) Close(tran *LedgerTransaction) {
	var amount TransactionAmountType = 0
	var volume uint64 = 0

	if tran != nil {
		amount = tran.Amount
		volume = tran.Volume
	}
	if s.income != nil {
		s.income.CloseTransfer(s.tran.Id, amount, volume)
	}
	s.spend.CloseTransfer(s.tran.Id, amount, volume)
}

func (s *transferChannelMemoryState) IsAutoClosable(timestamp int64) bool {
	if s.income == nil {
		return s.spend.IsOwn() && s.TimeLock() < timestamp
	} else {
		return s.income.IsOwn() && s.TimeLock() < timestamp
	}
}

func (s *transferChannelMemoryState) Tran() *LedgerTransaction {
	return s.tran
}

func (s *transferChannelMemoryState) TimeLock() int64 {
	return s.tran.TimeLock
}
