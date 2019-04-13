package store

type JsonTransaction struct {
	Amount				uint64	`json:"amount"`
	ContractId			string	`json:"contract_id"`
	From				string	`json:"from"`
	Id					string	`json:"tx_id"`
	ParentId			string	`json:"parent_tx_id"`
	PlannedQuantumCount	uint64	`json:"planned_quantum_count"`
	PriceAmount			uint64	`json:"price_amount"`
	PriceQuantumPower	uint8	`json:"price_quantum_power"`
	PublicKey			string	`json:"pk"`
	QuantumCount		uint64	`json:"quantum_count"`
	Sign				string	`json:"sign"`
	//Status				int16	`json:"status"`
	TimeLock			int64	`json:"timelock"`
	TimeStamp			int64	`json:"timestamp"`
	To					string	`json:"to"`
	Type				int8	`json:"tx_type"`
	Version				string	`json:"version"`
	Volume				uint64	`json:"volume"`
	//Data				string	`json:"data"`
}

func (json *JsonTransaction) ToSignedTransaction() (*SignedTransaction, error) {
	signedTran := &SignedTransaction{
		Transaction: &Transaction{
			TransactionBody: &TransactionBody{
				Amount: TransactionAmountType(json.Amount),
				ContractId: json.ContractId,
				From: json.From,
				ParentId: json.ParentId,
				PlannedQuantumCount: json.PlannedQuantumCount,
				PriceAmount: TransactionAmountType(json.PriceAmount),
				PriceQuantumPower: QuantumPowerType(json.PriceQuantumPower),
				QuantumCount: json.QuantumCount,
				TimeStamp: json.TimeStamp,
				TimeLock: json.TimeLock,
				To: json.To,
				Type: TransactionType(json.Type),
				Version: json.Version,
				Volume: json.Volume,
				//Data: json.Data,
			},
			Id: json.Id,
			PublicKey: json.PublicKey,
		},
		Sign: json.Sign,
	}

	if err := signedTran.validate(); err != nil {
		return nil, err
	}

	return signedTran, nil
}
