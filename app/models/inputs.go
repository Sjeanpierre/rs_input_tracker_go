package models

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
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
	Inactive  bool      `gorm:"inactive" json:"inactive"`
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

func ReactivateInput(input Input) {
	db := DB()
	db.Table(strconv.Itoa(input.AccountID)).Where("id = ?", strconv.Itoa(input.ID)).Updates(Input{UpdatedAt: time.Now(),Inactive:false})
}

//List all highest version inputs, including inputs marked as inactive
//this is needed for the worker in order to not start a new series for an input that was previously present
//but later marked inactive, then reintroduced.
//this function will only be needed internally
func ListCurrentInputsWithInactive(account string, array string) Inputs {
	db := DB()
	var inputs Inputs
	q := fmt.Sprintf("SELECT * FROM `%s` t WHERE array_id = t.array_id AND"+
		" input_name = t.input_name AND version = (SELECT  MAX(version) FROM"+
		" `%s` WHERE array_id = t.array_id AND input_name = t.input_name)"+
		" AND array_id = %s AND account_id = %s", account, account, array, account)
	db.Table(account).Raw(q).Find(&inputs)
	return inputs.enrich()
}

//List current inputs at their max version, excludes inputs marked as inactive
func ListCurrentInputs(account string, array string) Inputs {
	db := DB()
	var inputs Inputs
	q := fmt.Sprintf("SELECT * FROM `%s` t WHERE array_id = t.array_id AND"+
		" input_name = t.input_name AND version = (SELECT  MAX(version) FROM"+
		" `%s` WHERE array_id = t.array_id AND input_name = t.input_name)"+
		" AND array_id = %s AND account_id = %s AND inactive = false", account, account, array, account)
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

func CompareArrayInputs(account, array1, array2 string) map[string]Inputs {
	ar1 := ListCurrentInputs(account, array1)
	ar2 := ListCurrentInputs(account, array2)
	missingFromAR1 := ar1.diff(ar2)
	missingFromAR2 := ar2.diff(ar1)
	allMissing := append(missingFromAR1, missingFromAR2...)
	// ar1 = 1234
	// ar2 = 1256
	// ar1.diff(ar2) = 5,6
	// ar2.diff(ar1) = 3,4
	ret := make(map[string]Inputs)
	a1 := make(map[string]Input)
	a2 := make(map[string]Input)
	for _, input := range ar1 {
		a1[input.Name] = input
	}
	for _, input := range ar2 {
		a2[input.Name] = input
	}
	for _, ipt := range allMissing {
		ret[ipt.Name] = Inputs{a1[ipt.Name], a2[ipt.Name]}
	}
	return ret
}

func (ar1 Inputs) diff(ar2 Inputs) Inputs {
	lookUp := make(map[string]bool)
	for _, input := range ar1 {
		lookUp[input.inputSignature()] = true
	}
	var missing Inputs //exists in ar2, but not ar1
	for _, input := range ar2 {
		_, ok := lookUp[input.inputSignature()]
		if !ok {
			missing = append(missing, input)
		}
	}
	return missing
}

func (i Input) inputSignature() string {
	return fmt.Sprintf("%s|%s|%s", i.Name, i.Type, i.Value)
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

func CreateInactiveInputRecord(i Input, populateDate ...bool) {
	if len(populateDate) > 0 && populateDate[0] {
		currentTime := time.Now()
		i.CreatedAt = currentTime
		i.UpdatedAt = currentTime
	}
	i.Inactive = true
	db := DB()
	acc := strconv.Itoa(i.AccountID)
	db.Table(acc).Create(&i)
}