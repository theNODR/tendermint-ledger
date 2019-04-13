package store

import "common"

type SpendChannelState struct {
	*BaseChannelState
}

func NewSpendChannelState(tran *LedgerTransaction) *SpendChannelState {
	return &SpendChannelState{
		&BaseChannelState{
			Address: tran.To,
			Price: NewPrice(tran.PriceAmount, tran.PriceQuantumPower),
			State: NewSpecialChannelState(tran.Amount, tran.Volume),
			TimeLock: tran.TimeLock,
			LifeTime: tran.TimeLock - common.GetNowUnixMs(),
		},
	}
}

func (s *SpendChannelState) MaxPrice() Price {
	return s.Price
}

func (s *SpendChannelState) AddPlan(tran *LedgerTransaction) error {
	return s.State.DecreasePlan(tran.Id, tran.PlannedAmount(), tran.PlannedVolume())
}

func (s *SpendChannelState) AddFact(tran *LedgerTransaction) error {
	amount := tran.Amount
	quantumVolume := tran.QuantumCount * tran.PriceQuantumPower.Volume()
	volume := tran.Volume

	return s.State.DecreaseFact(tran.ParentId, amount, quantumVolume, volume)
}
