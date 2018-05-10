package models

import (
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

func ListArrays(account string) Arrays {
	db := DB()
	var arrays Arrays
	db.Table("server_arrays").Where("account_id = ?", account).Find(&arrays)
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
	matched := arrays.filterByArrayID(arrayID).filterByAccount(acct)
	if len(matched) != 0 {
		return matched[0]
	}
	return Array{}
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
