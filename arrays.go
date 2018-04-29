package main

import "time"

type Array struct {
	AccountID int `json:"account_id"`
	ArrayID int `json:"array_id"`
	Name string `json:"array_name" gorm:"column:array_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Arrays []Array

func listArrays(account string) Arrays{
	db := DB()
	var arrays Arrays
	db.Table("server_arrays").Where("account_id = ?",account).Find(&arrays)
	return arrays
}

func getArray(account string, array string) Array {
	db := DB()
	var a Array
	db.Table("server_arrays").Where("account_id = ? AND array_id = ?",account,array).First(&a)
	return a
}
