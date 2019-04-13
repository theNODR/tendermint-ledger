package store

type addressMemoryStater interface {
	Amount() TransactionAmountType
	CloseTransfer(id string, amount TransactionAmountType, volume uint64)
	IsAutoClosable(timestamp int64) bool
	IsOwn() bool
	Tran() *LedgerTransaction
	Volume() uint64
}

type baseAddressMemoryState struct {
	tran	*LedgerTransaction

	amount	TransactionAmountType
	isOwn	bool
	volume	uint64

	transfers map[string]transferChannelMemoryStater
}

func newBaseAddressMemoryState(
	tran *LedgerTransaction,
	amount TransactionAmountType,
	isOwn bool,
	volume uint64,
) *baseAddressMemoryState {
	return &baseAddressMemoryState{
		tran: tran,
		amount: amount,
		isOwn: isOwn,
		volume: volume,
		transfers: make(map[string]transferChannelMemoryStater),
	}
}

func (s *baseAddressMemoryState) Amount() TransactionAmountType {
	return s.amount
}

func (s *baseAddressMemoryState) IsAutoClosable(timestamp int64) bool {
	return s.IsOwn() && len(s.transfers) == 0 && s.TimeLock() < timestamp
}

func (s *baseAddressMemoryState) IsOwn() bool {
	return s.isOwn
}

func (s *baseAddressMemoryState) TimeLock() int64 {
	return s.tran.TimeLock
}

func (s *baseAddressMemoryState) Tran() *LedgerTransaction {
	return s.tran
}

func (s *baseAddressMemoryState) Volume() uint64 {
	return s.volume
}
