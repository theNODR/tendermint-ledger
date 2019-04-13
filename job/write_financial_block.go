package job

import "svcledger/store"

func WriteFinancialBlock(ledger *store.Ledger) func() {
	return func() {
		//ledger.WriteFinancialTransaction()
	}
}
