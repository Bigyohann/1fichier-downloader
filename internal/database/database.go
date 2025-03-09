package database

import (
	"bigyohann/apidownloader/internal/database/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var dbConn *gorm.DB

func InitDatabase() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("download.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	dbConn = db

	models.InitDbFileModel(db)

	return db
}

func GetDB() *gorm.DB { return dbConn }
