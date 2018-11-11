package models

import (
	"fmt"
	"time"
)

type Array struct {
	ID        int       `json:"-"`
	AccountID int       `json:"account_id"`
	ArrayID   int       `json:"array_id"`
	Name      string    `json:"array_name" gorm:"column:array_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Arrays []Array

func ListArraysByAccount(account int) Arrays {
	db := DB()
	var arrays Arrays
	db.Table("server_arrays").Where("account_id = ?", account).Find(&arrays)
	return arrays
}

func ListCurrentArraysByAccount(account string)  Arrays{
	db := DB()
	var arrays Arrays
	q := fmt.Sprintf("select * from server_arrays sa" + " where account_id = %s AND created_at ="+
		" (select MAX(created_at) from server_arrays where array_id = sa.array_id)",account)
	db.Table("server_arrays").Raw(q).Find(&arrays)
	return arrays
}

func ListArrayVersions(account,array string)  Arrays{
	db := DB()
	var arrays Arrays
	q := fmt.Sprintf("SELECT c.*, repeated FROM" +
		" (SELECT array_id, COUNT(array_id) AS repeated FROM" +
		" preprod_inputtracker.server_arrays GROUP BY array_id)" +
		" AS sa JOIN server_arrays c ON sa.array_id = c.array_id WHERE c.array_id = %s" +
		" AND c.account_id = %s",array,account)
	db.Table("server_arrays").Raw(q).Find(&arrays)
	return arrays
}

func ListAllArrays() Arrays {
	db := DB()
	var arrays Arrays
	db.Table("server_arrays").Find(&arrays)
	return arrays
}

func GetArray(account string, array string) Array {
	db := DB()
	var a Array
	db.Table("server_arrays").Where("account_id = ? AND array_id = ?", account, array).First(&a)
	return a
}

func CreateArray(array Array) {
	currentTime := time.Now()
	array.CreatedAt = currentTime
	array.UpdatedAt = currentTime
	db := DB()
	db.Table("server_arrays").Create(&array)
}

func UpdateArrays(arrays Arrays) {
	db := DB()
	var ids []int
	for _, array := range arrays {
		i := array.ID
		if i != 0 {
			ids = append(ids, array.ID)
		}
	}
	if len(ids) > 0 {
		db.Table("server_arrays").Where("id IN (?)", ids).Update("updated_at", time.Now())
	}
}

func (arrays Arrays) Exists(acct int, arrayID int) bool {
	matched := arrays.filterByArrayID(arrayID).filterByAccount(acct)
	if len(matched) == 1 {
		return true
	}
	if len(matched) > 1 {
		//log.Printf("Error: more than one array exists with ID: %v for account %v in DB", arrayID, acct)
		return true
	}
	return false
}

func (arrays Arrays) FindArrayByAccountArrayID(acct int, arrayID int) Array {
	return arrays.filterByArrayID(arrayID).filterByAccount(acct).filterLatest()
}

func (arrays Arrays) filterByAccount(accountID int) Arrays {
	var a2 Arrays
	for _, array := range arrays {
		if array.AccountID == accountID {
			a2 = append(a2, array)
		}
	}
	return a2
}

func (arrays Arrays) filterByArrayID(ArrayID int) Arrays {
	var a2 Arrays
	for _, array := range arrays {
		if array.ArrayID == ArrayID {
			a2 = append(a2, array)
		}
	}
	return a2
}

func (arrays Arrays) filterLatest() Array  {
	var a Array
	for _, arr := range arrays {
		if arr.ID > a.ID {
			a = arr
		}
	}
	return a
}
