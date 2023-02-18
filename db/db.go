package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type History struct {
	Hash    string `gorm:"index"`
	Line    uint64 `gorm:"index" gorm:"unique"`
	Message string
}

func ConnectToDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("gorc.db"), &gorm.Config{
		Logger:            logger.Default.LogMode(logger.Silent),
		AllowGlobalUpdate: true,
	})

	db.AutoMigrate(&History{})
	db.Delete(&History{})

	return db, err
}
