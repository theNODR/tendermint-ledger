package tendermint

import (
	"common"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/utils"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authtxb "github.com/cosmos/cosmos-sdk/x/auth/client/txbuilder"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"svcledger/blockchain"
	"svcledger/helpers"
	"svcnodr"
	"svcnodr/handlers"
	ntypes "svcnodr/types"
)

const (
	errorSendToTendermint = "tendermint: send"
	errorCreateTx         = "tendermint: create tx"
)

const queryPath = "custom/svcnodr/%s"

type Response = []byte

type TendermintClient struct {
	ctx       context.CLIContext
	cdc       *codec.Codec
	txBuilder authtxb.TxBuilder
	keys      helpers.KeyPair
}

func NewTendermintClient(chainId string, tendermintEndPoint string) *TendermintClient {
	viper.Set(client.FlagChainID, chainId)
	viper.Set(client.FlagNode, tendermintEndPoint)

	cdc := svcnodr.MakeCodec()

	return &TendermintClient{
		cdc:       cdc,
		ctx:       context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc),
		txBuilder: authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc)),
	}
}

func (client *TendermintClient) SendInit(trans *blockchain.Tran) error {
	_, err := client.sendTransactionSync(trans)
	return err
}

func (client *TendermintClient) SendOpenIncomeChannel(tran *blockchain.Tran) (*ntypes.IncomeChannelStateResponse, error) {
	res, err := client.sendTransactionSync(tran)
	if err != nil {
		return nil, err
	}

	var state *ntypes.IncomeChannelStateResponse
	client.cdc.MustUnmarshalJSON(res, &state)
	return state, nil
}

func (client *TendermintClient) SendCloseIncome(tran *blockchain.Tran) (*handlers.CloseIncomeResponse, error) {
	response, err := client.sendTransactionSync(tran)
	if err != nil {
		return nil, err
	}

	state := &handlers.CloseIncomeResponse{}
	err = client.cdc.UnmarshalJSON(response, state)
	if err != nil {
		return nil, err
	}

	return state, err
}

func (client *TendermintClient) SendAutoclose(tran *blockchain.Tran) error {
	_, err := client.sendTransactionSync(tran)
	return err
}

func (client *TendermintClient) SendOpenSpendChannel(tran *blockchain.Tran) (*ntypes.SpendChannelStateResponse, error) {
	response, err := client.sendTransactionSync(tran)
	if err != nil {
		return nil, err
	}
	state := &ntypes.SpendChannelStateResponse{}
	err = client.cdc.UnmarshalJSON(response, &state)
	if err != nil {
		return nil, err
	}

	return state, err
}

func (client *TendermintClient) SendOpenTransferChannel(tran *blockchain.Tran) (*handlers.OpenTransferResponse, error) {
	response, err := client.sendTransactionSync(tran)
	if err != nil {
		return nil, err
	}

	state := &handlers.OpenTransferResponse{}
	err = client.cdc.UnmarshalJSON(response, state)
	if err != nil {
		return nil, err
	}

	return state, err
}

func (client *TendermintClient) SendCloseTransfer(tran *blockchain.Tran) error {
	_, err := client.sendTransactionSync(tran)
	return err
}

func (client *TendermintClient) SendCloseSpend(tran *blockchain.Tran) (*handlers.CloseSpendResponse, error) {
	response, err := client.sendTransactionSync(tran)
	if err != nil {
		return nil, err
	}

	state := &handlers.CloseSpendResponse{}
	err = client.cdc.UnmarshalJSON(response, state)
	if err != nil {
		return nil, err
	}

	return state, err
}

func (client *TendermintClient) GetIncomeState(address string) (*ntypes.IncomeChannelStateResponse, error) {
	accAddress := createAccAddress(address)
	requestData := client.cdc.MustMarshalJSON(&ntypes.GetIncomeStateRequest{Address: string(accAddress)})
	responseBytes, err := client.ctx.QueryWithData(fmt.Sprintf(queryPath, svcnodr.QueryIncomeAddressState), requestData)
	if err != nil {
		return nil, err
	}

	state := &ntypes.IncomeChannelStateResponse{}
	err = client.cdc.UnmarshalJSON(responseBytes, state)
	if err != nil {
		return nil, err
	}

	return state, nil
}

