package job

import "svcledger/store"

func AutoCloseJob(ledger *store.Ledger) func() {
	return func() {
		ledger.AutoClose()
	}
}
