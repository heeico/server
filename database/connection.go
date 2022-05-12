package database

import (
	"github.com/iamrahultanwar/heeico/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	db, err := gorm.Open(sqlite.Open("url.sqlite"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	DB = db
	DB.AutoMigrate(&model.URLShortener{}, &model.URLTrack{}, &model.User{})
}
