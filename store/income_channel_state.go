package store

import "common"

type IncomeChannelState struct {
	*BaseChannelState
}

func NewIncomeChannelState(tran *LedgerTransaction) *IncomeChannelState {
	return &IncomeChannelState{
		&BaseChannelState{
			Address: tran.To,
			Price: NewPrice(tran.PriceAmount, tran.PriceQuantumPower),
			State: NewSpecialChannelState(0, 0),
			TimeLock: tran.TimeLock,
			LifeTime: tran.TimeLock - common.GetNowUnixMs(),
		},
	}
}

func (s *IncomeChannelState) MinPrice() Price {
	return s.Price
}

func (s *IncomeChannelState) AddPlan(tran *LedgerTransaction) error {
	return s.State.IncreasePlan(tran.Id, tran.PlannedAmount(), tran.PlannedVolume())
}

func (s *IncomeChannelState) AddFact(tran *LedgerTransaction) error {
	amount := tran.Amount
	quantumVolume := tran.QuantumCount * tran.PriceQuantumPower.Volume()
	volume := tran.Volume
	return s.State.IncreaseFact(tran.ParentId, amount, quantumVolume, volume)
}
