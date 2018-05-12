package app

import (
	"github.com/gin-gonic/gin"
	"github.com/sjeanpierre/rs_input_tracker_go/app/models"
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

func listArrayInputsEndpoint(c *gin.Context) {
	account := c.Param("account_id")
	array := c.Param("array_id")
	inputs := models.ListCurrentInputs(account,array)
	c.JSON(200,inputs)
}

func listArraysEndpoint(c *gin.Context) {
	account := c.Param("account_id")
	arrays := models.ListCurrentArraysByAccount(account)
	c.JSON(200,arrays)
}

func listArraysVersionsEndpoint(c *gin.Context) {
	account := c.Param("account_id")
	arrayID := c.Param("array_id")
	arrays := models.ListArrayVersions(account,arrayID)
	c.JSON(200,arrays)
}


func listAccountsEndpoint(c *gin.Context) {
	accounts := models.ListAccounts()
	c.JSON(200,accounts)
}

func getAccountEndpoint(c *gin.Context) {
	a := c.Param("account_id")
	account := models.GetAccount(a)
	c.JSON(200,account)
}

func listInputVersionsEndpoint(c *gin.Context) {
	account := c.Param("account_id")
	array := c.Param("array_id")
	input := c.Param("input_name")
	inputs := models.ListInputVersions(account,array,input)
	c.JSON(200,inputs)
}

func getArrayEndpoint(c *gin.Context)  {
	account := c.Param("account_id")
	a := c.Param("array_id")
	array := models.GetArray(account,a)
	c.JSON(200,array)
}
