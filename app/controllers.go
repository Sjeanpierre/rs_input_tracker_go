package app

import (
	"github.com/gin-gonic/gin"
	"github.com/sjeanpierre/rs_input_tracker_go/app/models"
)

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

func compareArrayInputs(c *gin.Context)  {
	account := c.Param("account_id")
	a1 := c.Param("array_1")
	a2 := c.Param("array_2")
	diff := models.CompareArrayInputs(account,a1,a2)
	c.JSON(200,diff)
}
