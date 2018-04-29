package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
	"log"
	"os"
)

var GDB *gorm.DB

//specify needed tablename
//http://gorm.io/docs/conventions.html#Specifying-The-Table-Name

func DB() *gorm.DB {
	if GDB != nil {
		log.Println("GDB already defined")
		return GDB
	}
	db, err := gorm.Open("mysql", os.Getenv("MYSQL_CONNECTION_STRING"))
	if err != nil {
		log.Fatal(err)
	}
	//https://github.com/go-sql-driver/mysql/issues/461
	db.DB().SetMaxOpenConns(5)
	db.DB().SetConnMaxLifetime(time.Minute * 10)
	db.LogMode(true)
	GDB = db
	return db
}

