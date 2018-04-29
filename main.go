package main

import (
	"fmt"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/static"
	//"net/http"
	"log"
)

func InputFunctions() {
	listTables()
	currentInputs := listCurrentInputs("61931","300594003")
	versionList := listInputVersions("61931","300594003","TARBALL")
	ci,err := json.Marshal(currentInputs)
	if err != nil {
		fmt.Println("Could not marshall current list")
	}
	vi,err := json.Marshal(versionList)
	if err != nil {
		fmt.Println("Could not marshall version list")
	}
	fmt.Printf("Current Inputs %s\n",ci)
	fmt.Printf("Version Inputs %s\n",vi)
	fmt.Print("\ndone")
}

func main() {
	log.Print("-Startup-")
	//InputFunctions()
	router := gin.Default()
	router.Use(static.Serve("/", static.LocalFile("./public/rs_audit/", true)))
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	})
	api := router.Group("/api")
	api.GET("/accounts/:account_id/arrays/:array_id/inputs", listArrayInputsEndpoint)
	api.GET("/accounts/:account_id/arrays/:array_id", getArrayEndpoint)
	api.GET("/accounts/:account_id/arrays",listArraysEndpoint)
	api.GET("/accounts/:account_id",getAccountEndpoint)
	api.GET("/accounts",listAccountsEndpoint)
	api.GET("/accounts/:account_id/arrays/:array_id/inputs/:input_name",listInputVersionsEndpoint)
	router.Run(":9080")
}

func listArrayInputsEndpoint(c *gin.Context) {
	account := c.Param("account_id")
	array := c.Param("array_id")
	inputs := listCurrentInputs(account,array)
	c.JSON(200,inputs)
}

func listArraysEndpoint(c *gin.Context) {
	account := c.Param("account_id")
	arrays := listArrays(account)
	c.JSON(200,arrays)
}

func listAccountsEndpoint(c *gin.Context) {
	accounts := listAccounts()
	c.JSON(200,accounts)
}

func getAccountEndpoint(c *gin.Context) {
	a := c.Param("account_id")
	account := getAccount(a)
	c.JSON(200,account)
}

func listInputVersionsEndpoint(c *gin.Context) {
	account := c.Param("account_id")
	array := c.Param("array_id")
	input := c.Param("input_name")
	inputs := listInputVersions(account,array,input)
	c.JSON(200,inputs)
}

func getArrayEndpoint(c *gin.Context)  {
	account := c.Param("account_id")
	a := c.Param("array_id")
	array := getArray(account,a)
	c.JSON(200,array)
}