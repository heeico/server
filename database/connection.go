package database

import (
	"fmt"

	"github.com/iamrahultanwar/heeico/config"
	"github.com/iamrahultanwar/heeico/model"
	"github.com/iamrahultanwar/heeico/types"
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

func Init() {
	var count int64
	DB.Model(&model.User{}).Count(&count)
	con := config.GetConfig()
	if count == 0 {
		adminUser := model.User{
			Name:     con.Admin.Name,
			Email:    con.Admin.Email,
			Password: con.Admin.Password,
			Role:     types.ADMIN,
		}
		adminUser.HashPassword()
		DB.Create(&adminUser)
		fmt.Println("Admin user created")
	}

	if len(con.Links) > 0 {
		var urlLinks []model.URLShortener
		for _, link := range con.Links {
			currentLink := model.URLShortener{
				Alias:       link.Alias,
				RedirectURL: link.RedirectURL,
			}
			if !AliasExist(currentLink) {
				urlLinks = append(urlLinks, currentLink)
			}
		}
		if len(urlLinks) > 0 {
			DB.CreateInBatches(urlLinks, 100)
		}
	}
}