func (client *TendermintClient) GetSpendChannelState(address string) (*ntypes.SpendChannelStateResponse, error) {
	accAddress := createAccAddress(address)
	requestData := client.cdc.MustMarshalJSON(&ntypes.GetSpendStateRequest{Address: string(accAddress)})
	response, err := client.ctx.QueryWithData(fmt.Sprintf(queryPath, svcnodr.QuerySpendAddressState), requestData)
	if err != nil {
		return nil, err
	}
	state := &ntypes.SpendChannelStateResponse{}
	err = client.cdc.UnmarshalJSON(response, state)
	if err != nil {
		return nil, err
	}

	return state, nil
}

func (client *TendermintClient) GetTransferChannelState(channelId string) (*ntypes.TransferChannelStateResponse, error) {
	requestData, err := client.cdc.MarshalJSON(&ntypes.GetTransferStateRequest{ChannelId: channelId})
	if err != nil {
		return nil, err
	}

	response, err := client.ctx.QueryWithData(fmt.Sprintf(queryPath, svcnodr.QueryTransferAddressState), requestData)
	if err != nil {
		return nil, err
	}

	state := &ntypes.TransferChannelStateResponse{}
	err = client.cdc.UnmarshalJSON(response, state)
	if err != nil {
		return nil, err
	}

	return state, nil
}

func (client *TendermintClient) sendTransaction(transactions *blockchain.Tran) chan Response {
	responseChan := make(chan Response)

	go func() {
		fmt.Printf("\nSend transactin: %s \n", transactions.ToString())

		msg := CreateMsgFromTransaction(transactions)

		coins := make([]sdk.Coin, 0)
		coins = append(coins, sdk.NewCoin(denom, sdk.ZeroInt()))
		fee := auth.NewStdFee(0, coins)

		stdTx := auth.NewStdTx([]sdk.Msg{msg}, fee, nil, "")
		tx, err := client.txBuilder.TxEncoder()(stdTx)
		if err != nil {
			common.Log.Error(errorCreateTx, common.Printf("can`t Create ecnoded tx: %v", err))

			return
		}

		response, err := client.ctx.BroadcastTx(tx)
		fmt.Printf("%d Transaction %s response code: %s", response.Code, transactions.Id, response.String())
		if err != nil {
			close(responseChan)
			return
		}

		if len(response.Data) > 0 {
			responseChan <- response.Data
			close(responseChan)
		}
	}()

	return responseChan
}

func (client *TendermintClient) sendTransactionSync(transactions *blockchain.Tran) (Response, error) {
	msg := CreateMsgFromTransaction(transactions)

	var err error
	err = msg.ValidateBasic()
	if err != nil {
		return nil, err
	}

	coins := make([]sdk.Coin, 0)
	coins = append(coins, sdk.NewCoin(denom, sdk.ZeroInt()))
	fee := auth.NewStdFee(0, coins)

	sigs := []auth.StdSignature{{}}
	stdTx := auth.NewStdTx([]sdk.Msg{msg}, fee, sigs, transactions.Id)
	tx, err := client.txBuilder.TxEncoder()(stdTx)

	if err != nil {
		common.Log.Error(errorCreateTx, common.Printf("can`t Create ecnoded tx: %v", err))
		return nil, err
	}

	response, err := client.ctx.BroadcastTx(tx)
	if err != nil {
		fmt.Printf("\n Send %s \nError: %v \n", transactions.ToString(), err)
	}

	if err != nil {
		return nil, err
	}

	if response.Code == 0 {
		return response.Data, nil
	}

	return nil, errors.New(response.Logs.String())
}
