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

}
