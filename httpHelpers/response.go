package httpHelpers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/json"

	"app/httputil"
	"common"
	"svcledger/netHelpers"
)

type Response struct {
	ctx			*gin.Context
	data		interface{}
}

func NewResponse(ctx *gin.Context, data interface{}) netHelpers.Responser {
	return &Response{
		ctx: ctx,
		data: data,
	}
}

func (resp *Response) Send() error{
	if resp.data == nil {
		ok, _ := json.Marshal("ok")
		httputil.SendApiResult(resp.ctx, ok)
	} else if x, err := json.Marshal(resp.data); err != nil {
		httputil.SendApiError(resp.ctx, http.StatusInternalServerError, "", err)
		return err
	} else {
		common.Log.PrintFull(common.Print("json: %v", string(x[:]), ), )
		httputil.SendApiResult(resp.ctx, x)
	}

	return nil
}
