package job

import (
	"common"
	"svcledger/blockchain"
	"svcledger/logger"
	"svcledger/store"
)

func Syncer(ledger *store.Ledger, client *blockchain.Client) func() {
	return func() {
		common.Log.Event(
			logger.EventLedgerSyncBegin,
			common.Print(logger.EventLedgerSyncBegin),
		)
		defer common.Log.Event(
			logger.EventLedgerSyncEnd,
			common.Print(logger.EventLedgerSyncEnd),
		)
		t, err := client.GetTrans()

		if err != nil {
			common.Log.Error(
				logger.ErrorLedgerSyncConnectFailed,
				common.Printf(
					"%s: %v",
					logger.ErrorLedgerSyncConnectFailed,
					err.Error(),
				),
			)
			return
		}

		ledger.AddTrans(t.Data)
	}
}
