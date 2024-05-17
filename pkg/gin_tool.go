package pkg

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var Validator = validator.New()

func BindAndValidateJsonRequest(c *gin.Context, obj any) (ok bool) {
	if err := c.ShouldBindJSON(obj); err != nil {
		Err := fmt.Errorf("bind payload: %w", err)
		ReplyErrorResponse(c, http.StatusBadRequest, Err)
		return false
	}

	var target *validator.InvalidValidationError
	err := Validator.StructCtx(c.Request.Context(), obj)
	if err != nil {
		if errors.As(err, &target) {
			return true
		}

		Err := fmt.Errorf("validate payload: %w", err)
		ReplyErrorResponse(c, http.StatusBadRequest, Err)
		return false
	}
	return true
}

func ReplyErrorResponse(c *gin.Context, httpCode int, err error) {
	resp := &HttpResponse{
		Message: err.Error(),
		Payload: struct{}{},
	}
	c.JSON(httpCode, resp)
	c.Abort()
}

func ReplySuccessResponse(c *gin.Context, httpCode int, payload any) {
	if payload == nil {
		payload = struct{}{}
	}
	resp := &HttpResponse{
		Message: "ok",
		Payload: payload,
	}
	c.JSON(httpCode, resp)
}

type HttpResponse struct {
	Message string `json:"msg"`
	Payload any    `json:"payload"`
}
