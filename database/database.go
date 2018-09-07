package database

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	// _ "github.com/jinzhu/gorm/dialects/postgres"
	// _ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
)

type Database struct {
	*gorm.DB
}

var DB *gorm.DB

func InitDB(dialect string, dbUri string) *gorm.DB {
	db, err := gorm.Open(dialect, fmt.Sprintf("%s?charset=utf8&parseTime=True&loc=Local", dbUri))
	if err != nil {
		log.Fatal("db err: ", err)
	}
	db.DB().SetMaxIdleConns(10)
	DB = db
	return DB
}

func Migrate(db *gorm.DB) {
	log.Print("Migrating the database...")
	db.AutoMigrate(&Archive{})
}

// Using this function to get a connection, you can create your connection pool here.
func GetDB() *gorm.DB {
	return DB
}
