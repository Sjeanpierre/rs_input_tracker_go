package models

type Account struct {
	AccountID int `json:"account_id"`
	Name string `json:"name" gorm:"column:account_name"`
}

type Accounts []Account

func ListAccounts() Accounts {
  db := DB()
  var accounts Accounts
  db.Table("accounts").Find(&accounts)
  return accounts
}

func GetAccount(accountID string) Account {
	db := DB()
	var account Account
	db.Table("accounts").Where("account_id = ?",accountID).First(&account)
	return account
}