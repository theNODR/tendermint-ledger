package httpHelpers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"app/httputil"
	"svcledger/netHelpers"
)

func Handler(createRequest netHelpers.CreateRequester) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		req, err := createRequest(NewBaseRequest())

		if err != nil {
			httputil.SendApiError(ctx, http.StatusBadRequest, "bad data", err)
			return
		}

		data, err := req.Handle()

		if err != nil {
			httputil.SendApiError(ctx, http.StatusInternalServerError, "", err)
			return
		}

		resp := NewResponse(ctx, data)
		err = resp.Send()

		if err != nil {
			httputil.SendApiError(ctx, http.StatusInternalServerError, "", err)
			return
		}
	}
}