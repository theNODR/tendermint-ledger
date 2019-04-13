package store

import (
	"common"
)

type memoryState struct {
	pk					string
	timestamp			int64

	incomeAddresses		map[string]addressMemoryStater
	spendAddresses		map[string]addressMemoryStater
	transferChannels	map[string]transferChannelMemoryStater
}

func newMemoryState(pk string, timestamp int64) *memoryState {
	return &memoryState{
		pk: pk,
		timestamp: timestamp,
		incomeAddresses: make(map[string]addressMemoryStater),
		spendAddresses: make(map[string]addressMemoryStater),
		transferChannels: make(map[string]transferChannelMemoryStater),
	}
}

func (s *memoryState) Init(trans LedgerTransactions) {
	for i, n := 0, len(trans); i < n; i++ {
		item := trans[i]
		if item.Status != ValidTransactionStatus {
			continue
		}

		switch item.Type {
		case OpenIncomeChannelTransactionType:
			s.openIncomeAddress(item)
			break
		case OpenSpendChannelTransactionType:
			s.openSpendAddress(item)
			break
		case OpenTransferChannelTransactionType:
			s.openTransferChannel(item)
			break
		case CloseTransferChannelTransactionType:
			s.closeTransferChannel(i, item)
			break
		case CloseSpendChannelTransactionType:
			s.closeSpendAddress(item)
			break
		case CloseIncomeChannelTransactionType:
			s.closeIncomeAddress(item)
			break
		}
	}
}

func (s *memoryState) GetTransferChannelsForClose() []transferChannelMemoryStater {
	res := make([]transferChannelMemoryStater, 0, len(s.transferChannels))

	for _, item := range s.transferChannels {
		if item.IsAutoClosable(s.timestamp) {
			res = append(res, item)
		}
	}

	return res
}

func (s *memoryState) getAddressesForClose(data map[string]addressMemoryStater) []addressMemoryStater {
	res := make([]addressMemoryStater, 0, len(data))

	for _, item := range data {
		if item.IsAutoClosable(s.timestamp) {
			res = append(res, item)
		}
	}

	return res
}

func (s *memoryState) GetSpendAddressesForClose() []addressMemoryStater {
	return s.getAddressesForClose(s.spendAddresses)
}

func (s *memoryState) GetIncomeAddressesForClose() []addressMemoryStater {
	return s.getAddressesForClose(s.incomeAddresses)
}

func (s *memoryState) openIncomeAddress(tran *LedgerTransaction) {
	s.incomeAddresses[tran.To] = newIncomeAddressMemoryState(
		tran,
		tran.PublicKey == s.pk,
	)
}

func (s *memoryState) closeIncomeAddress(tran *LedgerTransaction) {
	delete(s.incomeAddresses, tran.From)
}

func (s *memoryState) openSpendAddress(tran *LedgerTransaction) {
	s.spendAddresses[tran.To] = newSpendAddressMemoryState(
		tran,
		tran.PublicKey == s.pk,
	)
}

func (s *memoryState) closeSpendAddress(tran *LedgerTransaction) {
	delete(s.spendAddresses, tran.From)
}

func (s *memoryState) openTransferChannel(tran *LedgerTransaction) {
	if _, ok := s.transferChannels[tran.Id]; ok {
		common.Log.Panic(common.Printf("inconsistent tran log: tran has already existed in log", ))
	}
	var incomeAddress addressMemoryStater = nil
	if addr, ok := s.incomeAddresses[tran.To]; ok {
		incomeAddress = addr
	}
	var spendAddress addressMemoryStater = nil
	if addr, ok := s.spendAddresses[tran.From]; ok {
		spendAddress = addr
	} else {
		common.Log.Panic(common.Printf("inconsistent tran log: spend address didn't find"))
	}
	s.transferChannels[tran.Id] = newTransferChannelMemoryState(
		tran,
		incomeAddress,
		spendAddress,
	)
}

func (s *memoryState) closeTransferChannel(i int, tran *LedgerTransaction) {
	if item, ok := s.transferChannels[tran.ParentId]; ok {
		item.Close(tran)

		delete(s.transferChannels, tran.ParentId)
	} else {
		common.Log.Panic(common.Printf("inconsistent tran log: open tran hasn't existed in log", ))
	}
}
