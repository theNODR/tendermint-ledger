package httpHandlers

import (
	"app/httputil"
	"common"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"svcledger/blockchain/tendermint"
)

type IncomeRequest struct {
	Addresses []string `json:"addresses"`
}

type IncomeResponse struct {
	AddressesWithBalances map[string]int64 `json:"balances"`
}

func HandleIncomeState(client *tendermint.TendermintClient) gin.HandlerFunc {
	return func(context *gin.Context) {
		incomeRequest := &IncomeRequest{}
		err := context.ShouldBindJSON(incomeRequest)
		if err != nil {
			httputil.SendApiError(context, http.StatusBadRequest, "Can`t parse income addresses", err)
		}

		balances := make(map[string]int64)
		for i := 0; i < len(incomeRequest.Addresses); i++ {
			state, err := client.GetIncomeState(incomeRequest.Addresses[i])
			if err != nil {
				common.Log.Error("ErrorGetIncomeState",
					common.Printf("Erorr getting income state from tndrmn: %v", err))
				continue
			}

			balances[state.Address] = state.Current.Amount
		}

		response := &IncomeResponse{}
		response.AddressesWithBalances = balances
		bytes, err := json.Marshal(response)
		if err != nil {
			common.Log.Error("ErrorCreateResponse",
				common.Printf("Erorr create json: %v", err))
			httputil.SendApiError(context, http.StatusInternalServerError, "Can`t create response:", err)
		}

		httputil.SendApiResult(context, bytes)
	}
}
