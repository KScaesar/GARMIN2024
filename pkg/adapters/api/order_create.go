package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/KScaesar/GARMIN2024/pkg"
	"github.com/KScaesar/GARMIN2024/pkg/app"
)

func createOrder(svc *app.Service) func(*gin.Context) {
	return func(c *gin.Context) {
		var req app.CreateOrderParam
		if !pkg.BindAndValidateJsonRequest(c, &req) {
			return
		}

		err := svc.OrderService.CreateOrder(c.Request.Context(), &req)
		if err != nil {
			pkg.ReplyErrorResponse(c, http.StatusInternalServerError, err)
			return
		}

		pkg.ReplySuccessResponse(c, http.StatusOK, nil)
	}
}
