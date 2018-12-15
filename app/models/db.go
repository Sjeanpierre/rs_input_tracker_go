package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"os"
	"time"
)

var GDB *gorm.DB

func DB() *gorm.DB {
	if GDB != nil {
		return GDB
	}
	log.Println("No DB connection found, opening new one")
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



