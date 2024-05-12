package json

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type (
	jsonResponse struct {
		Code    string      `json:"responseCode"`
		Message string      `json:"message"`
		Data    interface{} `json:"data,omitempty"`
	}

	jsonResponsePage struct {
		Code    string      `json:"responseCode"`
		Message string      `json:"message,omitempty"`
		Data    interface{} `json:"data,omitempty"`
		Paging  *paging     `json:"paging,omitempty"`
	}

	jsonErrorResponse struct {
		Code    string `json:"responseCode"`
		Message string `json:"message"`
		Error   string `json:"error,omitempty"`
	}

	ValidationField struct {
		FieldName string `json:"field"`
		Message   string `json:"message"`
	}

	jsonBadRequestResponse struct {
		Code             string            `json:"responseCode"`
		Message          string            `json:"message"`
		ErrorDescription []ValidationField `json:"errorDescription,omitempty"`
	}

	paging struct {
		Page      int `json:"page"`
		TotalData int `json:"size"`
	}
)

func NewResponseSuccess(c *gin.Context, result interface{}, message, serviceCode, responseCode string) {
	c.JSON(http.StatusOK, jsonResponse{
		Code:    "200" + serviceCode + responseCode,
		Message: message,
		Data:    result,
	})
}

func NewResponseSuccessPage(c *gin.Context, result interface{}, page, count int, message, serviceCode, responseCode string) {
	c.JSON(http.StatusOK, jsonResponsePage{
		Code:    "200" + serviceCode + responseCode,
		Message: message,
		Data:    result,
		Paging: &paging{ // Note the use of address to pass it as a pointer
			Page:      page,
			TotalData: count,
		},
	})
}

func NewResponseBadRequest(c *gin.Context, validationField []ValidationField, message, serviceCode, errorCode string) {
	c.JSON(http.StatusBadRequest, jsonBadRequestResponse{
		Code:             "400" + serviceCode + errorCode,
		Message:          message,
		ErrorDescription: validationField,
	})
}

func NewResponseError(c *gin.Context, err, serviceCode, errorCode string) {
	log.Error().Msg(err)
	c.JSON(http.StatusInternalServerError, jsonErrorResponse{
		Code:    "500" + serviceCode + errorCode,
		Message: "internal server error",
		Error:   err,
	})
}

func NewResponseForbidden(c *gin.Context, message, serviceCode, errorCode string) {
	c.JSON(http.StatusForbidden, jsonResponse{
		Code:    "403" + serviceCode + errorCode,
		Message: message,
	})
}

func NewResponseUnauthorized(c *gin.Context, message, serviceCode, errorCode string) {
	c.JSON(http.StatusUnauthorized, jsonResponse{
		Code:    "401" + serviceCode + errorCode,
		Message: message,
	})
}
