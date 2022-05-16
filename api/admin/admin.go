package api

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"

	"github.com/iamrahultanwar/heeico/database"
	"github.com/iamrahultanwar/heeico/model"
	"github.com/iamrahultanwar/heeico/types"
)

func AdminApi(app *fiber.App) {
	// admin operations
	admin := app.Group("/api/v1/admin")

	admin.Get("/dashboard", func(c *fiber.Ctx) error {
		var totalLinks int64 = 0  // url shortener
		var totalClicks int64 = 0 // url track

		linkTrack := []model.TotalVisit{}

		database.DB.Model(&model.URLShortener{}).Count(&totalLinks)
		database.DB.Model(&model.URLTrack{}).Count(&totalClicks)
		database.DB.Model(&model.URLTrack{}).Select("created_at, count(*) as total").Group("cast(created_at as date)").Find(&linkTrack)

		return c.JSON(types.ResponseData{
			"track":       linkTrack,
			"totalLinks":  totalLinks,
			"totalClicks": totalClicks,
		})

	})

	admin.Post("/create-url", func(c *fiber.Ctx) error {
		var urlShortener model.URLShortener
		if err := c.BodyParser(&urlShortener); err != nil {
			c.SendStatus(http.StatusBadRequest)
			return c.JSON(types.SuccessResponse{Message: "Invalid data", Status: false})
		}
		result := database.DB.Create(&urlShortener)
		if result.Error != nil {
			c.SendStatus(http.StatusBadRequest)
			return c.JSON(types.FailResponse{Error: result.Error.Error(), Status: false})
		} else {
			return c.JSON(types.SuccessResponse{Message: "Alias Created", Status: true})
		}
	})

	type Query struct {
		Search string `json:"search"`
		Limit  int    `json:"limit default:10"`
	}
	admin.Get("/get-url", func(c *fiber.Ctx) error {
		query := Query{}
		c.QueryParser(&query)
		fmt.Println(query)
		var urlShorteners []model.URLShortener
		result := database.DB.Where("alias LIKE ?", "%"+query.Search+"%").Order("created_at desc").Find(&urlShorteners)
		if result.Error != nil {
			c.SendStatus(http.StatusBadRequest)
			return c.JSON(types.FailResponse{Error: result.Error.Error(), Status: false})
		}
		return c.JSON(urlShorteners)
	})

	admin.Get("/get-url-tracking/:urlId", func(c *fiber.Ctx) error {
		// var urlTracks []URLTrack
		urlID := c.Params("urlId")

		var totalVisits []model.TotalVisit

		result := database.DB.Model(&model.URLTrack{}).Select("created_at, count(*) as total").Where("url_shortener_id", urlID).Group("created_at").Find(&totalVisits)

		if result.Error != nil {
			c.SendStatus(http.StatusBadRequest)
			return c.JSON(types.FailResponse{Error: result.Error.Error(), Status: false})
		}
		return c.JSON(totalVisits)

	})

	admin.Delete("/delete-url/:urlId", func(c *fiber.Ctx) error {
		urlID := c.Params("urlId")
		result := database.DB.Delete(&model.URLShortener{}, urlID)

		if result.Error != nil {
			c.SendStatus(http.StatusBadRequest)
			return c.JSON(types.FailResponse{Error: result.Error.Error(), Status: false})
		}

		return c.JSON(types.SuccessResponse{Message: "Link deleted", Status: true, Data: types.ResponseData{"message": fmt.Sprintf("URL with id %s deleted", urlID)}})
	})

	admin.Put("/update-url/:urlId", func(c *fiber.Ctx) error {
		urlID := c.Params("urlId")
		var urlShortener model.URLShortener

		result := database.DB.Where("id", urlID).First(&urlShortener)

		if result.Error != nil {
			c.SendStatus(http.StatusBadRequest)
			return c.JSON(types.FailResponse{Error: result.Error.Error(), Status: false})
		}

		if err := c.BodyParser(&urlShortener); err != nil {
			c.SendStatus(http.StatusBadRequest)
			return c.JSON(types.SuccessResponse{Message: "Invalid data", Status: false})
		}

		go database.DB.Save(&urlShortener)

		return c.JSON(types.SuccessResponse{Message: "Link Updated", Status: true, Data: types.ResponseData{"updatedFields": urlShortener}})
	})

	admin.Get("/user-profile", func(c *fiber.Ctx) error {
		userId := c.Locals("userId")
		if userId == nil {
			c.SendStatus(http.StatusBadRequest)
			return c.JSON(types.FailResponse{Status: false, Error: "user not valid"})
		}
		user := model.User{}
		database.DB.Where("id", userId).First(&user)
		return c.JSON(user)
	})

	admin.Put("/update-profile", func(c *fiber.Ctx) error {
		userId := c.Locals("userId")
		if userId == nil {
			c.SendStatus(http.StatusBadRequest)
			return c.JSON(types.FailResponse{Status: false, Error: "user not valid"})
		}
		user := model.User{}
		c.BodyParser(&user)
		go database.DB.Model(&model.User{}).Where("id", userId).Update("name", user.Name).Update("email", user.Email)
		return c.JSON(types.SuccessResponse{Status: true, Message: "Profile Updated"})
	})

	admin.Put("/update-profile-password", func(c *fiber.Ctx) error {
		userId := c.Locals("userId")
		if userId == nil {
			c.SendStatus(http.StatusBadRequest)
			return c.JSON(types.FailResponse{Status: false, Error: "user not valid"})
		}
		type PasswordInput struct {
			OldPassword string `json:"oldPassword"`
			NewPassword string `json:"newPassword"`
		}
		userPassword := PasswordInput{}
		user := model.User{}
		c.BodyParser(&userPassword)
		oldUser := model.User{
			Password: userPassword.OldPassword,
		}
		database.DB.Where("id", userId).First(&user)
		e := oldUser.CheckPasswordHash(user.Password)
		if e != nil {
			c.SendStatus(http.StatusBadRequest)
			return c.JSON(types.FailResponse{Status: false, Error: "password does not match"})
		}
		user.Password = userPassword.NewPassword
		user.HashPassword()
		go database.DB.Model(&model.User{}).Where("id", userId).Update("password", user.Password)
		return c.JSON(types.SuccessResponse{Status: true, Message: "Profile Updated"})
	})
}
