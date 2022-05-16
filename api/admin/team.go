package api

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/iamrahultanwar/heeico/database"
	"github.com/iamrahultanwar/heeico/model"
	"github.com/iamrahultanwar/heeico/types"
)

func TeamApi(app *fiber.App) {
	team := app.Group("/api/v1/admin/team")

	team.Get("/", func(c *fiber.Ctx) error {
		var users []model.User
		database.DB.Where("role", types.TEAM).Select("id", "email", "name", "created_at").Find(&users)
		return c.JSON(users)
	})

	team.Delete("/:id", func(c *fiber.Ctx) error {
		userId, _ := c.ParamsInt("id")
		database.DB.Where("id", userId).Unscoped().Delete(&model.User{})
		return c.JSON(types.SuccessResponse{Status: true, Message: "team member removed"})
	})

	team.Post("/create", func(c *fiber.Ctx) error {
		user := model.User{}
		err := c.BodyParser(&user)
		if err != nil {
			c.SendStatus(http.StatusBadRequest)
			return c.JSON(types.FailResponse{Status: false, Error: "Valid data not supplied"})
		}
		var userCount int64
		database.DB.Model(&user).Where("email", user.Email).Count(&userCount)
		if userCount > 0 {
			c.SendStatus(http.StatusBadRequest)
			return c.JSON(types.FailResponse{Status: false, Error: "User already created", Data: types.ResponseData{
				"error": fmt.Sprintf("email address %s already registered", user.Email),
			}})
		}
		user.HashPassword()
		user.Role = types.TEAM
		result := database.DB.Create(&user)
		if result.Error != nil {
			c.SendStatus(http.StatusInternalServerError)
			return c.JSON(types.FailResponse{Status: false, Error: result.Error.Error()})
		}
		return c.JSON(types.SuccessResponse{Status: true, Message: "Team member added"})
	})
}
