package tendermint

import (
	"common"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"svcledger/blockchain"
	"svcledger/blockchain/trans"
	"svcnodr/msgs"
	"svcnodr/types"
)

const denom = "ndr"

func CreateMsgFromTransaction(tran *blockchain.Tran) sdk.Msg {
	switch tran.Type {
	case blockchain.InitTransactionType:
		initTran := tran.Data.(*trans.InitTran)
		return CreateInitMsg(initTran)
	case blockchain.AutoCloseTransactionType:
		return CreateAutoCloseMsg(tran.Data.(*trans.AutoCloseTran))
	case blockchain.CloseSpendChannelTransactionType:
		return CreateCloseSpendMsg(tran.Data.(*trans.CloseSpendTran))
	case blockchain.CloseIncomeChannelTransactionType:
		return CreateCloseIncomeMsg(tran.Data.(*trans.CloseIncomeTran))
	case blockchain.CloseTransferChannelTransactionType:
		return CreateCloseTransferMsg(tran.Data.(*trans.CloseTransferTran))
	case blockchain.OpenIncomeChannelTransactionType:
		return CreateOpenIncomeMsg(tran)
	case blockchain.OpenSpendChannelTransactionType:
		return CreateOpenSpendMsg(tran)
	case blockchain.OpenTransferChannelTransactionType:
		return CreateOpenTransferMsg(tran.Data.(*trans.OpenTransferTran))
	default:
		return msgs.NewMsgDebug(tran.ToString())
	}
}

func CreateAutoCloseMsg(tran *trans.AutoCloseTran) msgs.MsgAutoclose {
	return msgs.NewMsgAutoclose(
		common.GetNowUnixMs(),
		createAccAddress(tran.LedgerAddress))
}

func CreateInitMsg(tran *trans.InitTran) msgs.MsgInit {
	addr := createAccAddress(tran.To)
	coins := createCoins(tran.Amount)
	return msgs.NewMsgInit(addr, coins)
}

func CreateOpenSpendMsg(t *blockchain.Tran) msgs.MsgOpenSpend {
	tran := t.Data.(*trans.OpenSpendTran)
	return msgs.NewMsgOpenSpend(
		t.PublicKey,
		tran.PeerPublicKey,
		createAccAddress(tran.From),
		createCoins(tran.MaxAmount),
		createCoins(tran.PriceAmount),
		tran.PriceQuantumPower,
		tran.LifeTime)
}

func CreateOpenTransferMsg(tran *trans.OpenTransferTran) msgs.MsgOpenTransfer {
	return msgs.NewMsgOpenTransfer(
		tran.ChannelId,
		createAccAddress(tran.FromAddress),
		tran.FromPk,
		createAccAddress(tran.ToAddress),
		tran.ToPk,
		createCoins(tran.PriceAmount),
		tran.PriceQuantumPower,
		tran.PlannedQuantumCount,
		tran.LifeTime,
	)
}

func CreateOpenIncomeMsg(t *blockchain.Tran) msgs.MsgOpenIncome {
	tran := t.Data.(*trans.OpenIncomeTran)
	return msgs.NewMsgOpenIncome(
		t.PublicKey,
		tran.PeerPublicKey,
		createAccAddress(tran.From),
		createCoins(tran.PriceAmount),
		types.QuantumPowerType(tran.PriceQuantumPower),
		tran.LifeTime,
	)
}

func CreateCloseSpendMsg(tran *trans.CloseSpendTran) msgs.MsgCloseSpend {
	return msgs.NewMsgCloseSpend(
		tran.PeerPublicKey,
		createAccAddress(tran.From),
		createAccAddress(tran.To),
	)
}

func CreateCloseIncomeMsg(tran *trans.CloseIncomeTran) msgs.MsgCloseIncome {
	return msgs.NewMsgCloseIncome(
		tran.PeerPublicKey,
		createAccAddress(tran.From),
		createAccAddress(tran.To),
	)
}

func CreateCloseTransferMsg(tran *trans.CloseTransferTran) msgs.MsgCloseTransfer {
	return msgs.NewMsgCloseTransfer(

		tran.ChannelId,
		createAccAddress(tran.FromAddress),
		tran.FromPk,
		createAccAddress(tran.ToAddress),
		tran.ToPk,
		tran.PriceAmount,
		tran.PriceQuantumPower,
		createCoins(tran.Amount),
		tran.QuantumCount,
		tran.Volume,
	)
}

func createAccAddress(addr string) sdk.AccAddress {
	return sdk.AccAddress(addr)
}

func createCoins(amount uint64) sdk.Coins {
	coin := sdk.NewInt64Coin(denom, int64(amount))
	coins := sdk.Coins{}
	coins = append(coins, coin)
	return coins
}
