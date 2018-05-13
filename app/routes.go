package app

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(api *gin.RouterGroup) *gin.RouterGroup{
	api.GET("/accounts/:account_id/arrays/:array_id/inputs", listArrayInputsEndpoint)
	api.GET("/accounts/:account_id/arrays/:array_id", getArrayEndpoint)
	api.GET("/accounts/:account_id/arrays/:array_id/history", listArraysVersionsEndpoint)
	api.GET("/accounts/:account_id/arrays",listArraysEndpoint)
	api.GET("/accounts/:account_id",getAccountEndpoint)
	api.GET("/accounts",listAccountsEndpoint)
	api.GET("/accounts/:account_id/arrays/:array_id/inputs/:input_name",listInputVersionsEndpoint)
	return api
}
