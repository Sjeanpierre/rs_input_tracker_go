package models

import (
	"strconv"
	"time"
	"strings"
	"log"
	"fmt"
)

type Input struct {
	ID        int       `json:"-"`
	AccountID int       `json:"account_id"`
	ArrayID   int       `json:"array_id"`
	Type      string    `gorm:"-" json:"type"`
	Name      string    `gorm:"column:input_name" json:"name"`
	Value     string    `gorm:"-" json:"value"`
	RawValue  string    `gorm:"column:input_value" json:"-"`
	Version   int       `json:"version"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Inputs []Input

func (inputs Inputs) enrich() Inputs {
	for i, input := range inputs {
		valParts := strings.Split(input.RawValue, ":")
		if len(valParts) >= 2 {
			inputs[i].Type = valParts[0]
			inputs[i].Value = strings.Join(valParts[1:], ":")
		}
	}
	return inputs
}

func listTables() {
	log.Println("Listing tables")
	db := DB()
	rows, err := db.Raw("show tables").Rows()
	if err != nil {
		log.Fatalln("encountered error")
	}
	for rows.Next() {
		var res string
		rows.Scan(&res)
		log.Println(res)
	}
}

func ListInputVersions(account, array, input_name string) Inputs {
	db := DB()
	var inputs Inputs
	queryParams := map[string]interface{}{"account_id": account, "array_id": array, "input_name": input_name}
	db.Table(account).Where(queryParams).Find(&inputs)
	return inputs.enrich()
}

func ListInputs(account string, array string) Inputs {
	db := DB()
	var inputs Inputs
	db.Table(account).Where("array_id = ? AND account_id = ?", array, account).Find(&inputs)
	return inputs.enrich()
}

func ListCurrentInputs(account string, array string) Inputs {
	db := DB()
	var inputs Inputs
	q := fmt.Sprintf("SELECT * FROM `%s` t WHERE array_id = t.array_id AND"+
		" input_name = t.input_name AND version = (SELECT  MAX(version) FROM"+
		" `%s` WHERE array_id = t.array_id AND input_name = t.input_name)"+
		" AND array_id = %s AND account_id = %s", account, account, array, account)
	db.Table(account).Raw(q).Find(&inputs)
	return inputs.enrich()
}

func ListInputsAsOf(account string, array string, datetime time.Time) Inputs {
	db := DB()
	var inputs Inputs
	q := fmt.Sprintf("SELECT * FROM `%s` t WHERE array_id = t.array_id AND"+
		" input_name = t.input_name AND version = (SELECT  MAX(version) FROM"+
		" `%s` WHERE array_id = t.array_id AND input_name = t.input_name AND created_at <= %s)"+
		" AND array_id = %s AND account_id = %s", account, account, datetime, array, account)
	db.Table(account).Raw(q).Find(&inputs)
	return inputs.enrich()
}

func CreateInput(i Input, populateDate ...bool) {
	if len(populateDate) > 0 && populateDate[0] {
		currentTime := time.Now()
		i.CreatedAt = currentTime
		i.UpdatedAt = currentTime
	}
	db := DB()
	acc := strconv.Itoa(i.AccountID)
	db.Table(acc).Create(&i)
}
