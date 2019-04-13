package store

type EarnTokens struct {
	Fact	TransactionAmountType
	Plan	TransactionAmountType
}

func (e *EarnTokens) ToState() *EarnTokensState {
	return &EarnTokensState{
		Fact: e.Fact.ToString(),
		Plan: e.Plan.ToString(),
	}
}

type CommonAddressBalance struct {
	Amount		TransactionAmountType
	Profit		TransactionAmountType
	OwnEarn		*EarnTokens
	ForeignEarn	*EarnTokens
}

type EarnTokensState struct {
	Fact	string	`json:"fact"`
	Plan	string	`json:"plan"`
}

type CommonAddressBalanceState struct {
	Profit	string				`json:"profit"`
	Own		*EarnTokensState	`json:"own"`
	Foreign	*EarnTokensState	`json:"foreign"`
}

func NewCommonAddressBalance() *CommonAddressBalance {
	return &CommonAddressBalance{
		Amount: 0,
		Profit: 0,
		OwnEarn: &EarnTokens{
			Fact: 0,
			Plan: 0,
		},
		ForeignEarn: &EarnTokens{
			Fact: 0,
			Plan: 0,
		},
	}
}

func (b *CommonAddressBalance) ToState() *CommonAddressBalanceState {
	return &CommonAddressBalanceState{
		Profit: b.Profit.ToString(),
		Own: b.OwnEarn.ToState(),
		Foreign: b.ForeignEarn.ToState(),
	}
}
